package database

import (
	_const "MyGoRedis/const"
	"MyGoRedis/interface/resp"
	"MyGoRedis/lib/utils"
	"MyGoRedis/lib/wildcard"
	"MyGoRedis/resp/reply"
)

// 检查完成

// ===实现命令===

// 含义：删除给定的一个或多个键。
// 用法：DEL key [key ...]
// 返回值：被删除键的数量。
func exec_KEY_DEL(db *DB, args [][]byte) resp.Reply {
	count := 0
	for _, arg := range args {
		count += db.Remove(string(arg))
	}
	db.aofAdd(utils.ToCmdLine2(_const.CMD_KEY_DEL, args...))
	return reply.MakeIntReply(int64(count))
}

// 含义：检查给定键是否存在。
// 用法：EXISTS key
// 返回值：若键存在，则返回1，否则返回0。
func exec_KEY_EXISTS(db *DB, args [][]byte) resp.Reply {
	if db.IsExists(string(args[0])) {
		return reply.MakeIntReply(1)
	}
	return reply.MakeIntReply(0)
}

// 含义：设置键的过期时间，单位为秒。
// 用法：EXPIRE key seconds
// 返回值：若设置成功，则返回1，否则返回0。
func exec_KEY_EXPIRE(db *DB, args [][]byte) resp.Reply {
	// todo
	db.aofAdd(utils.ToCmdLine2(_const.CMD_KEY_EXPIRE, args...))
	return reply.MakeUnknownErrReply()
}

// 含义：设置键在指定的Unix时间戳（秒级）过期。
// 用法：EXPIREAT key timestamp
// 返回值：若设置成功，则返回1，否则返回0。
func exec_KEY_EXPIREAT(db *DB, args [][]byte) resp.Reply {
	// todo
	db.aofAdd(utils.ToCmdLine2(_const.CMD_KEY_EXPIREAT, args...))
	return reply.MakeUnknownErrReply()
}

// 含义：查找所有符合给定模式的键。
// 用法：KEYS pattern
// 返回值：符合模式的键组成的列表。
func exec_KEY_KEYS(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.data.Range(func(key any, val any) bool {
		if pattern.IsMatch(key.(string)) {
			result = append(result, []byte(key.(string)))
		}
		return true
	})
	return reply.MakeMultiBulkReply(result)
}

// 含义：将指定键移到另一个数据库。
// 用法：MOVE key db
// 返回值：若移动成功，则返回1，否则返回0。
func exec_KEY_MOVE(db *DB, args [][]byte) resp.Reply {
	//todo
	db.aofAdd(utils.ToCmdLine2(_const.CMD_KEY_MOVE, args...))
	return reply.MakeIntReply(0)
}

// 含义：移除键的过期时间，使其永不过期。
// 用法：PERSIST key
// 返回值：若键成功移除过期时间，则返回1，否则返回0。
func exec_KEY_PERSIST(db *DB, args [][]byte) resp.Reply {
	// todo
	db.aofAdd(utils.ToCmdLine2(_const.CMD_KEY_PERSIST, args...))
	return reply.MakeIntReply(0)
}

// 含义：设置键的过期时间，单位为毫秒。
// 用法：PEXPIRE key milliseconds
// 返回值：若设置成功，则返回1，否则返回0。
func exec_KEY_PEXPIRE(db *DB, args [][]byte) resp.Reply {
	// todo
	db.aofAdd(utils.ToCmdLine2(_const.CMD_KEY_PEXPIRE, args...))
	return reply.MakeIntReply(0)
}

// 含义：设置键在指定的Unix时间戳（毫秒级）过期。
// 用法：PEXPIREAT key milliseconds-timestamp
// 返回值：若设置成功，则返回1，否则返回0。
func exec_KEY_PEXPIREAT(db *DB, args [][]byte) resp.Reply {
	// todo
	db.aofAdd(utils.ToCmdLine2(_const.CMD_KEY_PEXPIREAT, args...))
	return reply.MakeUnknownErrReply()
}

// 含义：获取键的剩余过期时间，单位为毫秒。
// 用法：PTTL key
// 返回值：若键存在且有剩余过期时间，则返回剩余过期时间，若键不存在或不过期，则返回-1。
func exec_KEY_PTTL(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：从当前数据库中随机返回一个键。
// 用法：RANDOMKEY
// 返回值：随机键的名字，如果数据库为空，则返回nil。
func exec_KEY_RANDOMKEY(db *DB, args [][]byte) resp.Reply {
	sl := db.RandomKey(1)
	if len(sl) <= 0 {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeBulkReply([]byte(sl[0]))
}

// 含义：将键重命名为新键名
// 用法：RENAME key newkey
// 返回值：若重命名成功，则返回OK
func exec_KEY_RENAME(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])
	entity, exists := db.GetEntity(src)
	if !exists {
		return reply.MakeStandardErrReply("no such key")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)
	db.aofAdd(utils.ToCmdLine2(_const.CMD_KEY_RENAME, args...))
	return reply.MakeOkReply()
}

// 含义：仅当新键不存在时，将键重命名为新键名。
// 用法：RENAMENX key newkey
// 返回值：若重命名成功，则返回1，否则返回0。
func exec_KEY_RENAMENX(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])
	if _, ok := db.GetEntity(dest); ok {
		return reply.MakeIntReply(0)
	}
	entity, exists := db.GetEntity(src)
	if !exists {
		return reply.MakeStandardErrReply("no such key")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)
	db.aofAdd(utils.ToCmdLine2(_const.CMD_KEY_RENAMENX, args...))
	return reply.MakeIntReply(1)
}

// 含义：获取键的剩余过期时间，单位为秒。
// 用法：TTL key
// 返回值：若键存在且有剩余过期时间，则返回剩余过期时间，若键不存在或不过期，则返回-1。
func exec_KEY_TTL(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：获取键存储的值的数据类型。
// 用法：TYPE key
// 返回值：键的数据类型，包括"string"、"list"、"set"、"zset"、"hash"或"nil"（键不存在）。
func exec_KEY_TYPE(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	return reply.MakeBulkReply([]byte(db.KeyType(key)))
}

func init() {
	// 注册
	RegisterCommand(_const.CMD_KEY_DEL, exec_KEY_DEL, -2)
	RegisterCommand(_const.CMD_KEY_EXISTS, exec_KEY_EXISTS, 2)
	RegisterCommand(_const.CMD_KEY_EXPIRE, exec_KEY_EXPIRE, 3)
	RegisterCommand(_const.CMD_KEY_EXPIREAT, exec_KEY_EXPIREAT, 3)
	RegisterCommand(_const.CMD_KEY_KEYS, exec_KEY_KEYS, 2)
	RegisterCommand(_const.CMD_KEY_MOVE, exec_KEY_MOVE, 3)
	RegisterCommand(_const.CMD_KEY_PERSIST, exec_KEY_PERSIST, 2)
	RegisterCommand(_const.CMD_KEY_PEXPIRE, exec_KEY_PEXPIRE, 3)
	RegisterCommand(_const.CMD_KEY_PEXPIREAT, exec_KEY_PEXPIREAT, 3)
	RegisterCommand(_const.CMD_KEY_PTTL, exec_KEY_PTTL, 2)
	RegisterCommand(_const.CMD_KEY_RANDOMKEY, exec_KEY_RANDOMKEY, 1)
	RegisterCommand(_const.CMD_KEY_RENAME, exec_KEY_RENAME, 3)
	RegisterCommand(_const.CMD_KEY_RENAMENX, exec_KEY_RENAMENX, 3)
	RegisterCommand(_const.CMD_KEY_TTL, exec_KEY_TTL, 2)
	RegisterCommand(_const.CMD_KEY_TYPE, exec_KEY_TYPE, 2)
}
