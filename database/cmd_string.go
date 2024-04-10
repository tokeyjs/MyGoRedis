package database

import (
	_const "MyGoRedis/const"
	"MyGoRedis/interface/resp"
	"MyGoRedis/resp/reply"
)

// 实现命令

// 含义：将指定值追加到键的当前值的末尾。
// 用法：APPEND key value
// 返回值：追加后的字符串长度。
func exec_STRING_APPEND(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeUnknownErrReply()
}

// 含义：将键的值减1。
// 用法：DECR key
// 返回值：减少后的值。
func exec_STRING_DECR(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeUnknownErrReply()
}

// 含义：将键的值减去指定的整数值。
// 用法：DECRBY key decrement
// 返回值：减少后的值。
func exec_STRING_DECRBY(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeUnknownErrReply()
}

// 含义：获取指定键的值。
// 用法：GET key
// 返回值：指定键的值，如果键不存在则返回nil。
func exec_STRING_GET(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeUnknownErrReply()
}

// 含义：获取指定键值的子字符串。
// 用法：GETRANGE key start end
// 返回值：子字符串。
func exec_STRING_GETRANGE(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeUnknownErrReply()
}

// 含义：设置指定键的值，并返回原来的值。
// 用法：GETSET key value
// 返回值：原来的值。
func exec_STRING_GETSET(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeUnknownErrReply()
}

// 含义：将键的值加1。
// 用法：INCR key
// 返回值：增加后的值。
func exec_STRING_INCR(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeUnknownErrReply()
}

// 含义：将键的值加上指定的整数值。
// 用法：INCRBY key increment
// 返回值：增加后的值。
func exec_STRING_INCRBY(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeUnknownErrReply()
}

// 含义：将键的值加上指定的浮点数值。
// 用法：INCRBYFLOAT key increment
// 返回值：增加后的值。
func exec_STRING_INCRBYFLOAT(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeUnknownErrReply()
}

// 含义：获取一个或多个键的值。
// 用法：MGET key [key ...]
// 返回值：包含指定键值的数组。
func exec_STRING_MGET(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeUnknownErrReply()
}

// 含义：设置多个键的值。
// 用法：MSET key value [key value ...]
// 返回值：始终返回OK。
func exec_STRING_MSET(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeUnknownErrReply()
}

// 含义：设置多个键的值，当且仅当所有指定的键都不存在时才执行设置操作。
// 用法：MSETNX key value [key value ...]
// 返回值：若所有键设置成功，则返回1；若至少一个键已存在，则返回0。
func exec_STRING_MSETNX(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeUnknownErrReply()
}

// 含义：设置键的值并指定过期时间（以毫秒为单位）。
// 用法：PSETEX key milliseconds value
// 返回值：始终返回OK。
func exec_STRING_PSETEX(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeUnknownErrReply()
}

// 含义：设置键值。
// 用法：SET key value
// 返回值：OK
func exec_STRING_SET(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeUnknownErrReply()
}

// 含义：设置键的值并指定过期时间（以秒为单位）。
// 用法：SETEX key seconds value
// 返回值：始终返回OK。
func exec_STRING_SETEX(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeUnknownErrReply()
}

// 含义：设置键的值，仅当键不存在时才执行设置操作。
// 用法：SETNX key value
// 返回值：若设置成功，则返回1；若键已存在，则返回0。
func exec_STRING_SETNX(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeUnknownErrReply()
}

// 含义：获取指定键值的长度。
// 用法：STRLEN key
// 返回值：字符串的长度。
func exec_STRING_STRLEN(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeUnknownErrReply()
}

func init() {
	RegisterCommand(_const.CMD_STRING_APPEND, exec_STRING_APPEND, 3)
	RegisterCommand(_const.CMD_STRING_DECR, exec_STRING_DECR, 2)
	RegisterCommand(_const.CMD_STRING_DECRBY, exec_STRING_DECRBY, 3)
	RegisterCommand(_const.CMD_STRING_GET, exec_STRING_GET, 2)
	RegisterCommand(_const.CMD_STRING_GETRANGE, exec_STRING_GETRANGE, 4)
	RegisterCommand(_const.CMD_STRING_GETSET, exec_STRING_GETSET, 3)
	RegisterCommand(_const.CMD_STRING_INCR, exec_STRING_INCR, 2)
	RegisterCommand(_const.CMD_STRING_INCRBY, exec_STRING_INCRBY, 3)
	RegisterCommand(_const.CMD_STRING_INCRBYFLOAT, exec_STRING_INCRBYFLOAT, 3)
	RegisterCommand(_const.CMD_STRING_MGET, exec_STRING_MGET, -2)
	RegisterCommand(_const.CMD_STRING_MSET, exec_STRING_MSET, -3)
	RegisterCommand(_const.CMD_STRING_MSETNX, exec_STRING_MSETNX, -3)
	RegisterCommand(_const.CMD_STRING_PSETEX, exec_STRING_PSETEX, 4)
	RegisterCommand(_const.CMD_STRING_SET, exec_STRING_SET, 3)
	RegisterCommand(_const.CMD_STRING_SETEX, exec_STRING_SETEX, 4)
	RegisterCommand(_const.CMD_STRING_SETNX, exec_STRING_SETNX, 3)
	RegisterCommand(_const.CMD_STRING_STRLEN, exec_STRING_STRLEN, 2)
}
