package myzset

import (
	"MyGoRedis/lib/basedatastruct"
)

type ZSet struct {
	data *basedatastruct.BaseSkipLink
}

func MakeZSet() *ZSet {
	return &ZSet{data: basedatastruct.MakeSkipLink()}
}

// 向有序集合添加一个成员
func (zset *ZSet) Add(value string, score float32) error {
	return zset.data.AddElem(value, score)
}

// 更新已存在成员的分数
func (zset *ZSet) Update(value string, score float32) error {
	return zset.data.UpdateScoreElem(value, score)
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
func (zset *ZSet) Incr(value string, incr float32) error {
	return zset.data.Incr(value, incr)
}

// 通过索引区间返回有序集合指定区间内的成员
func (zset *ZSet) RankRange(indexStart, indexEnd int32) ([]*basedatastruct.Valtmp, error) {
	return zset.data.GetRankElem(indexStart, indexEnd)
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
func (zset *ZSet) Remove(value string) error {
	return zset.data.DelElem(value)
}

// 移除有序集合中给定的排名区间的所有成员
func (zset *ZSet) RemoveRangeRank(rankStart, rankEnd int32) {
	zset.data.RemoveRangeRank(rankStart, rankEnd)
}

// 移除有序集合中给定的分数区间的所有成员
func (zset *ZSet) RemoveRangeScore(scoreStart, scoreEnd float32) {
	zset.data.RemoveRangeScore(scoreStart, scoreEnd)
}

// 返回有序集中，成员的分数值
func (zset *ZSet) GetScore(value string) (float32, error) {
	return zset.data.GetScore(value)
}
