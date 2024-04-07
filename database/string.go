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

//SET key value            	设置指定 key 的值。
//GET key 					获取指定 key 的值。
//GETSET key value 			将给定 key 的值设为 value ，并返回 key 的旧值(old value)。
//MGET key1 [key2..] 			获取所有(一个或多个)给定 key 的值。
//SETEX key seconds value 	将值 value 关联到 key ，并将 key 的过期时间设为 seconds (以秒为单位)。
//SETNX key value 			只有在 key 不存在时设置 key 的值。
//STRLEN key 					返回 key 所储存的字符串值的长度。
//MSET key value [key value ...] 同时设置一个或多个 key-value 对。
//MSETNX key value [key value ...] 	同时设置一个或多个 key-value 对，当且仅当所有给定 key 都不存在。
//PSETEX key milliseconds value 这个命令和 SETEX 命令相似，但它以毫秒为单位设置 key 的生存时间，而不是像 SETEX 命令那样，以秒为单位。
//INCR key 					将 key 中储存的数字值增一。
//INCRBY key increment 		将 key 所储存的值加上给定的增量值（increment） 。
//DECR key 					将 key 中储存的数字值减一。
//DECRBY key decrement key 	所储存的值减去给定的减量值（decrement） 。
//APPEND key value 			如果 key 已经存在并且是一个字符串， APPEND 命令将指定的 value 追加到该 key 原来值（value）的末尾。

func init() {
	RegisterCommand("GET", execGET, 2)
	RegisterCommand("SET", execSET, 3)
	RegisterCommand("SETNX", execSETNX, 3)
	RegisterCommand("GETSET", execGETSET, 3)
	RegisterCommand("STRLEN", execSTRLEN, 2)
}
