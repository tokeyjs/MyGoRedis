package database

import (
	_const "MyGoRedis/const"
	"MyGoRedis/datastruct/myhash"
	"MyGoRedis/interface/resp"
	"MyGoRedis/lib/utils"
	"MyGoRedis/resp/reply"
)

// 实现命令

// 检查完成

// 含义：删除哈希表中一个或多个字段。
// 用法：HDEL key field1 [field2 ...]
// 返回值：被成功移除的字段的数量，不包括被忽略的字段。
func exec_HASH_HDEL(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_HASH_HDEL, args...))
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeIntReply(0)
	}
	typeHash := _const.DataToHash(it)
	if typeHash == nil {
		// 错误
		return reply.MakeUnknownErrReply()
	}
	index := 1
	count := int32(0)
	for index < len(args) {
		field := string(args[index])
		count += typeHash.DelFiled(field)
		index++
	}
	return reply.MakeIntReply(int64(count))
}

// 含义：检查哈希表中是否存在指定字段。
// 用法：HEXISTS key field
// 返回值：如果字段存在于哈希表中，则返回1；如果字段不存在或者哈希表不存在，则返回0。
func exec_HASH_HEXISTS(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	field := string(args[1])
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeIntReply(0)
	}
	typeHash := _const.DataToHash(it)
	if typeHash == nil {
		// 错误
		return reply.MakeUnknownErrReply()
	}
	if typeHash.IsExists(field) {
		return reply.MakeIntReply(1)
	}
	return reply.MakeIntReply(0)
}

// 含义：获取哈希表中指定字段的值。
// 用法：HGET key field
// 返回值：如果字段存在，则返回字段的值；如果字段不存在或者哈希表不存在，则返回nil。
func exec_HASH_HGET(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	field := string(args[1])
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeHash := _const.DataToHash(it)
	if typeHash == nil {
		return reply.MakeNullBulkReply()
	}
	val, ok := typeHash.Get(field)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeBulkReply([]byte(val))
}

// 含义：获取哈希表中所有字段和值。
// 用法：HGETALL key
// 返回值：返回一个包含所有字段和值的列表。
func exec_HASH_HGETALL(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeIntReply(0)
	}
	typeHash := _const.DataToHash(it)
	if typeHash == nil {
		return reply.MakeUnknownErrReply()
	}
	slic := typeHash.GetAllKV()
	strSlic := make([]string, 0, len(slic)*2)
	for _, v := range slic {
		strSlic = append(strSlic, v.Field)
		strSlic = append(strSlic, v.Value)
	}
	return reply.MakeMultiBulkReply(utils.ToCmdLine(strSlic...))
}

// 含义：为哈希表中指定字段的值加上增量。
// 用法：HINCRBY key field increment
// 返回值：增加后的字段值。
func exec_HASH_HINCRBY(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	field := string(args[1])
	incr, err := utils.StringToFloat64(string(args[2]))
	if err != nil {
		return reply.MakeStandardErrReply("increment is error")
	}
	db.aofAdd(utils.ToCmdLine2(_const.CMD_HASH_HINCRBY, args...))
	it, ok := db.GetEntity(key)
	if !ok {
		// 新建一个
		myh := myhash.MakeHash()
		myh.Set(field, utils.Float64ToString(incr))
		db.PutEntity(key, myh)
		return reply.MakeBulkReply(utils.Float64ToByte(incr))
	}
	typeHash := _const.DataToHash(it)
	if typeHash == nil {
		// 错误
		return reply.MakeUnknownErrReply()
	}
	ret, err := typeHash.Incr(field, incr)
	if err != nil {
		return reply.MakeStandardErrReply(err.Error())
	}
	return reply.MakeBulkReply(utils.Float64ToByte(ret))
}

// 含义：为哈希表中指定字段的值加上浮点数增量。
// 用法：HINCRBYFLOAT key field increment
// 返回值：增加后的字段值。
func exec_HASH_HINCRBYFLOAT(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	field := string(args[1])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_HASH_HINCRBYFLOAT, args...))
	incr, err := utils.StringToFloat64(string(args[2]))
	if err != nil {
		return reply.MakeStandardErrReply("increment is error")
	}
	it, ok := db.GetEntity(key)
	if !ok {
		// 新建一个
		myh := myhash.MakeHash()
		myh.Set(field, utils.Float64ToString(incr))
		db.PutEntity(key, myh)
		return reply.MakeBulkReply(utils.Float64ToByte(incr))
	}
	typeHash := _const.DataToHash(it)
	if typeHash == nil {
		return reply.MakeUnknownErrReply()
	}
	ret, err := typeHash.Incr(field, incr)
	if err != nil {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeBulkReply(utils.Float64ToByte(ret))
}

// 含义：获取哈希表中所有字段的列表。
// 用法：HKEYS key
// 返回值：返回一个包含所有字段的列表。
func exec_HASH_HKEYS(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeHash := _const.DataToHash(it)
	if typeHash == nil {
		return reply.MakeUnknownErrReply()
	}
	fields := typeHash.GetAllField()
	if len(fields) == 0 {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeMultiBulkReply(utils.ToCmdLine(fields...))
}

// 含义：获取哈希表中字段的数量。
// 用法：HLEN key
// 返回值：哈希表中字段的数量。
func exec_HASH_HLEN(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeIntReply(0)
	}
	typeHash := _const.DataToHash(it)
	if typeHash == nil {
		// 错误
		return reply.MakeUnknownErrReply()
	}
	return reply.MakeIntReply(int64(typeHash.Len()))
}

// 含义：获取哈希表中一个或多个字段的值。
// 用法：HMGET key field1 [field2 ...]
// 返回值：一个包含指定字段值的列表。
func exec_HASH_HMGET(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeHash := _const.DataToHash(it)
	if typeHash == nil {
		return reply.MakeUnknownErrReply()
	}
	slic := make([]string, 0)
	index := 1
	for index < len(args) {
		field := string(args[index])
		val, ok := typeHash.Get(field)
		if !ok {
			slic = append(slic, "nil")
		} else {
			slic = append(slic, val)
		}
		index++
	}
	return reply.MakeMultiBulkReply(utils.ToCmdLine(slic...))
}

// 含义：同时设置哈希表中的多个字段值。
// 用法：HMSET key field1 value1 [field2 value2 ...]
// 返回值：始终返回OK。
func exec_HASH_HMSET(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_HASH_HMSET, args...))
	var typeHash *myhash.Hash
	it, ok := db.GetEntity(key)
	if !ok {
		typeHash = myhash.MakeHash()
	} else {
		typeHash = _const.DataToHash(it)
		if typeHash == nil {
			typeHash = myhash.MakeHash()
		}
	}
	index := 1
	for index < len(args) {
		filed := string(args[index])
		index++
		if index >= len(args) {
			break
		}
		value := string(args[index])
		index++
		typeHash.Set(filed, value)
	}
	db.PutEntity(key, typeHash)
	return reply.MakeOkReply()
}

// 含义：设置哈希表中的一个字段值。
// 用法：HSET key field value
// 返回值：如果字段是一个新字段并成功设置了值，则返回1；如果字段已经存在，则更新值，返回0。
func exec_HASH_HSET(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	field := string(args[1])
	value := string(args[2])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_HASH_HSET, args...))
	it, ok := db.GetEntity(key)
	if !ok {
		// 新建一个
		myh := myhash.MakeHash()
		myh.Set(field, value)
		db.PutEntity(key, myh)
		return reply.MakeIntReply(1)
	}
	typeHash := _const.DataToHash(it)
	if typeHash == nil {
		return reply.MakeUnknownErrReply()
	}
	return reply.MakeIntReply(int64(typeHash.Set(field, value)))
}

// 含义：只在哈希表中的字段不存在时，设置字段的值。
// 用法：HSETNX key field value
// 返回值：如果字段是一个新字段并成功设置了值，则返回1；如果字段已经存在，则不执行任何操作，返回0。
func exec_HASH_HSETNX(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	field := string(args[1])
	value := string(args[2])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_HASH_HSETNX, args...))
	it, ok := db.GetEntity(key)
	if !ok {
		// 新建一个
		myh := myhash.MakeHash()
		myh.Set(field, value)
		db.PutEntity(key, myh)
		return reply.MakeIntReply(1)
	}
	typeHash := _const.DataToHash(it)
	if typeHash == nil {
		return reply.MakeUnknownErrReply()
	}
	if typeHash.IsExists(field) {
		return reply.MakeIntReply(0)
	}
	return reply.MakeIntReply(int64(typeHash.Set(field, value)))
}

// 含义：获取哈希表中所有值的列表。
// 用法：HVALS key
// 返回值：返回一个包含所有值的列表。
func exec_HASH_HVALS(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeHash := _const.DataToHash(it)
	if typeHash == nil {
		return reply.MakeUnknownErrReply()
	}
	values := typeHash.GetAllValue()
	if len(values) == 0 {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeMultiBulkReply(utils.ToCmdLine(values...))
}

func init() {
	RegisterCommand(_const.CMD_HASH_HDEL, exec_HASH_HDEL, -3)
	RegisterCommand(_const.CMD_HASH_HEXISTS, exec_HASH_HEXISTS, 3)
	RegisterCommand(_const.CMD_HASH_HGET, exec_HASH_HGET, 3)
	RegisterCommand(_const.CMD_HASH_HGETALL, exec_HASH_HGETALL, 2)
	RegisterCommand(_const.CMD_HASH_HINCRBY, exec_HASH_HINCRBY, 4)
	RegisterCommand(_const.CMD_HASH_HINCRBYFLOAT, exec_HASH_HINCRBYFLOAT, 4)
	RegisterCommand(_const.CMD_HASH_HKEYS, exec_HASH_HKEYS, 2)
	RegisterCommand(_const.CMD_HASH_HLEN, exec_HASH_HLEN, 2)
	RegisterCommand(_const.CMD_HASH_HMGET, exec_HASH_HMGET, -3)
	RegisterCommand(_const.CMD_HASH_HMSET, exec_HASH_HMSET, -4)
	RegisterCommand(_const.CMD_HASH_HSET, exec_HASH_HSET, 4)
	RegisterCommand(_const.CMD_HASH_HSETNX, exec_HASH_HSETNX, 4)
	RegisterCommand(_const.CMD_HASH_HVALS, exec_HASH_HVALS, 2)

}
