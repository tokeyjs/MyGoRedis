package cluster

import (
	"MyGoRedis/interface/resp"
	"MyGoRedis/resp/reply"
	"github.com/sirupsen/logrus"
)

func makeRouter() map[string]CmdFunc {
	routerMap := make(map[string]CmdFunc)
	routerMap["exists"] = defaultFunc
	routerMap["type"] = defaultFunc
	routerMap["set"] = defaultFunc
	routerMap["setnx"] = defaultFunc
	routerMap["get"] = defaultFunc
	routerMap["getset"] = defaultFunc
	routerMap["ping"] = pingFunc
	routerMap["rename"] = renameFunc
	routerMap["renamenx"] = renameFunc
	routerMap["flushdb"] = flushdbFunc
	routerMap["del"] = delFunc
	routerMap["select"] = selectFunc

	return routerMap
}

// GET SET
func defaultFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	key := string(cmdArgs[0])
	peer := cluster.peerPicker.PickNode(key)
	logrus.Infof("cmd run at peer:[%v]\n", peer)
	return cluster.relay(peer, c, cmdArgs)
}

// ping
func pingFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	return cluster.db.Exec(c, cmdArgs)
}

// select
func selectFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	return cluster.db.Exec(c, cmdArgs)
}

// rename renamenx
func renameFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.MakeStandardErrReply("wrong unmber args")
	}
	src := string(cmdArgs[1])
	dest := string(cmdArgs[2])
	// TODO: 优化让处在不同节点也能成功
	srcPeer := cluster.peerPicker.PickNode(src)
	destPeer := cluster.peerPicker.PickNode(dest)
	if srcPeer != destPeer {
		return reply.MakeStandardErrReply("rename must within on peer")
	}

	return cluster.relay(srcPeer, c, cmdArgs)
}

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

// del k1 k2 k3
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
