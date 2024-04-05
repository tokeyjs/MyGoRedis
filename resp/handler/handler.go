package handler

import (
	"MyGoRedis/cluster"
	"MyGoRedis/config"
	"MyGoRedis/database"
	databaseface "MyGoRedis/interface/database"
	"MyGoRedis/lib/logger"
	"MyGoRedis/resp/connection"
	"MyGoRedis/resp/parser"
	"MyGoRedis/resp/reply"
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
)

type RespHandler struct {
	activeConn sync.Map
	db         databaseface.DataBase
	closing    atomic.Bool
}

func MakeHandler() *RespHandler {
	var db databaseface.DataBase
	if len(config.Properties.Self) > 0 && len(config.Properties.Peers) > 0 {
		// 启动集群版redis
		db = cluster.MakeClusterDatabase()
		logrus.Infof("redis 集群版启动...")
	} else {
		// 单机版redis
		db = database.NewStandaloneDatabase()
		logrus.Infof("redis 单机版启动...")
	}
	return &RespHandler{
		db: db,
	}
}

func (r *RespHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	r.db.AfterClientClose(client)
	r.activeConn.Delete(client)
}

func (r *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if r.closing.Load() {
		_ = conn.Close()
		return
	}
	client := connection.NewConn(conn)
	r.activeConn.Store(client, struct{}{})
	ch := parser.ParseStream(conn)
	for payload := range ch {
		// error
		if payload.Err != nil {
			if payload.Err == io.EOF || payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of close network connection") {
				r.closeClient(client)
				logrus.Infof("connection closed: %v", client.RemoteAddr().String())
				return
			}
			// protocol error
			errReply := reply.MakeStandardErrReply(payload.Err.Error())
			err := client.Write(errReply.ToBytes())
			if err != nil {
				r.closeClient(client)
				logrus.Infof("connection closed: %v", client.RemoteAddr().String())
				return
			}
			continue
		}

		// exec
		if payload.Data == nil {
			continue
		}
		reply_, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logrus.Error("require multi bulk reply")
			continue
		}
		result := r.db.Exec(client, reply_.Args)
		if result != nil {
			_ = client.Write(result.ToBytes())
		} else {
			_ = client.Write(reply.MakeUnknownErrReply().ToBytes())
		}
	}
}

func (r *RespHandler) Close() error {
	logger.Info("handler shutting down")
	r.closing.Store(true)
	r.activeConn.Range(func(key interface{}, value interface{}) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		return true
	})
	r.db.Close()
	return nil
}
