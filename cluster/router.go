package cluster

import (
	_const "MyGoRedis/const"
	"MyGoRedis/interface/resp"
	"MyGoRedis/lib/utils"
	"MyGoRedis/resp/reply"
)

func makeRouter() map[string]CmdFunc {
	return map[string]CmdFunc{
		// server
		_const.CMD_SERVER_TIME:    curNodeFunc,
		_const.CMD_SERVER_FLUSHDB: flushdbFunc,
		// conn
		_const.CMD_CONN_PING:   curNodeFunc,
		_const.CMD_CONN_AUTH:   curNodeFunc,
		_const.CMD_CONN_ECHO:   curNodeFunc,
		_const.CMD_CONN_SELECT: curNodeFunc,
		// key
		_const.CMD_KEY_DEL:       delFunc,
		_const.CMD_KEY_EXISTS:    selectNodeFunc,
		_const.CMD_KEY_EXPIRE:    selectNodeFunc,
		_const.CMD_KEY_EXPIREAT:  selectNodeFunc,
		_const.CMD_KEY_KEYS:      keysFunc,
		_const.CMD_KEY_MOVE:      selectNodeFunc,
		_const.CMD_KEY_PERSIST:   selectNodeFunc,
		_const.CMD_KEY_PEXPIRE:   selectNodeFunc,
		_const.CMD_KEY_PEXPIREAT: selectNodeFunc,
		_const.CMD_KEY_PTTL:      selectNodeFunc,
		_const.CMD_KEY_RANDOMKEY: curNodeFunc,
		_const.CMD_KEY_RENAME:    renameFunc,
		_const.CMD_KEY_RENAMENX:  renamenxFunc,
		_const.CMD_KEY_TTL:       selectNodeFunc,
		_const.CMD_KEY_TYPE:      selectNodeFunc,

		// string
		_const.CMD_STRING_APPEND:      selectNodeFunc,
		_const.CMD_STRING_DECR:        selectNodeFunc,
		_const.CMD_STRING_DECRBY:      selectNodeFunc,
		_const.CMD_STRING_GET:         selectNodeFunc,
		_const.CMD_STRING_GETRANGE:    selectNodeFunc,
		_const.CMD_STRING_GETSET:      selectNodeFunc,
		_const.CMD_STRING_INCR:        selectNodeFunc,
		_const.CMD_STRING_INCRBY:      selectNodeFunc,
		_const.CMD_STRING_INCRBYFLOAT: selectNodeFunc,
		_const.CMD_STRING_MGET:        mgetFunc,
		_const.CMD_STRING_MSET:        msetFunc,
		_const.CMD_STRING_MSETNX:      msetnxFunc,
		_const.CMD_STRING_PSETEX:      selectNodeFunc,
		_const.CMD_STRING_SET:         selectNodeFunc,
		_const.CMD_STRING_SETEX:       selectNodeFunc,
		_const.CMD_STRING_SETNX:       selectNodeFunc,
		_const.CMD_STRING_STRLEN:      selectNodeFunc,

		// list
		_const.CMD_LIST_BLPOP:      selectNodeFunc,
		_const.CMD_LIST_BRPOP:      selectNodeFunc,
		_const.CMD_LIST_BRPOPLPUSH: selectNodeFunc,
		_const.CMD_LIST_LINDEX:     selectNodeFunc,
		_const.CMD_LIST_LINSERT:    selectNodeFunc,
		_const.CMD_LIST_LLEN:       selectNodeFunc,
		_const.CMD_LIST_LPOP:       selectNodeFunc,
		_const.CMD_LIST_LPUSH:      selectNodeFunc,
		_const.CMD_LIST_LPUSHX:     selectNodeFunc,
		_const.CMD_LIST_LREM:       selectNodeFunc,
		_const.CMD_LIST_LSET:       selectNodeFunc,
		_const.CMD_LIST_RPOP:       selectNodeFunc,
		_const.CMD_LIST_RPOPLPUSH:  rpoppushFunc,
		_const.CMD_LIST_RPUSH:      selectNodeFunc,
		_const.CMD_LIST_RPUSHX:     selectNodeFunc,
		_const.CMD_LIST_LRANGE:     selectNodeFunc,

		// hash
		_const.CMD_HASH_HDEL:         selectNodeFunc,
		_const.CMD_HASH_HEXISTS:      selectNodeFunc,
		_const.CMD_HASH_HGET:         selectNodeFunc,
		_const.CMD_HASH_HGETALL:      selectNodeFunc,
		_const.CMD_HASH_HINCRBY:      selectNodeFunc,
		_const.CMD_HASH_HINCRBYFLOAT: selectNodeFunc,
		_const.CMD_HASH_HKEYS:        selectNodeFunc,
		_const.CMD_HASH_HLEN:         selectNodeFunc,
		_const.CMD_HASH_HMGET:        selectNodeFunc,
		_const.CMD_HASH_HMSET:        selectNodeFunc,
		_const.CMD_HASH_HSET:         selectNodeFunc,
		_const.CMD_HASH_HSETNX:       selectNodeFunc,
		_const.CMD_HASH_HVALS:        selectNodeFunc,

		// set
		_const.CMD_SET_SADD:        selectNodeFunc,
		_const.CMD_SET_SCARD:       selectNodeFunc,
		_const.CMD_SET_SINTER:      sinterFunc,
		_const.CMD_SET_SINTERSTORE: sinterstoreFunc,
		_const.CMD_SET_SISMEMBER:   selectNodeFunc,
		_const.CMD_SET_SMEMBERS:    selectNodeFunc,
		_const.CMD_SET_SMOVE:       smoveFunc,
		_const.CMD_SET_SPOP:        selectNodeFunc,
		_const.CMD_SET_SREM:        selectNodeFunc,
		_const.CMD_SET_SRANDMEMBER: selectNodeFunc,

		// zset
		_const.CMD_ZSET_ZADD:             selectNodeFunc,
		_const.CMD_ZSET_ZCARD:            selectNodeFunc,
		_const.CMD_ZSET_ZCOUNT:           selectNodeFunc,
		_const.CMD_ZSET_ZINCRBY:          selectNodeFunc,
		_const.CMD_ZSET_ZRANGE:           selectNodeFunc,
		_const.CMD_ZSET_ZRANGEBYSCORE:    selectNodeFunc,
		_const.CMD_ZSET_ZRANK:            selectNodeFunc,
		_const.CMD_ZSET_ZREM:             selectNodeFunc,
		_const.CMD_ZSET_ZREMRANGEBYRANK:  selectNodeFunc,
		_const.CMD_ZSET_ZREMRANGEBYSCORE: selectNodeFunc,
		_const.CMD_ZSET_ZREVRANGE:        selectNodeFunc,
		_const.CMD_ZSET_ZREVRANGEBYSCORE: selectNodeFunc,
		_const.CMD_ZSET_ZREVRANK:         selectNodeFunc,
		_const.CMD_ZSET_ZSCORE:           selectNodeFunc,
	}
}

// 根据key将命令发送到指定节点进行执行
func selectNodeFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	key := string(cmdArgs[1])
	peer, _ := cluster.peerPicker.PickNode(key)
	return cluster.relay(peer, c, cmdArgs)
}

// 当前节点执行命令
func curNodeFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	return cluster.db.Exec(c, cmdArgs)
}

//------ 特例化集群命令执行方法 ------
// ===>[server]

// flushdb
func flushdbFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcast(c, cmdArgs)
	for _, r := range replies {
		if reply.IsErrReply(r) {
			rep := r.(reply.ErrorReply)
			return reply.MakeStandardErrReply("err: " + rep.Error())
		}
	}
	return reply.MakeOkReply()
}

// ==>[key]

// 含义：删除给定的一个或多个键。
// 用法：DEL key [key ...]
// 返回值：被删除键的数量。
func delFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcast(c, cmdArgs)
	var deleted int64 = 0
	for _, r := range replies {
		if reply.IsErrReply(r) {
			rep := r.(reply.ErrorReply)
			return reply.MakeStandardErrReply("err: " + rep.Error())
		}
		intReply, ok := r.(*reply.IntReply)
		if !ok {
			return reply.MakeStandardErrReply("error")
		}
		deleted += intReply.Code
	}
	return reply.MakeIntReply(deleted)
}

// 含义：查找所有符合给定模式的键。
// 用法：KEYS pattern
// 返回值：符合模式的键组成的列表。
func keysFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcast(c, cmdArgs)
	allRes := make([][]byte, 0)
	for _, r := range replies {
		if reply.IsErrReply(r) {
			continue
		}
		multiBulkReply, ok := r.(*reply.MultiBulkReply)
		if !ok {
			return reply.MakeStandardErrReply("error")
		}
		allRes = append(allRes, multiBulkReply.Args...)
	}
	return reply.MakeMultiBulkReply(allRes)
}

// 含义：将键重命名为新键名
// 用法：RENAME key newkey
// 返回值：若重命名成功，则返回OK
func renameFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.MakeStandardErrReply("wrong unmber args")
	}
	src := string(cmdArgs[1])
	dest := string(cmdArgs[2])
	srcPeer, _ := cluster.peerPicker.PickNode(src)
	destPeer, _ := cluster.peerPicker.PickNode(dest)
	if srcPeer != destPeer {
		// TODO
		// 获取value类型

		// src节点删除键

		// dest节点新增键
		return reply.MakeStandardErrReply("key and newKey is not exists same node")
	}
	return cluster.relay(srcPeer, c, cmdArgs)
}

// 含义：仅当新键不存在时，将键重命名为新键名。
// 用法：RENAMENX key newkey
// 返回值：若重命名成功，则返回1，否则返回0。
func renamenxFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.MakeStandardErrReply("wrong unmber args")
	}
	src := string(cmdArgs[1])
	dest := string(cmdArgs[2])
	srcPeer, _ := cluster.peerPicker.PickNode(src)
	destPeer, _ := cluster.peerPicker.PickNode(dest)
	if srcPeer != destPeer {
		// TODO
		// 获取value类型

		// src节点删除键

		// dest节点新增键
		return reply.MakeStandardErrReply("key and newKey is not exists same node")
	}
	return cluster.relay(srcPeer, c, cmdArgs)
}

// ==>[string]

// 含义：获取一个或多个键的值。
// 用法：MGET key [key ...]
// 返回值：包含指定键值的数组。
func mgetFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) < 2 {
		return reply.MakeArgNumErrReply("mget")
	}
	allRes := make([][]byte, 0)
	index := 1
	for index < len(cmdArgs) {
		key := string(cmdArgs[index])
		index++
		// 获取执行节点
		peer, _ := cluster.peerPicker.PickNode(key)
		res := cluster.relay(peer, c, utils.ToCmdLine2(_const.CMD_STRING_MGET, []byte(key)))
		if reply.IsErrReply(res) {
			allRes = append(allRes, nil)
			continue
		}
		multiBulkReply, ok := res.(*reply.MultiBulkReply)
		if !ok {
			return reply.MakeStandardErrReply("error")
		}
		allRes = append(allRes, multiBulkReply.Args...)
	}
	return reply.MakeMultiBulkReply(allRes)
}

// 含义：设置多个键的值
// 用法：MSET key value [key value ...]
// 返回值：始终返回OK
func msetFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) < 3 {
		return reply.MakeArgNumErrReply("mget")
	}
	index := 1
	for index < len(cmdArgs) {
		key := string(cmdArgs[index])
		index++
		if index >= len(cmdArgs) {
			break
		}
		value := string(cmdArgs[index])
		index++
		// 获取执行节点
		peer, _ := cluster.peerPicker.PickNode(key)
		res := cluster.relay(peer, c, utils.ToCmdLine2(_const.CMD_STRING_MSET, []byte(key), []byte(value)))
		if reply.IsErrReply(res) {
			continue
		}
	}
	return reply.MakeOkReply()
}

// 含义：设置多个键的值，当且仅当所有指定的键都不存在时才执行设置操作。
// 用法：MSETNX key value [key value ...]
// 返回值：若所有键设置成功，则返回1；若至少一个键已存在，则返回0。
func msetnxFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) < 3 {
		return reply.MakeArgNumErrReply("mget")
	}
	index := 1
	for index < len(cmdArgs) {
		key := string(cmdArgs[index])
		index++
		index++
		// 获取执行节点
		peer, _ := cluster.peerPicker.PickNode(key)
		res := cluster.relay(peer, c, utils.ToCmdLine2(_const.CMD_KEY_EXISTS, []byte(key)))
		if reply.IsErrReply(res) {
			return reply.MakeIntReply(0)
		}
		if val, ok := res.(*reply.IntReply); ok == false || val.Code == 0 {
			return reply.MakeIntReply(0)
		}
	}
	index = 1
	for index < len(cmdArgs) {
		key := string(cmdArgs[index])
		index++
		if index >= len(cmdArgs) {
			break
		}
		value := string(cmdArgs[index])
		index++
		// 获取执行节点
		peer, _ := cluster.peerPicker.PickNode(key)
		res := cluster.relay(peer, c, utils.ToCmdLine2(_const.CMD_STRING_MSET, []byte(key), []byte(value)))
		if reply.IsErrReply(res) {
			return reply.MakeIntReply(0)
		}
	}
	return reply.MakeIntReply(1)
}

// ==>[list]

// 含义：弹出一个列表最右边的元素，并将它推入另一个列表的最左边。
// 用法：RPOPLPUSH source destination
// 返回值：被弹出的元素。
func rpoppushFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.MakeStandardErrReply("wrong unmber args")
	}
	src := string(cmdArgs[1])
	dest := string(cmdArgs[2])
	srcPeer, _ := cluster.peerPicker.PickNode(src)
	destPeer, _ := cluster.peerPicker.PickNode(dest)
	if srcPeer != destPeer {
		// TODO
		return reply.MakeStandardErrReply("rpoppush must within on peer")
	}

	return cluster.relay(srcPeer, c, cmdArgs)
}

// ==>[hash]

// ==>[set]

// 含义：返回给定所有集合的交集。
// 用法：SINTER key [key ...]
// 返回值：包含交集成员的列表。
func sinterFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	// TODO
	return reply.MakeUnknownErrReply()
}

// 含义：将给定所有集合的交集存储到指定的目标集合中。
// 用法：SINTERSTORE destination key [key ...]
// 返回值：存储到目标集合的元素数量。
func sinterstoreFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	// TODO
	return reply.MakeUnknownErrReply()
}

// 含义：将指定成员从源集合移动到目标集合。
// 用法：SMOVE source destination member
// 返回值：如果成员成功移动，则返回1；如果成员不存在于源集合中，则返回0。
func smoveFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	// TODO

	return reply.MakeUnknownErrReply()
}

// ==>[zset]
