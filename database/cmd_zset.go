package database

import (
	_const "MyGoRedis/const"
	"MyGoRedis/datastruct/myzset"
	"MyGoRedis/interface/resp"
	"MyGoRedis/lib/utils"
	"MyGoRedis/resp/reply"
	"strconv"
	"strings"
)

// 实现命令

// 检查完成

// 含义：向有序集合添加一个或多个成员，或更新已存在成员的分数
// 用法：ZADD key score member [score member ...]
// 返回值：成功添加到有序集合的新成员数量，不包括已经存在但分数被更新的成员
func exec_ZSET_ZADD(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_ZSET_ZADD, args...))
	var typeZSet *myzset.ZSet
	it, ok := db.GetEntity(key)
	if !ok {
		typeZSet = myzset.MakeZSet()
	} else {
		typeZSet = _const.DataToZSET(it)
		if typeZSet == nil {
			typeZSet = myzset.MakeZSet()
		}
	}
	index := 1
	count := int32(0)
	for index < len(args) {
		score, err := utils.StringToFloat64(string(args[index]))
		if err != nil {
			break
		}
		index++
		if index >= len(args) {
			break
		}
		member := string(args[index])
		index++
		// 设置
		if 0 == typeZSet.Add(member, float32(score)) {
			typeZSet.Update(member, float32(score))
		} else {
			count++
		}
	}
	db.PutEntity(key, typeZSet)
	return reply.MakeIntReply(int64(count))
}

// 含义：获取有序集合的成员数量。
// 用法：ZCARD key
// 返回值：有序集合的基数（成员数量）。
func exec_ZSET_ZCARD(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeIntReply(0)
	}
	typeZSet := _const.DataToZSET(it)
	if typeZSet == nil {
		return reply.MakeUnknownErrReply()
	}
	return reply.MakeIntReply(int64(typeZSet.Size()))
}

// 含义：计算有序集合中指定分数范围内的成员数量。
// 用法：ZCOUNT key min max
// 返回值：指定分数范围内的成员数量。
func exec_ZSET_ZCOUNT(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	minScore, err := utils.StringToFloat64(string(args[1]))
	if err != nil {
		return reply.MakeIntReply(0)
	}
	maxScore, err := utils.StringToFloat64(string(args[2]))
	if err != nil {
		return reply.MakeIntReply(0)
	}
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeIntReply(0)
	}
	typeZSet := _const.DataToZSET(it)
	if typeZSet == nil {
		return reply.MakeIntReply(0)
	}
	return reply.MakeIntReply(int64(typeZSet.LenScoreRange(float32(minScore), float32(maxScore))))
}

// 含义：为有序集合的成员增加分数
// 用法：ZINCRBY key increment member
// 返回值：增加后的成员分数
func exec_ZSET_ZINCRBY(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_ZSET_ZINCRBY, args...))
	incr, err := utils.StringToFloat64(string(args[1]))
	if err != nil {
		return reply.MakeIntReply(0)
	}
	member := string(args[2])
	var typeZSet *myzset.ZSet
	it, ok := db.GetEntity(key)
	if !ok {
		typeZSet = myzset.MakeZSet()
	} else {
		typeZSet = _const.DataToZSET(it)
		if typeZSet == nil {
			typeZSet = myzset.MakeZSet()
		}
	}
	// 判断是否存在
	var newScore float32
	if typeZSet.IsExists(member) {
		newScore = typeZSet.Incr(member, float32(incr))
	} else {
		typeZSet.Add(member, float32(incr))
		newScore = float32(incr)
	}
	db.PutEntity(key, typeZSet)
	return reply.MakeBulkReply(utils.Float64ToByte(float64(newScore)))
}

// 含义：通过索引区间返回有序集合指定区间内的成员。
// 用法：ZRANGE key start stop [WITHSCORES]
// 返回值：指定区间内的成员列表。
func exec_ZSET_ZRANGE(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	var typeZSet *myzset.ZSet
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeZSet = _const.DataToZSET(it)
	if typeZSet == nil {
		return reply.MakeNullBulkReply()
	}
	start, err := strconv.Atoi(string(args[1]))
	if err != nil {
		return reply.MakeNullBulkReply()
	}
	stop, err := strconv.Atoi(string(args[2]))
	if err != nil {
		return reply.MakeNullBulkReply()
	}
	valtmps, err := typeZSet.RankRange(int32(start), int32(stop))
	if err != nil {
		return reply.MakeUnknownErrReply()
	}
	slic := make([]string, 0)
	if len(args) > 3 && strings.ToLower(string(args[3])) == "withscores" {
		for _, val := range valtmps {
			slic = append(slic, val.Value)
			slic = append(slic, utils.Float64ToString(float64(val.Score)))
		}
	} else {
		for _, val := range valtmps {
			slic = append(slic, val.Value)
		}
	}
	return reply.MakeMultiBulkReply(utils.ToCmdLine(slic...))
}

// 含义：通过分数范围返回有序集合的成员。
// 用法：ZRANGEBYSCORE key min max [WITHSCORES]
// 返回值：指定分数范围内的成员列表。
func exec_ZSET_ZRANGEBYSCORE(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	var typeZSet *myzset.ZSet
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeZSet = _const.DataToZSET(it)
	if typeZSet == nil {
		return reply.MakeNullBulkReply()
	}
	minScore, err := utils.StringToFloat64(string(args[1]))
	if err != nil {
		return reply.MakeNullBulkReply()
	}
	maxScore, err := utils.StringToFloat64(string(args[2]))
	if err != nil {
		return reply.MakeNullBulkReply()
	}
	valtmps, err := typeZSet.ScoreRange(float32(minScore), float32(maxScore))
	if err != nil {
		return reply.MakeUnknownErrReply()
	}
	slic := make([]string, 0)
	if len(args) > 3 && strings.ToLower(string(args[3])) == "withscores" {
		for _, val := range valtmps {
			slic = append(slic, val.Value)
			slic = append(slic, utils.Float64ToString(float64(val.Score)))
		}
	} else {
		for _, val := range valtmps {
			slic = append(slic, val.Value)
		}
	}
	return reply.MakeMultiBulkReply(utils.ToCmdLine(slic...))
}

// 含义：获取有序集合中指定成员的排名（从0开始）。
// 用法：ZRANK key member
// 返回值：指定成员的排名。
func exec_ZSET_ZRANK(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	member := string(args[1])
	var typeZSet *myzset.ZSet
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeZSet = _const.DataToZSET(it)
	if typeZSet == nil {
		return reply.MakeNullBulkReply()
	}
	rank, err := typeZSet.GetRank(member)
	if err != nil {
		return reply.MakeStandardErrReply(err.Error())
	}
	return reply.MakeIntReply(int64(rank))
}

// 含义：从有序集合中移除一个或多个成员。
// 用法：ZREM key member [member ...]
// 返回值：被成功移除的成员数量，不包括不存在于集合中的成员。
func exec_ZSET_ZREM(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_ZSET_ZREM, args...))
	var typeZSet *myzset.ZSet
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeIntReply(0)
	}
	typeZSet = _const.DataToZSET(it)
	if typeZSet == nil {
		return reply.MakeIntReply(0)
	}
	index := int32(1)
	count := int32(0)
	for index < int32(len(args)) {
		member := string(args[index])
		index++
		count += typeZSet.Remove(member)
	}
	return reply.MakeIntReply(int64(count))
}

// 含义：移除有序集合中给定排名范围内的所有成员。
// 用法：ZREMRANGEBYRANK key start stop
// 返回值：被移除的成员数量。
func exec_ZSET_ZREMRANGEBYRANK(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_ZSET_ZREMRANGEBYRANK, args...))
	var typeZSet *myzset.ZSet
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeIntReply(0)
	}
	typeZSet = _const.DataToZSET(it)
	if typeZSet == nil {
		return reply.MakeIntReply(0)
	}
	start, err := strconv.Atoi(string(args[1]))
	if err != nil {
		return reply.MakeStandardErrReply(err.Error())
	}
	stop, err := strconv.Atoi(string(args[2]))
	if err != nil {
		return reply.MakeStandardErrReply(err.Error())
	}
	count := typeZSet.RemoveRangeRank(int32(start), int32(stop))
	return reply.MakeIntReply(int64(count))
}

// 含义：移除有序集合中给定分数范围内的所有成员。
// 用法：ZREMRANGEBYSCORE key min max
// 返回值：被移除的成员数量。
func exec_ZSET_ZREMRANGEBYSCORE(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_ZSET_ZREMRANGEBYSCORE, args...))
	var typeZSet *myzset.ZSet
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeIntReply(0)
	}
	typeZSet = _const.DataToZSET(it)
	if typeZSet == nil {
		return reply.MakeIntReply(0)
	}
	minScore, err := utils.StringToFloat64(string(args[1]))
	if err != nil {
		return reply.MakeStandardErrReply(err.Error())
	}
	maxScore, err := utils.StringToFloat64(string(args[2]))
	if err != nil {
		return reply.MakeStandardErrReply(err.Error())
	}
	return reply.MakeIntReply(int64(typeZSet.RemoveRangeScore(float32(minScore), float32(maxScore))))
}

// 含义：返回有序集合中指定区间内的成员，按分数从高到低排序。
// 用法：ZREVRANGE key start stop [WITHSCORES]
// 返回值：指定区间内的成员列表。
func exec_ZSET_ZREVRANGE(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	var typeZSet *myzset.ZSet
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeZSet = _const.DataToZSET(it)
	if typeZSet == nil {
		return reply.MakeNullBulkReply()
	}
	start, err := strconv.Atoi(string(args[1]))
	if err != nil {
		return reply.MakeStandardErrReply(err.Error())
	}
	stop, err := strconv.Atoi(string(args[2]))
	if err != nil {
		return reply.MakeStandardErrReply(err.Error())
	}
	valtmps, err := typeZSet.RankRange(int32(start), int32(stop))
	if err != nil {
		return reply.MakeUnknownErrReply()
	}
	slic := make([]string, 0)
	if len(args) > 3 && strings.ToLower(string(args[3])) == "withscores" {
		for _, val := range valtmps {
			slic = append(slic, val.Value)
			slic = append(slic, utils.Float64ToString(float64(val.Score)))
		}
	} else {
		for _, val := range valtmps {
			slic = append(slic, val.Value)
		}
	}
	return reply.MakeMultiBulkReply(utils.ToCmdLine(slic...))
}

// 含义：返回有序集合中指定分数范围内的成员，按分数从高到低排序。
// 用法：ZREVRANGEBYSCORE key max min [WITHSCORES]
// 返回值：指定分数范围内的成员列表。
func exec_ZSET_ZREVRANGEBYSCORE(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	var typeZSet *myzset.ZSet
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeZSet = _const.DataToZSET(it)
	if typeZSet == nil {
		return reply.MakeNullBulkReply()
	}
	minScore, err := utils.StringToFloat64(string(args[1]))
	if err != nil {
		return reply.MakeStandardErrReply(err.Error())
	}
	maxScore, err := utils.StringToFloat64(string(args[2]))
	if err != nil {
		return reply.MakeStandardErrReply(err.Error())
	}
	valtmps, err := typeZSet.ScoreRange(float32(minScore), float32(maxScore))
	if err != nil {
		return reply.MakeUnknownErrReply()
	}
	slic := make([]string, 0)
	if len(args) > 3 && strings.ToLower(string(args[3])) == "withscores" {
		for _, val := range valtmps {
			slic = append(slic, val.Value)
			slic = append(slic, utils.Float64ToString(float64(val.Score)))
		}
	} else {
		for _, val := range valtmps {
			slic = append(slic, val.Value)
		}
	}
	return reply.MakeMultiBulkReply(utils.ToCmdLine(slic...))
}

// 含义：获取有序集合中指定成员的倒序排名（从0开始）。
// 用法：ZREVRANK key member
// 返回值：指定成员的倒序排名。
func exec_ZSET_ZREVRANK(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	member := string(args[1])
	var typeZSet *myzset.ZSet
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeZSet = _const.DataToZSET(it)
	if typeZSet == nil {
		return reply.MakeNullBulkReply()
	}
	rank, err := typeZSet.GetRank(member)
	if err != nil {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeIntReply(int64(typeZSet.Size() - 1 - rank))
}

// 含义：获取有序集合中指定成员的分数。
// 用法：ZSCORE key member
// 返回值：指定成员的分数。
func exec_ZSET_ZSCORE(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	member := string(args[1])
	var typeZSet *myzset.ZSet
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeZSet = _const.DataToZSET(it)
	if typeZSet == nil {
		return reply.MakeNullBulkReply()
	}
	score, err := typeZSet.GetScore(member)
	if err != nil {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeBulkReply(utils.Float64ToByte(float64(score)))
}

func init() {
	RegisterCommand(_const.CMD_ZSET_ZADD, exec_ZSET_ZADD, -4)
	RegisterCommand(_const.CMD_ZSET_ZCARD, exec_ZSET_ZCARD, 2)
	RegisterCommand(_const.CMD_ZSET_ZCOUNT, exec_ZSET_ZCOUNT, 4)
	RegisterCommand(_const.CMD_ZSET_ZINCRBY, exec_ZSET_ZINCRBY, 4)
	RegisterCommand(_const.CMD_ZSET_ZRANGE, exec_ZSET_ZRANGE, -4)
	RegisterCommand(_const.CMD_ZSET_ZRANGEBYSCORE, exec_ZSET_ZRANGEBYSCORE, -4)
	RegisterCommand(_const.CMD_ZSET_ZRANK, exec_ZSET_ZRANK, 3)
	RegisterCommand(_const.CMD_ZSET_ZREM, exec_ZSET_ZREM, -3)
	RegisterCommand(_const.CMD_ZSET_ZREMRANGEBYRANK, exec_ZSET_ZREMRANGEBYRANK, 4)
	RegisterCommand(_const.CMD_ZSET_ZREMRANGEBYSCORE, exec_ZSET_ZREMRANGEBYSCORE, 4)
	RegisterCommand(_const.CMD_ZSET_ZREVRANGE, exec_ZSET_ZREVRANGE, -4)
	RegisterCommand(_const.CMD_ZSET_ZREVRANGEBYSCORE, exec_ZSET_ZREVRANGEBYSCORE, -4)
	RegisterCommand(_const.CMD_ZSET_ZREVRANK, exec_ZSET_ZREVRANK, 3)
	RegisterCommand(_const.CMD_ZSET_ZSCORE, exec_ZSET_ZSCORE, 3)
}
