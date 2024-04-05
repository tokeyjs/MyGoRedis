package database

import (
	"MyGoRedis/interface/resp"
	"MyGoRedis/lib/utils"
	"MyGoRedis/lib/wildcard"
	"MyGoRedis/resp/reply"
)

// DEL k1 k2 k3...
func execDEL(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	deleted := db.Removes(keys...)
	if deleted > 0 {
		db.aofAdd(utils.ToCmdLine2("del", args...))
	}
	return reply.MakeIntReply(int64(deleted))
}

// EXISTS k1 k2 k3...
func execEXISTS(db *DB, args [][]byte) resp.Reply {
	num := 0
	for _, v := range args {
		if _, ok := db.GetEntity(string(v)); ok {
			num++
		}
	}
	return reply.MakeIntReply(int64(num))
}

// KEYS *
func execKEYS(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.data.ForEach(func(key string, val interface{}) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})
	return reply.MakeMultiBulkReply(result)
}

// FLUSHDB a b c
func execFLUSHDB(db *DB, args [][]byte) resp.Reply {
	db.data.Clear()
	db.aofAdd(utils.ToCmdLine2("flushdb", args...))
	return reply.MakeOkReply()
}

// TYPE k1
func execTYPE(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeStatusReply("none")
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.MakeStatusReply("string")
	}

	return reply.MakeUnknownErrReply()
}

// RENAME k1 k2
func execRENAME(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])
	entity, exists := db.GetEntity(src)
	if !exists {
		return reply.MakeStandardErrReply("no such key")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)
	db.aofAdd(utils.ToCmdLine2("rename", args...))
	return reply.MakeOkReply()
}

// RENAMENX k1 k2
func execRENAMENX(db *DB, args [][]byte) resp.Reply {
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
	db.aofAdd(utils.ToCmdLine2("renamenx", args...))
	return reply.MakeIntReply(1)
}

func init() {
	// 注册
	RegisterCommand("DEL", execDEL, -2)
	RegisterCommand("EXISTS", execEXISTS, -2)
	RegisterCommand("FLUSHDB", execFLUSHDB, -1)
	RegisterCommand("TYPE", execTYPE, 2)
	RegisterCommand("RENAME", execRENAME, 3)
	RegisterCommand("RENAMENX", execRENAMENX, 3)
	RegisterCommand("KEYS", execKEYS, 2)

}
