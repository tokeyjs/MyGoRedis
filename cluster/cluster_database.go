package cluster

import (
	"MyGoRedis/config"
	_const "MyGoRedis/const"
	"MyGoRedis/database"
	databaseface "MyGoRedis/interface/database"
	"MyGoRedis/interface/resp"
	"MyGoRedis/lib/consistenthash"
	"MyGoRedis/lib/logger"
	"MyGoRedis/resp/reply"
	"context"
	pool "github.com/jolestar/go-commons-pool/v2"
	"strings"
)

type ClusterDatabase struct {
	self               string
	nodes              []string
	peerPicker         *consistenthash.NodeMap
	peerConnectionPool map[string]*pool.ObjectPool // 与其他节点的连接池
	db                 databaseface.DataBase
}

type CmdFunc func(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply

var router = makeRouter()

func MakeClusterDatabase() *ClusterDatabase {
	cluster := &ClusterDatabase{
		self:               config.Properties.Self,
		nodes:              make([]string, 0, len(config.Properties.Peers)+1),
		peerPicker:         consistenthash.NewNodeMap(),
		peerConnectionPool: make(map[string]*pool.ObjectPool),
		db:                 database.NewStandaloneDatabase(),
	}
	for _, peer := range config.Properties.Peers {
		cluster.nodes = append(cluster.nodes, peer)
	}
	cluster.nodes = append(cluster.nodes, config.Properties.Self)
	cluster.peerPicker.AddNode(cluster.nodes...)
	ctx := context.Background()
	for _, peer := range config.Properties.Peers {
		cluster.peerConnectionPool[peer] = pool.NewObjectPoolWithDefaultConfig(ctx, &connectionFactory{
			Peer: peer,
		})
	}
	return cluster
}

func (cDB *ClusterDatabase) Exec(client resp.Connection, args [][]byte) (result resp.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("panic: %v", err)
			result = reply.MakeUnknownErrReply()
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	// 判断该连接是否进行认证
	if !client.IsCertification() {
		if cmdName == _const.CMD_CONN_AUTH {
			// 进行认证
			if len(args) != 2 {
				return reply.MakeArgNumErrReply("auth")
			}
			return cDB.execAUTH(client, args[1:])
		} else {
			// 返回错误
			return reply.MakeStandardErrReply("NOAUTH Authentication required.")
		}
	}
	cmdFunc, ok := router[cmdName]
	if !ok {
		result = reply.MakeStandardErrReply("not supported cmd")
		return
	}
	result = cmdFunc(cDB, client, args)
	return
}

func (cDB *ClusterDatabase) Close() {
	cDB.db.Close()
}

func (cDB *ClusterDatabase) AfterClientClose(c resp.Connection) {
	cDB.db.AfterClientClose(c)
}

func (cDB *ClusterDatabase) execAUTH(c resp.Connection, args [][]byte) resp.Reply {
	password := string(args[0])
	c.CheckAuth(password)
	if c.IsCertification() {
		logger.Debugf("auth successful")
		return reply.MakeOkReply()
	}
	logger.Debugf("auth failed")
	return reply.MakeStandardErrReply("password error")
}
