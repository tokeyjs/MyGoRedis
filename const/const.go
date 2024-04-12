package _const

import (
	"MyGoRedis/datastruct/myhash"
	"MyGoRedis/datastruct/mylist"
	"MyGoRedis/datastruct/myset"
	"MyGoRedis/datastruct/mystring"
	"MyGoRedis/datastruct/myzset"
)

// 命令(不包括redis全部命令)
const (
	//server
	CMD_SERVER_TIME    = "time"
	CMD_SERVER_FLUSHDB = "flushdb"

	// connection
	CMD_CONN_PING   = "ping"
	CMD_CONN_AUTH   = "auth"
	CMD_CONN_ECHO   = "echo"
	CMD_CONN_SELECT = "select"

	// key
	CMD_KEY_DEL       = "del"
	CMD_KEY_EXISTS    = "exists"
	CMD_KEY_EXPIRE    = "expire"
	CMD_KEY_EXPIREAT  = "expireat"
	CMD_KEY_KEYS      = "keys"
	CMD_KEY_MOVE      = "move"
	CMD_KEY_PERSIST   = "persist"
	CMD_KEY_PEXPIRE   = "pexpire"
	CMD_KEY_PEXPIREAT = "pexpireat"
	CMD_KEY_PTTL      = "pttl"
	CMD_KEY_RANDOMKEY = "randomkey"
	CMD_KEY_RENAME    = "rename"
	CMD_KEY_RENAMENX  = "renamenx"
	CMD_KEY_TTL       = "ttl"
	CMD_KEY_TYPE      = "type"

	// string
	CMD_STRING_APPEND      = "append"
	CMD_STRING_DECR        = "decr"
	CMD_STRING_DECRBY      = "decrby"
	CMD_STRING_GET         = "get"
	CMD_STRING_GETRANGE    = "getrange"
	CMD_STRING_GETSET      = "getset"
	CMD_STRING_INCR        = "incr"
	CMD_STRING_INCRBY      = "incrby"
	CMD_STRING_INCRBYFLOAT = "incrbyfloat"
	CMD_STRING_MGET        = "mget"
	CMD_STRING_MSET        = "mset"
	CMD_STRING_MSETNX      = "msetnx"
	CMD_STRING_PSETEX      = "psetex"
	CMD_STRING_SET         = "set"
	CMD_STRING_SETEX       = "setex"
	CMD_STRING_SETNX       = "setnx"
	CMD_STRING_STRLEN      = "strlen"

	// list
	CMD_LIST_BLPOP      = "blpop"
	CMD_LIST_BRPOP      = "brpop"
	CMD_LIST_BRPOPLPUSH = "brpoplpush"
	CMD_LIST_LINDEX     = "lindex"
	CMD_LIST_LINSERT    = "linsert"
	CMD_LIST_LLEN       = "llen"
	CMD_LIST_LPOP       = "lpop"
	CMD_LIST_LPUSH      = "lpush"
	CMD_LIST_LPUSHX     = "lpushx"
	CMD_LIST_LREM       = "lrem"
	CMD_LIST_LSET       = "lset"
	CMD_LIST_RPOP       = "rpop"
	CMD_LIST_RPOPLPUSH  = "rpoplpush"
	CMD_LIST_RPUSH      = "rpush"
	CMD_LIST_RPUSHX     = "rpushx"

	// hash
	CMD_HASH_HDEL         = "hdel"
	CMD_HASH_HEXISTS      = "hexists"
	CMD_HASH_HGET         = "hget"
	CMD_HASH_HGETALL      = "hgetall"
	CMD_HASH_HINCRBY      = "hincrby"
	CMD_HASH_HINCRBYFLOAT = "hincrbyfloat"
	CMD_HASH_HKEYS        = "hkeys"
	CMD_HASH_HLEN         = "hlen"
	CMD_HASH_HMGET        = "hmget"
	CMD_HASH_HMSET        = "hmset"
	CMD_HASH_HSET         = "hset"
	CMD_HASH_HSETNX       = "hsetnx"
	CMD_HASH_HVALS        = "hvals"

	// set
	CMD_SET_SADD        = "sadd"
	CMD_SET_SCARD       = "scard"
	CMD_SET_SINTER      = "sinter"
	CMD_SET_SINTERSTORE = "sinterstore"
	CMD_SET_SISMEMBER   = "sismember"
	CMD_SET_SMEMBERS    = "smembers"
	CMD_SET_SMOVE       = "smove"
	CMD_SET_SPOP        = "spop"
	CMD_SET_SRANDMEMBER = "srandmember"

	// zset
	CMD_ZSET_ZADD             = "zadd"
	CMD_ZSET_ZCARD            = "zcard"
	CMD_ZSET_ZCOUNT           = "zcount"
	CMD_ZSET_ZINCRBY          = "zincrby"
	CMD_ZSET_ZRANGE           = "zrange"
	CMD_ZSET_ZRANGEBYSCORE    = "zrangebyscore"
	CMD_ZSET_ZRANK            = "zrank"
	CMD_ZSET_ZREM             = "zrem"
	CMD_ZSET_ZREMRANGEBYRANK  = "zremrangebyrank"
	CMD_ZSET_ZREMRANGEBYSCORE = "zremrangebyscore"
	CMD_ZSET_ZREVRANGE        = "zrevrange"
	CMD_ZSET_ZREVRANGEBYSCORE = "zrevrangebyscore"
	CMD_ZSET_ZREVRANK         = "zrevrank"
	CMD_ZSET_ZSCORE           = "zscore"
)

// value数据类型
const (
	TYPE_DATA_STRING = "string"
	TYPE_DATA_LIST   = "list"
	TYPE_DATA_HASH   = "hash"
	TYPE_DATA_SET    = "set"
	TYPE_DATA_ZSET   = "zset"
	TYPE_DATA_NONE   = "nil"
	TYPE_DATA_ANY    = "any" //以上几种类型都有可能
)

// interface ===> 转string, list, hash, set, zset

func DataToSTRING(data interface{}) *mystring.String {
	res, ok := data.(*mystring.String)
	if ok {
		return res
	}
	return nil
}

func DataToLIST(data interface{}) *mylist.List {
	res, ok := data.(*mylist.List)
	if ok {
		return res
	}
	return nil
}

func DataToHash(data interface{}) *myhash.Hash {
	res, ok := data.(*myhash.Hash)
	if ok {
		return res
	}
	return nil
}

func DataToSET(data interface{}) *myset.Set {
	res, ok := data.(*myset.Set)
	if ok {
		return res
	}
	return nil
}

func DataToZSET(data interface{}) *myzset.ZSet {
	res, ok := data.(*myzset.ZSet)
	if ok {
		return res
	}
	return nil
}
