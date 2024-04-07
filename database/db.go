package database

import (
	"MyGoRedis/datastruct/dict"
	"MyGoRedis/interface/database"
	"MyGoRedis/interface/resp"
	"MyGoRedis/resp/reply"
	"strings"
	"sync"
)

type DB struct {
	index   int
	data    dict.Dict
	expired sync.Map // map[string]int64 // 过期时间以毫秒为单位
	aofAdd  func(cmd CmdLine)
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply

type CmdLine = [][]byte

func makeDB() *DB {
	return &DB{
		data: dict.MakeSyncDict(),
		aofAdd: func(cmd CmdLine) {
			// 空实现：aof数据恢复阶段
		},
	}
}
func (db *DB) Exec(c resp.Connection, cmdLine CmdLine) resp.Reply {
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.MakeStandardErrReply("ERR unknown command " + cmdName)
	}
	if !validArity(cmd.arity, cmdLine) {
		return reply.MakeArgNumErrReply(cmdName)
	}
	return cmd.exector(db, cmdLine[1:])
}

// SET k v -> arity = 3 定长
// EXIST k1 k2 arity = -2 变长 ：至少为2个
func validArity(arity int, cmdArgs [][]byte) bool {
	argNum := len(cmdArgs)
	if arity >= 0 {
		return argNum == arity
	}
	return argNum >= -arity
}

func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	entity, _ := raw.(*database.DataEntity)
	return entity, true
}

func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

func (db *DB) Remove(key string) int {
	return db.data.Remove(key)
}

func (db *DB) Removes(keys ...string) int {
	deleted := 0
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.Remove(key)
			deleted++
		}
	}
	return deleted
}

func (db *DB) Flush() {
	db.data.Clear()
}
