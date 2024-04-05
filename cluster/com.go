package cluster

import (
	"MyGoRedis/interface/resp"
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
	if peer == cluster.self {
		return cluster.db.Exec(c, args)
	}
	peerClient, err := cluster.getPeerClient(peer)
	if err != nil {
		return reply.MakeStandardErrReply(err.Error())
	}
	defer func() {
		_ = cluster.returnPeerClient(peer, peerClient)
	}()
	peerClient.Send(utils.ToCmdLine("select", strconv.Itoa(c.GetDBIndex())))
	return peerClient.Send(args)
}

// 广播
func (cluster ClusterDatabase) broadcast(c resp.Connection, args [][]byte) map[string]resp.Reply {
	results := make(map[string]resp.Reply)
	for _, node := range cluster.nodes {
		results[node] = cluster.relay(node, c, args)
	}
	return results
}
