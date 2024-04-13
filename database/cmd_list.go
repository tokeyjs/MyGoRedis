package database

import (
	_const "MyGoRedis/const"
	"MyGoRedis/datastruct/mylist"
	"MyGoRedis/interface/resp"
	"MyGoRedis/lib/utils"
	"MyGoRedis/resp/reply"
	"strconv"
	"strings"
)

// 实现命令

// 检查完成

// 含义：阻塞式弹出列表最左边的元素。
// 用法：BLPOP key1 timeout
// 返回值：返回被弹出的元素和对应的键。
func exec_LIST_BLPOP(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：阻塞式弹出列表最右边的元素。
// 用法：BRPOP key1 timeout
// 返回值：返回被弹出的元素和对应的键。
func exec_LIST_BRPOP(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：阻塞式弹出一个列表最右边的元素，并将它推入另一个列表的最左边。
// 用法：BRPOPLPUSH source destination timeout
// 返回值：被弹出的元素。
func exec_LIST_BRPOPLPUSH(db *DB, args [][]byte) resp.Reply {
	// todo
	return reply.MakeUnknownErrReply()
}

// 含义：获取列表中指定位置的元素。
// 用法：LINDEX key index
// 返回值：列表中指定位置的元素。
func exec_LIST_LINDEX(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	index, err := strconv.Atoi(string(args[1]))
	if err != nil {
		return reply.MakeStandardErrReply("index is error")
	}
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeList := _const.DataToLIST(it)
	if typeList == nil {
		return reply.MakeUnknownErrReply()
	}
	data, err := typeList.GetByIndex(int32(index))
	if err != nil {
		return reply.MakeStandardErrReply("index is out range")
	}
	return reply.MakeBulkReply([]byte(data))
}

// 含义：在列表中指定元素的前面或后面插入新元素。
// 用法：LINSERT key BEFORE|AFTER pivot value
// 返回值：插入操作完成后列表的长度。
func exec_LIST_LINSERT(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	index := strings.ToLower(string(args[1]))
	if index != "before" && index != "after" {
		return reply.MakeStandardErrReply("index is error")
	}
	db.aofAdd(utils.ToCmdLine2(_const.CMD_LIST_LINSERT, args...))
	pivot := string(args[2])
	value := string(args[3])
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeList := _const.DataToLIST(it)
	if typeList == nil {
		return reply.MakeUnknownErrReply()
	}
	var err error
	if index == "before" {
		err = typeList.InsertByValue(value, pivot, true)
	} else {
		err = typeList.InsertByValue(value, pivot, false)
	}
	if err != nil {
		return reply.MakeStandardErrReply("insert err: " + err.Error())
	}
	return reply.MakeIntReply(int64(typeList.Size()))
}

// 含义：获取列表的长度。
// 用法：LLEN key
// 返回值：列表的长度。
func exec_LIST_LLEN(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeIntReply(0)
	}
	typeList := _const.DataToLIST(it)
	if typeList == nil {
		return reply.MakeUnknownErrReply()
	}
	return reply.MakeIntReply(int64(typeList.Size()))
}

// 含义：弹出列表最左边的元素。
// 用法：LPOP key
// 返回值：被弹出的元素。
func exec_LIST_LPOP(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_LIST_LPOP, args...))
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeList := _const.DataToLIST(it)
	if typeList == nil {
		return reply.MakeUnknownErrReply()
	}
	data, err := typeList.PopBegin()
	if err != nil {
		return reply.MakeStandardErrReply("index is out range")
	}
	return reply.MakeBulkReply([]byte(data))
}

// 含义：将一个或多个值插入列表的头部。
// 用法：LPUSH key value1 [value2 ...]
// 返回值：插入操作完成后列表的长度。
func exec_LIST_LPUSH(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_LIST_LPUSH, args...))
	var typeList *mylist.List
	it, ok := db.GetEntity(key)
	if !ok {
		typeList = mylist.MakeList()
	} else {
		typeList = _const.DataToLIST(it)
		if typeList == nil {
			typeList = mylist.MakeList()
		}
	}
	index := 1
	for index < len(args) {
		value := string(args[index])
		_ = typeList.PushBegin(value)
		index++
	}
	db.PutEntity(key, typeList)
	return reply.MakeIntReply(int64(typeList.Size()))
}

// 含义：将值插入到已存在的列表头部。
// 用法：LPUSHX key value
// 返回值：插入操作完成后列表的长度，若列表不存在则返回0。
func exec_LIST_LPUSHX(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_LIST_LPUSHX, args...))
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeIntReply(0)
	}
	typeList := _const.DataToLIST(it)
	if typeList == nil {
		// 错误
		return reply.MakeUnknownErrReply()
	}
	value := string(args[1])
	_ = typeList.PushBegin(value)
	return reply.MakeIntReply(int64(typeList.Size()))
}

// 含义：从列表中删除指定数量的指定元素。
// 用法：LREM key count value
// 返回值：被移除的元素数量。
func exec_LIST_LREM(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_LIST_LREM, args...))
	count, err := strconv.Atoi(string(args[1]))
	if err != nil {
		return reply.MakeStandardErrReply("count is error")
	}
	value := string(args[2])
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeIntReply(0)
	}
	typeList := _const.DataToLIST(it)
	if typeList == nil {
		return reply.MakeUnknownErrReply()
	}
	return reply.MakeIntReply(int64(typeList.Remove(int32(count), value)))
}

// 含义：设置列表指定位置的值。
// 用法：LSET key index value
// 返回值：操作成功则返回OK。
func exec_LIST_LSET(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_LIST_LSET, args...))
	index, err := strconv.Atoi(string(args[1]))
	if err != nil {
		return reply.MakeStandardErrReply("index is error")
	}
	value := string(args[2])
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeStandardErrReply("list not exists")
	}
	typeList := _const.DataToLIST(it)
	if typeList == nil {
		return reply.MakeUnknownErrReply()
	}
	err = typeList.SetByIndex(int32(index), value)
	if err != nil {
		return reply.MakeStandardErrReply("error :" + err.Error())
	}
	return reply.MakeOkReply()
}

// 含义：弹出列表最右边的元素。
// 用法：RPOP key
// 返回值：被弹出的元素。
func exec_LIST_RPOP(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_LIST_RPOP, args...))
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeList := _const.DataToLIST(it)
	if typeList == nil {
		return reply.MakeUnknownErrReply()
	}
	back, err := typeList.PopBack()
	if err != nil {
		return reply.MakeStandardErrReply("error :" + err.Error())
	}
	return reply.MakeBulkReply([]byte(back))
}

// 含义：弹出一个列表最右边的元素，并将它推入另一个列表的最左边。
// 用法：RPOPLPUSH source destination
// 返回值：被弹出的元素。
func exec_LIST_RPOPLPUSH(db *DB, args [][]byte) resp.Reply {
	source := string(args[0])
	dest := string(args[1])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_LIST_RPOPLPUSH, args...))
	it, ok := db.GetEntity(source)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeList := _const.DataToLIST(it)
	if typeList == nil {
		return reply.MakeUnknownErrReply()
	}
	it2, ok := db.GetEntity(dest)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeList2 := _const.DataToLIST(it2)
	if typeList2 == nil {
		return reply.MakeNullBulkReply()
	}
	// 弹出右边
	back, err := typeList.PopBack()
	if err != nil {
		return reply.MakeNullBulkReply()
	}
	// 放入左边
	_ = typeList2.PushBegin(back)
	return reply.MakeBulkReply([]byte(back))
}

// 含义：将一个或多个值插入到列表的右侧（尾部）。
// 用法：RPUSH key value1 [value2 ...]
// 返回值：执行命令后列表的长度。
func exec_LIST_RPUSH(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_LIST_RPUSH, args...))
	var typeList *mylist.List
	it, ok := db.GetEntity(key)
	if !ok {
		typeList = mylist.MakeList()
	} else {
		typeList = _const.DataToLIST(it)
		if typeList == nil {
			typeList = mylist.MakeList()
		}
	}
	index := 1
	for index < len(args) {
		value := string(args[index])
		typeList.PushBack(value)
		index++
	}
	db.PutEntity(key, typeList)
	return reply.MakeIntReply(int64(typeList.Size()))
}

// 含义：将一个值插入到已存在的列表的右侧（尾部）。
// 用法：RPUSHX key value
// 返回值：如果列表存在，则返回插入后列表的长度；如果列表不存在，则不执行插入操作，返回0。
func exec_LIST_RPUSHX(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	db.aofAdd(utils.ToCmdLine2(_const.CMD_LIST_LPUSHX, args...))
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeIntReply(0)
	}
	typeList := _const.DataToLIST(it)
	if typeList == nil {
		return reply.MakeUnknownErrReply()
	}
	value := string(args[1])
	typeList.PushBack(value)
	return reply.MakeIntReply(int64(typeList.Size()))
}

// LRANGE key start stop
// 返回列表 key 中指定区间内的元素，区间以偏移量 start 和 stop 指定。
// 返回值: 一个列表，包含指定区间内的元素。
func exec_LIST_LRANGE(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	start, err := strconv.Atoi(string(args[1]))
	if err != nil {
		return reply.MakeStandardErrReply("start is not int")
	}
	stop, err := strconv.Atoi(string(args[2]))
	if err != nil {
		return reply.MakeStandardErrReply("end is not int")
	}
	it, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	typeList := _const.DataToLIST(it)
	if typeList == nil {
		return reply.MakeUnknownErrReply()
	}
	data := typeList.GetRange(int32(start), int32(stop))
	if len(data) == 0 {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeMultiBulkReply(utils.ToCmdLine(data...))
}

func init() {
	RegisterCommand(_const.CMD_LIST_BLPOP, exec_LIST_BLPOP, 3)
	RegisterCommand(_const.CMD_LIST_BRPOP, exec_LIST_BRPOP, 3)
	RegisterCommand(_const.CMD_LIST_BRPOPLPUSH, exec_LIST_BRPOPLPUSH, 4)
	RegisterCommand(_const.CMD_LIST_LINDEX, exec_LIST_LINDEX, 3)
	RegisterCommand(_const.CMD_LIST_LINSERT, exec_LIST_LINSERT, 5)
	RegisterCommand(_const.CMD_LIST_LLEN, exec_LIST_LLEN, 2)
	RegisterCommand(_const.CMD_LIST_LPOP, exec_LIST_LPOP, 2)
	RegisterCommand(_const.CMD_LIST_LPUSH, exec_LIST_LPUSH, -3)
	RegisterCommand(_const.CMD_LIST_LPUSHX, exec_LIST_LPUSHX, 3)
	RegisterCommand(_const.CMD_LIST_LREM, exec_LIST_LREM, 4)
	RegisterCommand(_const.CMD_LIST_LSET, exec_LIST_LSET, 4)
	RegisterCommand(_const.CMD_LIST_RPOP, exec_LIST_RPOP, 2)
	RegisterCommand(_const.CMD_LIST_RPOPLPUSH, exec_LIST_RPOPLPUSH, 3)
	RegisterCommand(_const.CMD_LIST_RPUSH, exec_LIST_RPUSH, -3)
	RegisterCommand(_const.CMD_LIST_RPUSHX, exec_LIST_RPUSHX, 3)
	RegisterCommand(_const.CMD_LIST_LRANGE, exec_LIST_LRANGE, 4)

}
