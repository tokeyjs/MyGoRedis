package cluster

import (
	_const "MyGoRedis/const"
	"MyGoRedis/interface/resp"
	"MyGoRedis/lib/logger"
	"MyGoRedis/lib/utils"
	"MyGoRedis/resp/client"
	"MyGoRedis/resp/reply"
	"context"
	"errors"
	"strconv"
)

// 从连接池中获取一个连接
func (cluster *ClusterDatabase) getPeerClient(peer string) (*client.Client, error) {
	pool, ok := cluster.peerConnectionPool[peer]
	if !ok {
		return nil, errors.New("connection pool not found")
	}
	object, err := pool.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}
	c, ok := object.(*client.Client)
	if !ok {
		return nil, errors.New("wrong type")
	}
	return c, err
}

// 将一个客户端连接归还至连接池
func (cluster *ClusterDatabase) returnPeerClient(peer string, peerClient *client.Client) error {
	pool, ok := cluster.peerConnectionPool[peer]
	if !ok {
		return errors.New("connection pool not found")
	}
	return pool.ReturnObject(context.Background(), peerClient)
}

// 转发命令到集群指定节点中执行
func (cluster *ClusterDatabase) relay(peer string, c resp.Connection, args [][]byte) resp.Reply {
	logger.Infof("cmd[%v] run at peer:[%v]", string(args[0]), peer)
	if peer == cluster.self {
		return cluster.db.Exec(c, args)
	}
	peerClient, err := cluster.getPeerClient(peer)
	if err != nil {
		// 获取客户端连接错误，将该节点下线
		cluster.peerPicker.RemoveNode(peer)
		cluster.nodes.Delete(peer)
		logger.Infof("node[%v] disconnection.", peer)
		return reply.MakeStandardErrReply(err.Error())
	}
	defer func() {
		_ = cluster.returnPeerClient(peer, peerClient)
	}()
	peerClient.Send(utils.ToCmdLine(_const.CMD_CONN_SELECT, strconv.Itoa(c.GetDBIndex())))
	return peerClient.Send(args)
}

// 广播
func (cluster *ClusterDatabase) broadcast(c resp.Connection, args [][]byte) map[string]resp.Reply {
	results := make(map[string]resp.Reply)
	// 如果当前连接为集群的通信连接，则此命令是转发来的命令不需要进行广播
	if c.IsClusterClient() {
		results[cluster.self] = cluster.relay(cluster.self, c, args)
		return results
	}
	cluster.nodes.Range(func(key, value any) bool {
		results[key.(string)] = cluster.relay(key.(string), c, args)
		return true
	})
	return results
}
