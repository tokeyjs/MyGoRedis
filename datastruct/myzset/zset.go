package myzset

import (
	"MyGoRedis/lib/basedatastruct"
	"errors"
	"sync"
)

type ZSet struct {
	data *basedatastruct.BaseSkipLink
	set  sync.Map // map[string]float
}

func MakeZSet() *ZSet {
	return &ZSet{data: basedatastruct.MakeSkipLink()}
}

// 向有序集合添加一个成员 不存在添加成功1 存在添加0
func (zset *ZSet) Add(value string, score float32) int32 {
	if zset.IsExists(value) {
		return 0
	}
	zset.set.Store(value, score)
	zset.data.AddElem(value, score)
	return 1
}

func (zset *ZSet) IsExists(value string) bool {
	_, ok := zset.set.Load(value)
	return ok
}

// 更新已存在成员的分数
func (zset *ZSet) Update(value string, score float32) int32 {
	if zset.IsExists(value) {
		srcScore, _ := zset.set.Load(value)
		err := zset.data.UpdateScoreElem(value, srcScore.(float32), score)
		if err != nil {
			return 0
		}
		zset.set.Store(value, score)
		return 1
	}
	// 不存在不进行更新
	return 0
}

// 获取有序集合的成员数
func (zset *ZSet) Size() int32 {
	return zset.data.Size()
}

// 指定区间分数的成员数
func (zset *ZSet) LenScoreRange(scoreMin, scoreMax float32) int32 {
	return zset.data.CountRangeScore(scoreMin, scoreMax)
}

// 对指定成员的分数加上增量 increment
func (zset *ZSet) Incr(value string, incr float32) float32 {
	if !zset.IsExists(value) {
		return 0
	}
	score, _ := zset.set.Load(value)
	err := zset.data.Incr(value, score.(float32), incr)
	if err != nil {
		return 0
	}
	zset.set.Store(value, score.(float32)+incr)
	return score.(float32) + incr
}

// 通过索引区间返回有序集合指定区间内的成员
func (zset *ZSet) RankRange(indexStart, indexEnd int32) ([]*basedatastruct.Valtmp, error) {
	return zset.data.GetRankRangeElem(indexStart, indexEnd)
}

// 通过分数返回有序集合指定区间内的成员
func (zset *ZSet) ScoreRange(scoreMin, scoreMax float32) ([]*basedatastruct.Valtmp, error) {
	return zset.data.GetScoreRangeElem(scoreMin, scoreMax)
}

// 返回有序集合中指定成员的索引
func (zset *ZSet) GetRank(value string) (int32, error) {
	return zset.data.GetIndexElem(value)
}

// 移除有序集合中的一个成员
func (zset *ZSet) Remove(value string) int32 {
	if !zset.IsExists(value) {
		return 0
	}
	score, _ := zset.set.Load(value)
	ok := zset.data.DelElem(value, score.(float32))
	if !ok {
		return 0
	}
	zset.set.Delete(value)
	return 1
}

// 移除有序集合中给定的排名区间的所有成员
func (zset *ZSet) RemoveRangeRank(rankStart, rankEnd int32) int32 {
	return zset.data.RemoveRangeRank(rankStart, rankEnd)
}

// 移除有序集合中给定的分数区间的所有成员
func (zset *ZSet) RemoveRangeScore(scoreStart, scoreEnd float32) int32 {
	return zset.data.RemoveRangeScore(scoreStart, scoreEnd)
}

// 返回有序集中，成员的分数值
func (zset *ZSet) GetScore(value string) (float32, error) {
	if zset.IsExists(value) {
		val, _ := zset.set.Load(value)
		return val.(float32), nil
	}
	return 0, errors.New("not found value")
}
