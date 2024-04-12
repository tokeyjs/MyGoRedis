package basedatastruct

import (
	"MyGoRedis/lib/logger"
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const MaxH = 32

// 跳表 （含头节点）（元素默认从大到小）（使用时确保不存在重复元素）
type BaseSkipLink struct {
	head      *skNode // 头节点（不包含数据）
	tail      *skNode
	rmutex    sync.RWMutex //并发锁
	size      atomic.Int32 //跳表元素个数
	maxHeight atomic.Int32 // 最大高度
}

func MakeSkipLink() *BaseSkipLink {
	head := &skNode{
		score:     0,
		value:     "",
		levelNext: make([]*skNode, 1),
		prev:      nil,
	}
	Val := &BaseSkipLink{
		head: head,
		tail: head,
	}
	Val.size.Store(0)
	Val.maxHeight.Store(1)
	return Val
}

type Valtmp struct {
	Score float32
	Value string
}

type skNode struct {
	score     float32
	value     string
	levelNext []*skNode //下一个节点指针
	prev      *skNode   // 前一个节点指针
}

// 获取元素数量
func (skip *BaseSkipLink) Size() int32 {
	return skip.size.Load()
}

// 增加元素
func (skip *BaseSkipLink) AddElem(value string, score float32) {
	skip.rmutex.Lock()
	defer skip.rmutex.Unlock()
	skip.addElem(value, score)
}
func (skip *BaseSkipLink) addElem(value string, score float32) {
	node := skip.head
	// 获取随机高度
	rh := skip.randomHeight()
	h := skip.maxHeight.Load()
	if rh > h {
		// 更新最大高度
		skip.maxHeight.Store(rh)
		skip.head.levelNext = append(skip.head.levelNext, make([]*skNode, rh-h)...)
	}
	newNode := &skNode{
		score:     score,
		value:     value,
		levelNext: make([]*skNode, rh),
		prev:      nil,
	}
	for i := rh - 1; i >= 0; i-- {
		for node.levelNext[i] != nil && node.levelNext[i].score > newNode.score {
			node = node.levelNext[i]
		}
		last := node.levelNext[i]
		if last == nil {
			node.levelNext[i] = newNode
			newNode.prev = node
			newNode.levelNext[i] = last
		} else {
			node.levelNext[i] = newNode
			newNode.prev = node
			newNode.levelNext[i] = last
			last.prev = newNode
		}
	}
	skip.size.Add(1)
}

// 删除指定元素 删除成功为true
func (skip *BaseSkipLink) DelElem(value string, score float32) bool {
	skip.rmutex.Lock()
	defer skip.rmutex.Unlock()
	return skip.delElem(value, score)
}
func (skip *BaseSkipLink) delElem(value string, score float32) bool {
	// 搜索节点
	node := skip.secNode(value, score)
	if node == nil {
		return false
	}
	// 将指向该节点的所有指针替换
	// 该元素为尾节点
	if node.levelNext[0] == nil {
		h := len(node.levelNext) - 1
		tmpNode := skip.head
		for ; h >= 0; h-- {
			for tmpNode.levelNext[h] != node {
				tmpNode = tmpNode.levelNext[h]
			}
			tmpNode.levelNext[h] = nil
		}
		skip.tail = node.prev
	} else {
		h := len(node.levelNext) - 1
		tmpNode := skip.head
		for ; h >= 0; h-- {
			for tmpNode.levelNext[h] != node {
				tmpNode = tmpNode.levelNext[h]
			}
			tmpNode.levelNext[h] = node.levelNext[h]
		}
		node.levelNext[0].prev = node.prev
	}
	skip.size.Add(-1)
	return true
}

// 搜索指定节点 （非并发安全）
func (skip *BaseSkipLink) secNode(value string, score float32) *skNode {
	// 搜索value的node节点
	h := skip.maxHeight.Load() - 1
	node := skip.head
	for ; h >= 0; h-- {
		for node.levelNext[h] != nil && node.levelNext[h].score > score {
			node = node.levelNext[h]
		}
		if node.levelNext[h] == nil || node.levelNext[h].score < score {
			continue
		}
		if node.levelNext[h].score == score {
			if node.levelNext[h].value == value {
				return node.levelNext[h]
			} else {
				node = node.levelNext[h]
				continue
			}
		}
	}
	return nil
}

// 返回所有成员
func (skip *BaseSkipLink) GetAllElem() ([]*Valtmp, error) {
	skip.rmutex.RLock()
	defer skip.rmutex.RUnlock()
	res := make([]*Valtmp, 0, skip.Size())
	node := skip.head
	for node.levelNext[0] != nil {
		res = append(res, &Valtmp{
			Score: node.levelNext[0].score,
			Value: node.levelNext[0].value,
		})
		node = node.levelNext[0]
	}
	return res, nil
}

// 返回指定值区间的元素 [ ]
func (skip *BaseSkipLink) GetScoreRangeElem(minScore, maxScore float32) ([]*Valtmp, error) {
	skip.rmutex.RLock()
	defer skip.rmutex.RUnlock()
	res := make([]*Valtmp, 0)
	// 获取元素
	elems, err := skip.getScoreRangeElem(minScore, maxScore)
	if err != nil {
		return nil, err
	}
	for _, node := range elems {
		res = append(res, &Valtmp{
			Score: node.score,
			Value: node.value,
		})
	}
	return res, nil
}
func (skip *BaseSkipLink) getScoreRangeElem(minScore, maxScore float32) ([]*skNode, error) {
	res := make([]*skNode, 0)
	if minScore > maxScore {
		return nil, errors.New("minScore >= maxScore")
	}
	node := skip.head
	// 快速接近node
	h := skip.maxHeight.Load() - 1
	for ; h >= 0; h-- {
		for node.levelNext[h] != nil && node.levelNext[h].score > maxScore {
			node = node.levelNext[h]
		}
	}
	for node.levelNext[0] != nil {
		if node.levelNext[0].score > maxScore {
			node = node.levelNext[0]
		} else if node.levelNext[0].score < minScore {
			break
		} else {
			res = append(res, node.levelNext[0])
			node = node.levelNext[0]
		}
	}
	return res, nil
}

// 返回指定排名(0-n)(从大到小排)之间的元素 [ ]: 后续可以进行优化span
func (skip *BaseSkipLink) GetRankRangeElem(minRank, maxRank int32) ([]*Valtmp, error) {
	skip.rmutex.RLock()
	defer skip.rmutex.RUnlock()
	res := make([]*Valtmp, 0)
	elem, err := skip.getRankRangeElem(minRank, maxRank)
	if err != nil {
		return nil, err
	}
	for _, node := range elem {
		res = append(res, &Valtmp{
			Score: node.score,
			Value: node.value,
		})
	}
	return res, nil
}
func (skip *BaseSkipLink) getRankRangeElem(minRank, maxRank int32) ([]*skNode, error) {
	skip.rmutex.RLock()
	defer skip.rmutex.RUnlock()
	res := make([]*skNode, 0)
	if minRank > maxRank {
		return nil, errors.New("minScore >= maxScore")
	}
	idx := int32(0)
	node := skip.head
	for node.levelNext[0] != nil {
		if idx > maxRank {
			break
		} else if idx < minRank {
			node = node.levelNext[0]
		} else {
			res = append(res, node.levelNext[0])
			node = node.levelNext[0]
		}
		idx++
	}
	return res, nil
}

// 获取指定元素索引(排名)O(n) : 暂时未优化，加上span会降低时间复杂度至O(logn)
func (skip *BaseSkipLink) GetIndexElem(value string) (int32, error) {
	skip.rmutex.RLock()
	defer skip.rmutex.RUnlock()
	idx := int32(0)
	node := skip.head
	for node.levelNext[0] != nil {
		if node.levelNext[0].value == value {
			return idx, nil
		}
		idx++
	}
	return -1, errors.New("not found ")
}

// 更新成员分数
func (skip *BaseSkipLink) UpdateScoreElem(value string, srcScore float32, destScore float32) error {
	skip.rmutex.RLock()
	defer skip.rmutex.RUnlock()
	node := skip.secNode(value, srcScore)
	if node == nil {
		return errors.New("not found")
	}
	node.score = destScore
	// 调整位置
	skip.deal(node)
	return nil
}

// 为node找到新位置（调整跳表）
func (skip *BaseSkipLink) deal(node *skNode) {
	// 情况不需要更新: 处于头位置正确，处于尾位置正确，处于中间位置正确
	if (skip.Size() == 1) || (node.prev == skip.head && node.levelNext[0].score <= node.score) ||
		(node.levelNext[0] == nil && node.prev.score >= node.score) ||
		(node.prev != skip.head && node.levelNext[0] != nil && node.prev.score >= node.score && node.levelNext[0].score <= node.score) {
		return
	}
	//删除再插入
	skip.delElem(node.value, node.score)
	skip.addElem(node.value, node.score)
}

// 删除指定排名范围的成员：TODO 后续优化
func (skip *BaseSkipLink) RemoveRangeRank(rankStart, rankEnd int32) int32 {
	if rankStart > rankEnd {
		return 0
	}
	// 获取成员
	slice, err := skip.getRankRangeElem(rankStart, rankEnd)
	if err != nil {
		return 0
	}
	for _, node := range slice {
		skip.delElem(node.value, node.score)
	}
	return int32(len(slice))
}

// 删除指定分数区间的成员： TODO 后续优化
func (skip *BaseSkipLink) RemoveRangeScore(minScore, maxScore float32) int32 {
	if minScore > maxScore {
		return 0
	}
	// 获取成员
	slice, err := skip.getScoreRangeElem(minScore, maxScore)
	if err != nil {
		return 0
	}
	for _, node := range slice {
		skip.delElem(node.value, node.score)
	}
	return int32(len(slice))
}

// 获取一个随机高度
func (skip *BaseSkipLink) randomHeight() int32 {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(MaxH)
	return int32(randomNum)
}

// 统计分数区间人数数量
func (skip *BaseSkipLink) CountRangeScore(scoreMin, scoreMax float32) int32 {
	skip.rmutex.RLock()
	defer skip.rmutex.RUnlock()
	elem, err := skip.getScoreRangeElem(scoreMin, scoreMax)
	if err != nil {
		return 0
	}
	return int32(len(elem))
}

// 对指定成员score加增量
func (skip *BaseSkipLink) Incr(value string, score float32, incr float32) error {
	skip.rmutex.Lock()
	defer skip.rmutex.Unlock()
	node := skip.secNode(value, score)
	if node == nil {
		return errors.New("not found")
	}
	node.score += incr
	return nil
}

func Show_skip(sk *BaseSkipLink) {
	res, err := sk.GetAllElem()
	if err != nil {
		logger.Info("error :" + err.Error())
	}
	logger.Info("SHOW")
	for _, v := range res {
		logger.Infof("value:%v Score:%v", v.Value, v.Score)
	}
}

// 测试
func Test_skip() {

}
