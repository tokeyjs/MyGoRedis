package database

import (
	"MyGoRedis/interface/database"
	"MyGoRedis/interface/resp"
	"MyGoRedis/lib/utils"
	"MyGoRedis/resp/reply"
)

// GET k1
func execGET(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeBulkReply(entity.Data.([]byte))
}

// SET k1 v1
func execSET(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := args[1]
	entity := &database.DataEntity{Data: value}
	db.PutEntity(key, entity)
	db.aofAdd(utils.ToCmdLine2("set", args...))
	return reply.MakeOkReply()
}

// SETNX k1 v1
func execSETNX(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := args[1]
	entity := &database.DataEntity{Data: value}
	res := db.PutIfAbsent(key, entity)
	db.aofAdd(utils.ToCmdLine2("setnx", args...))
	return reply.MakeIntReply(int64(res))
}

// GETSET k1 v1
func execGETSET(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := args[1]
	entity, exists := db.GetEntity(key)
	db.PutEntity(key, &database.DataEntity{
		Data: value,
	})
	db.aofAdd(utils.ToCmdLine2("getset", args...))
	if !exists {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeBulkReply(entity.Data.([]byte))
}

// STRLEN k1
func execSTRLEN(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeNullBulkReply()
	}
	byes := entity.Data.([]byte)
	return reply.MakeIntReply(int64(len(byes)))
}

func init() {
	RegisterCommand("GET", execGET, 2)
	RegisterCommand("SET", execSET, 3)
	RegisterCommand("SETNX", execSETNX, 3)
	RegisterCommand("GETSET", execGETSET, 3)
	RegisterCommand("STRLEN", execSTRLEN, 2)
}
