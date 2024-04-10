package basedatastruct

import (
	"MyGoRedis/lib/logger"
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const MaxH = 256

// 跳表 （元素默认从大到小）
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
	score float32
	value string
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

// 获取分数
func (skip *BaseSkipLink) GetScore(value string) (float32, error) {
	skip.rmutex.RLock()
	defer skip.rmutex.RUnlock()
	return skip.getScore(value)
}
func (skip *BaseSkipLink) getScore(value string) (float32, error) {
	node := skip.head
	for node.levelNext[0] != nil && node.levelNext[0].value != value {
		node = node.levelNext[0]
	}
	if node.levelNext[0] == nil {
		return 0, errors.New("not found")
	}
	return node.levelNext[0].score, nil
}

// 增加元素
func (skip *BaseSkipLink) AddElem(value string, score float32) error {
	skip.rmutex.Lock()
	defer skip.rmutex.Unlock()
	return skip.addElem(value, score)
}
func (skip *BaseSkipLink) addElem(value string, score float32) error {
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
		node.levelNext[i] = newNode
		newNode.prev = node
		newNode.levelNext[i] = last
	}
	skip.size.Add(1)
	return nil
}

// 删除指定元素
func (skip *BaseSkipLink) DelElem(value string) error {
	skip.rmutex.Lock()
	defer skip.rmutex.Unlock()
	return skip.delElem(value)
}
func (skip *BaseSkipLink) delElem(value string) error {
	node := skip.head
	for node.levelNext[0] != nil && node.levelNext[0].value != value {
		node = node.levelNext[0]
	}
	if node.levelNext[0] == nil {
		// 没有此元素 无法进行删除
		return errors.New("not found elem, delete failed")
	}
	node = node.levelNext[0]
	// 将下一个节点的prev赋值
	if node.levelNext[0] != nil {
		node.levelNext[0].prev = node.prev
	}
	// 将levelNext赋值跳过node
	tmpNode := node.prev
	for i := 0; i < len(node.levelNext); {
		if i < len(tmpNode.levelNext) && tmpNode.levelNext[i] != nil {
			tmpNode.levelNext[i] = node.levelNext[i]
			i++
		} else {
			tmpNode = tmpNode.prev
		}
	}
	skip.size.Add(-1)
	return nil
}

func (skip *BaseSkipLink) secNode(value string, score float32) *skNode {
	// 搜索value的node节点
	h := skip.maxHeight.Load()
	node := skip.head
	for ; h >= 0; h-- {
		for node.levelNext[h] != nil && node.levelNext[h].score > score {
			node = node.levelNext[h]
		}
		if node.levelNext[h] == nil {
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
			score: node.levelNext[0].score,
			value: node.levelNext[0].value,
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
			res = append(res, &Valtmp{
				score: node.levelNext[0].score,
				value: node.levelNext[0].value,
			})
			node = node.levelNext[0]
		}
	}
	return res, nil
}

// 返回指定排名(0-n)(从大到小排)之间的元素 [ ]
func (skip *BaseSkipLink) GetRankElem(minRank, maxRank int32) ([]*Valtmp, error) {
	skip.rmutex.RLock()
	defer skip.rmutex.RUnlock()
	res := make([]*Valtmp, 0)
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
			res = append(res, &Valtmp{
				score: node.levelNext[0].score,
				value: node.levelNext[0].value,
			})
			node = node.levelNext[0]
		}
		idx++
	}
	return res, nil
}

// 获取指定元素索引(排名)
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
func (skip *BaseSkipLink) UpdateScoreElem(value string, score float32) error {
	skip.rmutex.RLock()
	defer skip.rmutex.RUnlock()
	node := skip.head
	for node.levelNext[0] != nil {
		if node.levelNext[0].value == value {
			oldScore := node.levelNext[0].score
			node.levelNext[0].score = score
			if err := skip.deal(node.levelNext[0]); err != nil {
				node.levelNext[0].score = oldScore
				return err
			}
			return nil
		}
		node = node.levelNext[0]
	}
	return errors.New("not found")
}

// 为node找到新位置（调整跳表）
func (skip *BaseSkipLink) deal(node *skNode) error {
	// 情况不需要更新: 处于头位置正确，处于尾位置正确，处于中间位置正确
	if (skip.Size() == 1) || (node.prev == skip.head && node.levelNext[0].score <= node.score) ||
		(node.levelNext[0] == nil && node.prev.score >= node.score) ||
		(node.prev != skip.head && node.levelNext[0] != nil && node.prev.score >= node.score && node.levelNext[0].score <= node.score) {
		return nil
	}
	// 删除节点插入节点
	if err := skip.delElem(node.value); err != nil {
		return errors.New("modify error")
	}
	// 插入节点
	if err := skip.addElem(node.value, node.score); err != nil {
		return errors.New("modify error")
	}
	return nil
}

// 删除指定排名范围的成员
func (skip *BaseSkipLink) RemoveRangeRank(rankStart, rankEnd int32) {
	// todo
}

// 删除指定分数区间的成员
func (skip *BaseSkipLink) RemoveRangeScore(minScore, maxScore float32) {
	// todo
	// 快速接近

	// 查找max第一个元素

	// 查找区间最后一个元素

}

// 获取一个随机高度
func (skip *BaseSkipLink) randomHeight() int32 {
	baseH := int32(1)
	for true {
		if baseH == MaxH {
			break
		}
		rand.Seed(time.Now().UnixNano())
		randomNum := rand.Intn(99) % 2
		if randomNum == 0 {
			break
		}
		baseH++
	}
	return baseH
}

// 统计分数区间人数数量
func (skip *BaseSkipLink) CountRangeScore(scoreMin, scoreMax float32) int32 {
	skip.rmutex.RLock()
	defer skip.rmutex.RUnlock()
	if skip.Size() == 0 || scoreMin > skip.head.levelNext[0].score || scoreMax < skip.tail.score {
		return 0
	}
	node := skip.head
	// 快速接近
	h := skip.maxHeight.Load() - 1
	for ; h >= 0; h-- {
		if node.levelNext[h] != nil && node.levelNext[h].score > scoreMax {
			node = node.levelNext[h]
		}
	}
	num := int32(0)
	for node.levelNext[0] != nil {
		if node.levelNext[0].score < scoreMin {
			return num
		}
		if node.levelNext[0].score >= scoreMin && node.levelNext[0].score <= scoreMax {
			num++
		}
		node = node.levelNext[0]
	}
	return num
}

// 对指定成员score加增量
func (skip *BaseSkipLink) Incr(value string, incr float32) error {
	skip.rmutex.Lock()
	defer skip.rmutex.Unlock()
	node := skip.head
	for node.levelNext[0] != nil && node.levelNext[0].value != value {
		node = node.levelNext[0]
	}
	if node.levelNext[0] != nil {
		node.levelNext[0].score += incr
		return nil
	}
	return errors.New("not found")
}

func Show_skip(sk *BaseSkipLink) {
	res, err := sk.GetAllElem()
	if err != nil {
		logger.Info("error :" + err.Error())
	}
	logger.Info("SHOW")
	for _, v := range res {
		logger.Infof("value:%v score:%v", v.value, v.score)
	}
}

// 测试
func Test_skip() {
	sk := MakeSkipLink()
	sk.AddElem("a", 1234)
	Show_skip(sk)
	sk.AddElem("b", 124)
	Show_skip(sk)
	sk.AddElem("c", 7234)
	Show_skip(sk)
	sk.AddElem("d", 1264)
	Show_skip(sk)
	sk.DelElem("c")
	Show_skip(sk)
	sk.UpdateScoreElem("d", 9.0)
	Show_skip(sk)
	sk.AddElem("tt1", -9.4)
	Show_skip(sk)
	sk.AddElem("b1", 124)
	Show_skip(sk)
	sk.AddElem("c1", 7234)
	Show_skip(sk)
	sk.AddElem("d1", 1264)
	Show_skip(sk)
}
