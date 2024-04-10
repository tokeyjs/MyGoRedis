package database

import (
	_const "MyGoRedis/const"
	"MyGoRedis/datastruct/myhash"
	"MyGoRedis/datastruct/mylist"
	"MyGoRedis/datastruct/myset"
	"MyGoRedis/datastruct/mystring"
	"MyGoRedis/datastruct/myzset"
	"MyGoRedis/interface/resp"
	"MyGoRedis/resp/reply"
	"strings"
	"sync"
)

//type DB struct {
//	index   int
//	Data    dict.Dict
//	expired sync.Map // map[mystring]int64 // 过期时间以毫秒为单位
//	aofAdd  func(cmd CmdLine)
//}
//
//type ExecFunc func(db *DB, args [][]byte) resp.Reply
//
//type CmdLine = [][]byte
//
//func makeDB() *DB {
//	return &DB{
//		Data: dict.MakeSyncDict(),
//		aofAdd: func(cmd CmdLine) {
//			// 空实现：aof数据恢复阶段
//		},
//	}
//}
//func (db *DB) Exec(c resp.Connection, cmdLine CmdLine) resp.Reply {
//	cmdName := strings.ToLower(mystring(cmdLine[0]))
//	cmd, ok := cmdTable[cmdName]
//	if !ok {
//		return reply.MakeStandardErrReply("ERR unknown command " + cmdName)
//	}
//	if !validArity(cmd.arity, cmdLine) {
//		return reply.MakeArgNumErrReply(cmdName)
//	}
//	return cmd.exector(db, cmdLine[1:])
//}
//
//// SET k v -> arity = 3 定长
//// EXIST k1 k2 arity = -2 变长 ：至少为2个
//func validArity(arity int, cmdArgs [][]byte) bool {
//	argNum := len(cmdArgs)
//	if arity >= 0 {
//		return argNum == arity
//	}
//	return argNum >= -arity
//}
//
//func (db *DB) GetEntity(key mystring) (*database.DataEntity, bool) {
//	raw, ok := db.Data.Get(key)
//	if !ok {
//		return nil, false
//	}
//	entity, _ := raw.(*database.DataEntity)
//	return entity, true
//}
//
//func (db *DB) PutEntity(key mystring, entity *database.DataEntity) int {
//	return db.Data.Put(key, entity)
//}
//
//func (db *DB) PutIfExists(key mystring, entity *database.DataEntity) int {
//	return db.Data.PutIfExists(key, entity)
//}
//
//func (db *DB) PutIfAbsent(key mystring, entity *database.DataEntity) int {
//	return db.Data.PutIfAbsent(key, entity)
//}
//
//func (db *DB) Remove(key mystring) int {
//	return db.Data.Remove(key)
//}
//
//func (db *DB) Removes(keys ...mystring) int {
//	deleted := 0
//	for _, key := range keys {
//		_, exists := db.Data.Get(key)
//		if exists {
//			db.Remove(key)
//			deleted++
//		}
//	}
//	return deleted
//}
//
//func (db *DB) Flush() {
//	db.Data.Clear()
//}

type DB struct {
	index   int
	data    sync.Map //map[string]interface{}   --> interface{}是各种类型的指针
	expired sync.Map // map[string]int64 // 过期时间以毫秒为单位
	aofAdd  func(cmd CmdLine)
}

// 设置index
func (db *DB) SetDBIndex(index int) {
	db.index = index
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply

type CmdLine = [][]byte

func makeDB() *DB {
	return &DB{
		aofAdd: func(cmd CmdLine) {
			// 空实现：aof数据恢复阶段
		},
	}
}

// 执行命令
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

// 获取元素
func (db *DB) GetEntity(key string) (interface{}, bool) {
	raw, ok := db.data.Load(key)
	if !ok {
		return nil, false
	}
	return raw, true
}

// 放入元素
func (db *DB) PutEntity(key string, entity interface{}) {
	db.data.Store(key, entity)
}

// 判断元素是否存在
func (db *DB) IsExists(key string) bool {
	_, ok := db.data.Load(key)
	return ok
}

// 如果存在替换值
func (db *DB) PutIfExists(key string, entity interface{}) int {
	if db.IsExists(key) {
		db.PutEntity(key, entity)
		return 1
	}
	return 0
}

// 如果不存在就设置kv
func (db *DB) PutIfAbsent(key string, entity interface{}) int {
	if !db.IsExists(key) {
		db.PutEntity(key, entity)
		return 1
	}
	return 0
}

// 删除元素
func (db *DB) Remove(key string) int {
	if !db.IsExists(key) {
		return 0
	}
	db.data.Delete(key)
	return 1
}

// 获取类型
func (db *DB) KeyType(key string) string {
	val, ok := db.GetEntity(key)
	if !ok {
		return ""
	}
	switch val.(type) {
	case *myhash.Hash:
		return _const.TYPE_DATA_HASH
	case *mystring.String:
		return _const.TYPE_DATA_STRING
	case *mylist.List:
		return _const.TYPE_DATA_LIST
	case *myset.Set:
		return _const.TYPE_DATA_SET
	case *myzset.ZSet:
		return _const.TYPE_DATA_ZSET
	default:
		return _const.TYPE_DATA_NONE
	}
}

// 返回n个随机的key
func (db *DB) RandomKey(count int32) []string {
	if count <= 0 {
		return nil
	}
	slic := make([]string, 0, count)
	for i := 0; i < int(count); i++ {
		db.data.Range(func(key, value any) bool {
			slic = append(slic, key.(string))
			return false
		})
	}
	return slic
}

// 清空元素
func (db *DB) Flush() {
	db.data = sync.Map{}
}
