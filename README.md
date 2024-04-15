# MyGoRedis
一个使用golang实现的redis

## 1.项目简介

本项目使用Go语言实现了一个分布式Redis，支持RESP协议的解析，实现String、List（双向链表实现）、Hash、Set、ZSet（跳跃表实现）多种数据结构。此外，项目还支持Key的有效期，实现了过期删除策略（懒惰删除和定期删除相结合），并采用了集群模式（一致性Hash算法），实现了数据在不同节点间的分布式存储。在持久化方面，本项目实现了AOF持久化，保证了数据的可靠性和一致性。


## 2.项目亮点及成果

采用Go语言实现，高并发、简洁易读的特点；支持多种数据结构，满足不同业务场景的需求；实现了Key的有效期和过期删除策略；采用一致性Hash算法实现集群模式，实现了数据在不同节点间的分布式存储，提高了系统的扩展性和可用性；实现了AOF持久化，保证了数据的可靠性和一致性。

## 3.已实现命令

### *==key==*
* DEL
* EXISTS
* EXPIRE
* EXPIREAT
* KEYS
* PERSIST
* PEXPIRE
* PEXPIREAT
* PTTL
* RANDOMKEY
* RENAME
* RENAMENX
* TTL
* TYPE
### *==string==*
* APPEND
* DECR
* DECRBY
* GET
* GETRANGE
* GETSET
* INCR
* INCRBY
* INCRBYFLOAT
* MGET
* MSET
* MSETNX
* PSETEX
* SET
* SETEX
* SETNX
* STRLEN
### *==list==*
* LINDEX
* LINSERT
* LLEN
* LPOP
* LPUSH
* LPUSHX
* LREM
* LSET
* RPOP
* RPOPLPUSH
* RPUSH
* RPUSHX
* LRANGE
### *==hash==*
* HDEL
* HEXISTS
* HGET
* HGETALL
* HINCRBY
* HINCRBYFLOAT
* HKEYS
* HLEN
* HMGET
* HMSET
* HSET
* HSETNX
* HVALS
### *==set==*
* SADD
* SCARD
* SISMEMBER
* SMEMBERS
* SMOVE
* SPOP
* SREM
* SRANDMEMBER
### *==zset==*
* ZADD
* ZCARD
* ZCOUNT
* ZINCRBY
* ZRANGE
* ZRANGEBYSCORE
* ZRANK
* ZREM
* ZREMRANGEBYRANK
* ZREMRANGEBYSCORE
* ZREVRANGE
* ZREVRANGEBYSCORE
* ZREVRANK
* ZSCORE
### *==else==*
* TIME
* FLUSHDB
* PING
* AUTH
* ECHO
* SELECT

## 其他
其他命令及Redis功能待后续进行完善...

## 特别感谢
《Redis设计与实现》

https://github.com/HDT3213/godis