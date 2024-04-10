package database

import (
	_const "MyGoRedis/const"
	"MyGoRedis/interface/resp"
	"MyGoRedis/resp/reply"
)

// 实现命令

// 含义：删除哈希表中一个或多个字段。
// 用法：HDEL key field1 [field2 ...]
// 返回值：被成功移除的字段的数量，不包括被忽略的字段。
func exec_HASH_HDEL(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：检查哈希表中是否存在指定字段。
// 用法：HEXISTS key field
// 返回值：如果字段存在于哈希表中，则返回1；如果字段不存在或者哈希表不存在，则返回0。
func exec_HASH_HEXISTS(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：获取哈希表中指定字段的值。
// 用法：HGET key field
// 返回值：如果字段存在，则返回字段的值；如果字段不存在或者哈希表不存在，则返回nil。
func exec_HASH_HGET(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：获取哈希表中所有字段和值。
// 用法：HGETALL key
// 返回值：返回一个包含所有字段和值的列表。
func exec_HASH_HGETALL(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：为哈希表中指定字段的值加上增量。
// 用法：HINCRBY key field increment
// 返回值：增加后的字段值。
func exec_HASH_HINCRBY(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：为哈希表中指定字段的值加上浮点数增量。
// 用法：HINCRBYFLOAT key field increment
// 返回值：增加后的字段值。
func exec_HASH_HINCRBYFLOAT(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：获取哈希表中所有字段的列表。
// 用法：HKEYS key
// 返回值：返回一个包含所有字段的列表。
func exec_HASH_HKEYS(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：获取哈希表中字段的数量。
// 用法：HLEN key
// 返回值：哈希表中字段的数量。
func exec_HASH_HLEN(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：获取哈希表中一个或多个字段的值。
// 用法：HMGET key field1 [field2 ...]
// 返回值：一个包含指定字段值的列表。
func exec_HASH_HMGET(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：同时设置哈希表中的多个字段值。
// 用法：HMSET key field1 value1 [field2 value2 ...]
// 返回值：始终返回OK。
func exec_HASH_HMSET(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：设置哈希表中的一个字段值。
// 用法：HSET key field value
// 返回值：如果字段是一个新字段并成功设置了值，则返回1；如果字段已经存在，则更新值，返回0。
func exec_HASH_HSET(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：只在哈希表中的字段不存在时，设置字段的值。
// 用法：HSETNX key field value
// 返回值：如果字段是一个新字段并成功设置了值，则返回1；如果字段已经存在，则不执行任何操作，返回0。
func exec_HASH_HSETNX(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：获取哈希表中所有值的列表。
// 用法：HVALS key
// 返回值：返回一个包含所有值的列表。
func exec_HASH_HVALS(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
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
