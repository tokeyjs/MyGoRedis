package database

import (
	_const "MyGoRedis/const"
	"MyGoRedis/datastruct/myset"
	"MyGoRedis/interface/resp"
	"MyGoRedis/lib/utils"
	"MyGoRedis/resp/reply"
	"strconv"
)

// 实现命令

// 检查完成

// 含义：向集合中添加一个或多个成员。
// 用法：SADD key member1 [member2 ...]
// 返回值：添加到集合中的新元素的数量，不包括已经存在于集合中的元素。
func exec_SET_SADD(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_SET_SADD, args...))
	var typeSet *myset.Set
	it, ok := db.GetEntity(key)
	if !ok {
		typeSet = myset.MakeSet()
	} else {
		typeSet = _const.DataToSET(it)
		if typeSet == nil {
			typeSet = myset.MakeSet()
		}
	}
	index := 1
	count := 0
	for index < len(args) {
		value := string(args[index])
		if !typeSet.IsExists(value) {
			count++
			typeSet.Add(value)
		}
		index++
	}
	db.PutEntity(key, typeSet)
	return reply.MakeIntReply(int64(count))
}

// 含义：获取集合中的成员数量。
// 用法：SCARD key
// 返回值：集合的基数（元素数量）。
func exec_SET_SCARD(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeIntReply(0)
	}
	typeSet := _const.DataToSET(it)
	if typeSet == nil {
		return reply.MakeUnknownErrReply()
	}
	return reply.MakeIntReply(int64(typeSet.Size()))
}

// 含义：返回给定所有集合的交集。
// 用法：SINTER key [key ...]
// 返回值：包含交集成员的列表。
func exec_SET_SINTER(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：将给定所有集合的交集存储到指定的目标集合中。
// 用法：SINTERSTORE destination key [key ...]
// 返回值：存储到目标集合的元素数量。
func exec_SET_SINTERSTORE(db *DB, args [][]byte) resp.Reply {
	// todo
	db.aofAdd(utils.ToCmdLine2(_const.CMD_SET_SINTERSTORE, args...))
	return reply.MakeUnknownErrReply()
}

// 含义：判断指定成员是否存在于集合中。
// 用法：SISMEMBER key member
// 返回值：如果成员存在于集合中，则返回1；如果成员不存在或集合不存在，则返回0。
func exec_SET_SISMEMBER(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := string(args[1])
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeIntReply(0)
	}
	typeSet := _const.DataToSET(it)
	if typeSet == nil {
		return reply.MakeUnknownErrReply()
	}
	if typeSet.IsExists(value) {
		return reply.MakeIntReply(1)
	}
	return reply.MakeIntReply(0)
}

// 含义：返回集合中的所有成员。
// 用法：SMEMBERS key
// 返回值：包含所有成员的列表。
func exec_SET_SMEMBERS(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeSet := _const.DataToSET(it)
	if typeSet == nil {
		// 错误
		return reply.MakeUnknownErrReply()
	}
	slic := typeSet.GetAll()
	return reply.MakeMultiBulkReply(utils.ToCmdLine(slic...))
}

// 含义：将指定成员从源集合移动到目标集合。
// 用法：SMOVE source destination member
// 返回值：如果成员成功移动，则返回1；如果成员不存在于源集合中，则返回0。
func exec_SET_SMOVE(db *DB, args [][]byte) resp.Reply {
	// todo
	db.aofAdd(utils.ToCmdLine2(_const.CMD_SET_SMOVE, args...))
	return reply.MakeUnknownErrReply()
}

// 含义：随机移除并返回集合中的一个成员。
// 用法：SPOP key
// 返回值：移除的成员。
func exec_SET_SPOP(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_SET_SPOP, args...))
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeStandardErrReply("no key")
	}
	typeSet := _const.DataToSET(it)
	if typeSet == nil {
		return reply.MakeUnknownErrReply()
	}
	str := typeSet.GetRandom(1)
	if len(str) <= 0 {
		return reply.MakeNullBulkReply()
	}
	typeSet.Remove(str[0])
	return reply.MakeBulkReply([]byte(str[0]))
}

// 含义：随机获取集合中的一个或多个成员。
// 用法：SRANDMEMBER key [count]
// 返回值：返回一个或多个随机成员，不移除成员。
func exec_SET_SRANDMEMBER(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	var err error
	count := 1
	if len(args) == 2 {
		count, err = strconv.Atoi(string(args[1]))
		if err != nil {
			return reply.MakeStandardErrReply("argment 'count' is not int")
		}
	}
	if len(args) > 1 {
		count, _ = strconv.Atoi(string(args[1]))
	}
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeIntReply(0)
	}
	typeSet := _const.DataToSET(it)
	if typeSet == nil {
		// 错误
		return reply.MakeUnknownErrReply()
	}
	slic := typeSet.GetRandom(int32(count))
	return reply.MakeMultiBulkReply(utils.ToCmdLine(slic...))
}

// 含义：移除集合中一个或多个成员。
// 用法：SREM key member1 [member2 ...]
// 返回值：移除的成员数量，不包括不存在于集合中的成员。
func exec_SET_SREM(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_SET_SREM, args...))
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeIntReply(0)
	}
	typeSet := _const.DataToSET(it)
	if typeSet == nil {
		return reply.MakeUnknownErrReply()
	}
	index := 1
	count := int32(0)
	for index < len(args) {
		member := string(args[index])
		count += typeSet.Remove(member)
		index++
	}
	return reply.MakeIntReply(int64(count))
}

func init() {
	RegisterCommand(_const.CMD_SET_SADD, exec_SET_SADD, -3)
	RegisterCommand(_const.CMD_SET_SCARD, exec_SET_SCARD, 2)
	RegisterCommand(_const.CMD_SET_SINTER, exec_SET_SINTER, -2)
	RegisterCommand(_const.CMD_SET_SINTERSTORE, exec_SET_SINTERSTORE, -3)
	RegisterCommand(_const.CMD_SET_SISMEMBER, exec_SET_SISMEMBER, 3)
	RegisterCommand(_const.CMD_SET_SMEMBERS, exec_SET_SMEMBERS, 2)
	RegisterCommand(_const.CMD_SET_SMOVE, exec_SET_SMOVE, 4)
	RegisterCommand(_const.CMD_SET_SPOP, exec_SET_SPOP, 2)
	RegisterCommand(_const.CMD_SET_SRANDMEMBER, exec_SET_SRANDMEMBER, -2)
	RegisterCommand(_const.CMD_SET_SREM, exec_SET_SREM, -3)
}
