package basedatastruct

import (
	"MyGoRedis/lib/logger"
	"errors"
	"math/rand"
	"time"
)

// 跳表 （元素默认从大到小）
type BaseSkipLink struct {
	head      *skNode // 头节点（不包含数据）
	tail      *skNode
	size      int32 //跳表元素个数
	maxHeight int32 // 最大高度
}

func MakeSkipLink() *BaseSkipLink {
	head := &skNode{
		score:     0,
		value:     "",
		levelNext: make([]*skNode, 1),
		levelPre:  make([]*skNode, 1),
	}
	return &BaseSkipLink{
		head:      head,
		tail:      head,
		size:      0,
		maxHeight: 1,
	}
}

type valtmp struct {
	score float32
	value string
}

type skNode struct {
	score     float32
	value     string
	levelNext []*skNode //下一个节点指针
	levelPre  []*skNode // 前一个节点指针
}

// 获取元素数量
func (skip *BaseSkipLink) Size() int32 {
	return skip.size
}

// 增加元素
func (skip *BaseSkipLink) AddElem(value string, score float32) error {
	node := skip.head
	// 获取随机高度
	rh := skip.randomHeight()
	logger.Infof("h=%v", rh)
	if rh > skip.maxHeight {
		// 更新最大高度
		skip.head.levelPre = append(skip.head.levelPre, make([]*skNode, rh-skip.maxHeight)...)
		skip.head.levelNext = append(skip.head.levelNext, make([]*skNode, rh-skip.maxHeight)...)
		skip.maxHeight = rh
	}
	newNode := &skNode{
		score:     score,
		value:     value,
		levelNext: make([]*skNode, rh),
		levelPre:  make([]*skNode, rh),
	}
	for i := rh - 1; i >= 0; i-- {
		for node.levelNext[i] != nil && node.levelNext[i].score > newNode.score {
			node = node.levelNext[i]
		}
		last := node.levelNext[i]
		node.levelNext[i] = newNode
		newNode.levelPre[i] = node
		if last != nil {
			last.levelPre[i] = newNode
		}
		newNode.levelNext[i] = last
	}
	skip.size++
	return nil
}

// 删除指定元素
func (skip *BaseSkipLink) DelElem(value string) error {
	node := skip.head
	for node.levelNext[0] != nil && node.levelNext[0].value != value {
		node = node.levelNext[0]
	}
	if node.levelNext[0] == nil {
		// 没有此元素 无法进行删除
		return errors.New("not found elem, delete failed")
	}
	node = node.levelNext[0]
	for i := 0; i < len(node.levelNext); i++ {
		if node.levelNext[i] != nil {
			node.levelNext[i].levelPre[i] = node.levelPre[i]
			node.levelPre[i].levelNext[i] = node.levelNext[i]
		} else {
			node.levelPre[i].levelNext[i] = nil
		}
	}
	skip.size--
	return nil
}

// 返回所有成员
func (skip *BaseSkipLink) GetAllElem() ([]*valtmp, error) {
	res := make([]*valtmp, 0, skip.Size())
	node := skip.head
	for node.levelNext[0] != nil {
		res = append(res, &valtmp{
			score: node.levelNext[0].score,
			value: node.levelNext[0].value,
		})
		node = node.levelNext[0]
	}
	return res, nil
}

// 返回指定值区间的元素 [ ]
func (skip *BaseSkipLink) GetScoreRangeElem(minScore, maxScore float32) ([]*valtmp, error) {
	res := make([]*valtmp, 0)
	if minScore > maxScore {
		return nil, errors.New("minScore >= maxScore")
	}
	node := skip.head
	for node.levelNext[0] != nil {
		if node.levelNext[0].score > maxScore {
			node = node.levelNext[0]
		} else if node.levelNext[0].score < minScore {
			break
		} else {
			res = append(res, &valtmp{
				score: node.levelNext[0].score,
				value: node.levelNext[0].value,
			})
			node = node.levelNext[0]
		}
	}
	return res, nil
}

// 返回指定排名(0-n)(从大到小排)之间的元素 [ ]
func (skip *BaseSkipLink) GetRankElem(minRank, maxRank int32) ([]*valtmp, error) {
	res := make([]*valtmp, 0)
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
			res = append(res, &valtmp{
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
	if (skip.Size() == 1) || (node.levelPre[0] == skip.head && node.levelNext[0].score <= node.score) ||
		(node.levelNext[0] == nil && node.levelPre[0].score >= node.score) ||
		(node.levelPre[0] != skip.head && node.levelNext[0] != nil && node.levelPre[0].score >= node.score && node.levelNext[0].score <= node.score) {
		return nil
	}
	// 删除节点插入节点
	if err := skip.DelElem(node.value); err != nil {
		return errors.New("modify error")
	}
	// 插入节点
	if err := skip.AddElem(node.value, node.score); err != nil {
		return errors.New("modify error")
	}
	return nil
}

// 获取一个随机高度
func (skip *BaseSkipLink) randomHeight() int32 {
	baseH := int32(1)
	for true {
		rand.Seed(time.Now().UnixNano())
		randomNum := rand.Intn(99) % 2
		if randomNum == 0 {
			break
		}
		baseH++
	}
	return baseH
}

// 测试
func Test_skip() {
	sk := MakeSkipLink()
	sk.AddElem("a", 1234)
	sk.AddElem("b", 124)
	sk.AddElem("c", 7234)
	sk.AddElem("d", 1264)
	sk.DelElem("c")
	sk.UpdateScoreElem("d", 9.0)
	sk.AddElem("tt1", -9.4)
	sk.AddElem("b1", 124)
	sk.AddElem("c1", 7234)
	sk.AddElem("d1", 1264)
}
