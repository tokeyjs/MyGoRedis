package database

import (
	_const "MyGoRedis/const"
	"MyGoRedis/interface/resp"
	"MyGoRedis/lib/utils"
	"MyGoRedis/resp/reply"
	"strconv"
	"time"
)

// 实现命令
// 含义：获取服务器的当前时间，以UNIX时间戳和微秒的格式返回。
// 用法：TIME
// 返回值：返回一个包含两个元素的数组，第一个元素是当前时间的UNIX时间戳（以秒为单位），第二个元素是微秒偏移量。
func exec_SERVER_TIME(db *DB, args [][]byte) resp.Reply {
	// 获取当前时间
	now := time.Now()
	// 计算当前时间的UNIX时间戳（秒）
	secondsSinceEpoch := now.Unix()
	// 计算微秒偏移量（纳秒转换为微秒）
	microsecondOffset := now.Nanosecond() / 1000 // 因为1微秒 = 1000纳秒
	arr := make([][]byte, 0)
	arr = append(arr, []byte(strconv.FormatInt(secondsSinceEpoch, 10)))
	arr = append(arr, []byte(strconv.FormatInt(int64(microsecondOffset), 10)))
	return reply.MakeMultiBulkReply(arr)
}

// 含义：清空当前数据库的所有数据。
// 用法：FLUSHDB
// 返回值：执行成功时返回OK
func exec_SERVER_FLUSHDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	db.aofAdd(utils.ToCmdLine2(_const.CMD_SERVER_FLUSHDB, args...))
	return reply.MakeOkReply()
}

func init() {
	// 注册
	RegisterCommand(_const.CMD_SERVER_TIME, exec_SERVER_TIME, 1)
	RegisterCommand(_const.CMD_SERVER_FLUSHDB, exec_SERVER_FLUSHDB, 1)
}
