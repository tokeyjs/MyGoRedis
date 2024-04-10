package database

import (
	_const "MyGoRedis/const"
	"MyGoRedis/interface/resp"
	"MyGoRedis/resp/reply"
)

// 实现命令

//// DEL k1 k2 k3...
//func execDEL(db *DB, args [][]byte) resp.Reply {
//	keys := make([]mystring, len(args))
//	for i, v := range args {
//		keys[i] = mystring(v)
//	}
//	deleted := db.Removes(keys...)
//	if deleted > 0 {
//		db.aofAdd(utils.ToCmdLine2("del", args...))
//	}
//	return reply.MakeIntReply(int64(deleted))
//}
//
//// EXISTS k1 k2 k3...
//func execEXISTS(db *DB, args [][]byte) resp.Reply {
//	num := 0
//	for _, v := range args {
//		if _, ok := db.GetEntity(mystring(v)); ok {
//			num++
//		}
//	}
//	return reply.MakeIntReply(int64(num))
//}
//
//// KEYS *
//func execKEYS(db *DB, args [][]byte) resp.Reply {
//	pattern := wildcard.CompilePattern(mystring(args[0]))
//	result := make([][]byte, 0)
//	db.Data.ForEach(func(key mystring, val interface{}) bool {
//		if pattern.IsMatch(key) {
//			result = append(result, []byte(key))
//		}
//		return true
//	})
//	return reply.MakeMultiBulkReply(result)
//}
//
//// FLUSHDB a b c
//func execFLUSHDB(db *DB, args [][]byte) resp.Reply {
//	db.Data.Clear()
//	db.aofAdd(utils.ToCmdLine2("flushdb", args...))
//	return reply.MakeOkReply()
//}
//
//// TYPE k1
//func execTYPE(db *DB, args [][]byte) resp.Reply {
//	key := mystring(args[0])
//	entity, exists := db.GetEntity(key)
//	if !exists {
//		return reply.MakeStatusReply("none")
//	}
//	switch entity.Data.(type) {
//	case []byte:
//		return reply.MakeStatusReply("mystring")
//	}
//
//	return reply.MakeUnknownErrReply()
//}
//
//// RENAME k1 k2
//func execRENAME(db *DB, args [][]byte) resp.Reply {
//	src := mystring(args[0])
//	dest := mystring(args[1])
//	entity, exists := db.GetEntity(src)
//	if !exists {
//		return reply.MakeStandardErrReply("no such key")
//	}
//	db.PutEntity(dest, entity)
//	db.Remove(src)
//	db.aofAdd(utils.ToCmdLine2("rename", args...))
//	return reply.MakeOkReply()
//}
//
//// RENAMENX k1 k2
//func execRENAMENX(db *DB, args [][]byte) resp.Reply {
//	src := string(args[0])
//	dest := string(args[1])
//	if _, ok := db.GetEntity(dest); ok {
//		return reply.MakeIntReply(0)
//	}
//	entity, exists := db.GetEntity(src)
//	if !exists {
//		return reply.MakeStandardErrReply("no such key")
//	}
//	db.PutEntity(dest, entity)
//	db.Remove(src)
//	db.aofAdd(utils.ToCmdLine2("renamenx", args...))
//	return reply.MakeIntReply(1)
//}
//
//func init() {
//	// 注册
//	RegisterCommand("DEL", execDEL, -2)
//	RegisterCommand("EXISTS", execEXISTS, -2)
//	RegisterCommand("FLUSHDB", execFLUSHDB, -1)
//	RegisterCommand("TYPE", execTYPE, 2)
//	RegisterCommand("RENAME", execRENAME, 3)
//	RegisterCommand("RENAMENX", execRENAMENX, 3)
//	RegisterCommand("KEYS", execKEYS, 2)
//
//}

// ===实现命令===

// 含义：删除给定的一个或多个键。
// 用法：DEL key [key ...]
// 返回值：被删除键的数量。
func exec_KEY_DEL(db *DB, args [][]byte) resp.Reply {
	// todo
	count := 0
	for _, arg := range args {
		count += db.Remove(string(arg))
	}
	return reply.MakeUnknownErrReply()
}

// 含义：检查给定键是否存在。
// 用法：EXISTS key
// 返回值：若键存在，则返回1，否则返回0。
func exec_KEY_EXISTS(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：设置键的过期时间，单位为秒。
// 用法：EXPIRE key seconds
// 返回值：若设置成功，则返回1，否则返回0。
func exec_KEY_EXPIRE(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：设置键在指定的Unix时间戳（秒级）过期。
// 用法：EXPIREAT key timestamp
// 返回值：若设置成功，则返回1，否则返回0。
func exec_KEY_EXPIREAT(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：查找所有符合给定模式的键。
// 用法：KEYS pattern
// 返回值：符合模式的键组成的列表。
func exec_KEY_KEYS(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：将指定键移到另一个数据库。
// 用法：MOVE key db
// 返回值：若移动成功，则返回1，否则返回0。
func exec_KEY_MOVE(db *DB, args [][]byte) resp.Reply {
	//todo
	return reply.MakeIntReply(0)
}

// 含义：移除键的过期时间，使其永不过期。
// 用法：PERSIST key
// 返回值：若键成功移除过期时间，则返回1，否则返回0。
func exec_KEY_PERSIST(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：设置键的过期时间，单位为毫秒。
// 用法：PEXPIRE key milliseconds
// 返回值：若设置成功，则返回1，否则返回0。
func exec_KEY_PEXPIRE(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：设置键在指定的Unix时间戳（毫秒级）过期。
// 用法：PEXPIREAT key milliseconds-timestamp
// 返回值：若设置成功，则返回1，否则返回0。
func exec_KEY_PEXPIREAT(db *DB, args [][]byte) resp.Reply {
	// todo
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
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：将键重命名为新键名。
// 用法：RENAME key newkey
// 返回值：若重命名成功，则返回OK。
func exec_KEY_RENAME(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：仅当新键不存在时，将键重命名为新键名。
// 用法：RENAMENX key newkey
// 返回值：若重命名成功，则返回1，否则返回0。
func exec_KEY_RENAMENX(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
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
	// todo
	return reply.MakeUnknownErrReply()
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