package database

import (
	_const "MyGoRedis/const"
	"MyGoRedis/datastruct/myhash"
	"MyGoRedis/datastruct/mylist"
	"MyGoRedis/datastruct/myset"
	"MyGoRedis/datastruct/mystring"
	"MyGoRedis/datastruct/myzset"
	"MyGoRedis/interface/resp"
	"MyGoRedis/lib/logger"
	"MyGoRedis/lib/utils"
	"MyGoRedis/resp/reply"
	"strings"
	"sync"
	"time"
)

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
	db := &DB{
		aofAdd: func(cmd CmdLine) {
			// 空实现：aof数据恢复阶段
		},
	}
	// 启动定期删除协程
	go db.deleteRegularlyWorker()
	return db
}

// 定期删除过期key的协程 [定期删除] todo aof添加del语句
func (db *DB) deleteRegularlyWorker() {
	ticker := time.NewTicker(10 * time.Second)
	for true {
		select {
		case <-ticker.C:
			go func() {
				// 删除过期的key
				delKey := make([]string, 0)
				// 当前时间
				nowmsec := time.Now().UnixNano() / 1e6
				db.expired.Range(func(key, value any) bool {
					if value.(int64) <= nowmsec {
						delKey = append(delKey, key.(string))
					}
					return true
				})
				for _, key := range delKey {
					// 删除过期
					db.Remove(key)
					db.aofAdd(utils.ToCmdLine2(_const.CMD_KEY_DEL, []byte(key)))
					logger.Debugf("deleteRegularly: db=[%v] key=[%v]", db.index, key)
				}
			}()
		}
	}
}

// 提前校验key是否过期 (过期会自动删除)  [懒惰删除]
func (db *DB) isOverdue(key string) bool {
	oTimeMil, ok := db.expired.Load(key)
	if !ok {
		return false
	}
	timestampmil, _ := oTimeMil.(int64)
	curTime := time.Now().UnixNano() / 1e6
	if curTime >= timestampmil {
		// 过期了进行删除
		db.Remove(key)
		db.aofAdd(utils.ToCmdLine2(_const.CMD_KEY_DEL, []byte(key)))
		logger.Debugf("overdue delete: db=[%v] key=[%v]", db.index, key)
		return true
	}
	return false
}

// 设置key的过期时间（毫秒）
func (db *DB) SetExpiredMSec(key string, msec int64) {
	if msec < 0 {
		// 永久
		db.expired.Delete(key)
		return
	}
	if _, ok := db.data.Load(key); !ok {
		return
	}
	db.expired.Store(key, msec+time.Now().UnixNano()/1e6)
}

// 设置key的过期时间戳（毫秒）
func (db *DB) SetExpiredTimestampsec(key string, timestampmsec int64) {
	if timestampmsec < 0 {
		db.expired.Delete(key)
		return
	}
	if _, ok := db.data.Load(key); !ok {
		return
	}
	db.expired.Store(key, timestampmsec)
}

// 返回key的过期时间戳（毫秒级别）
func (db *DB) GetExpiredTimestampMSec(key string) int64 {
	timestampMSec, ok := db.expired.Load(key)
	if !ok {
		return -1
	}
	return timestampMSec.(int64)
}

// 返回key还有多少时间过期（毫秒）
func (db *DB) GetExpiredMSec(key string) int64 {
	timestampMSec, ok := db.expired.Load(key)
	if !ok {
		return -1
	}
	// 剩余毫秒
	msec := timestampMSec.(int64) - time.Now().UnixNano()/1e6
	if msec <= 0 {
		// 过期了当作不存在
		return -2
	}
	return msec
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
	if !ok || db.isOverdue(key) {
		return nil, false
	}
	return raw, true
}

// 放入元素
func (db *DB) PutEntity(key string, entity interface{}) {
	db.data.Store(key, entity)
}

// 放入元素并设置过期时间msec
func (db *DB) PutEntityWithMsec(key string, entity interface{}, msec int64) {
	db.data.Store(key, entity)
	db.SetExpiredMSec(key, msec)
}

// 放入元素并设置过期时间毫秒级时间戳
func (db *DB) PutEntityWithTimestamp(key string, entity interface{}, timestampmsec int64) {
	db.data.Store(key, entity)
	db.SetExpiredTimestampsec(key, timestampmsec)
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
	db.expired.Delete(key)
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
	db.expired = sync.Map{}
}
