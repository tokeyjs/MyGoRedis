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
	connNum    atomic.Int32 // 连接数量
}

func MakeHandler() *RespHandler {
	var db databaseface.DataBase
	if len(config.Properties.Self) > 0 && len(config.Properties.Peers) > 0 {
		// 启动集群版redis
		db = cluster.MakeClusterDatabase()
		logger.Infof("mygoredis 集群版启动 bind:[%v:%v] peers:[%v]...", config.Properties.Bind, config.Properties.Port, config.Properties.Peers)
	} else {
		// 单机版redis
		db = database.NewStandaloneDatabase()
		logger.Infof("mygoredis 单机版启动 bind:[%v:%v]...", config.Properties.Bind, config.Properties.Port)
	}
	return &RespHandler{
		db: db,
	}
}

func (r *RespHandler) closeClient(client *connection.Connection) {
	r.connNum.Add(-1)
	_ = client.Close()
	r.db.AfterClientClose(client)
	r.activeConn.Delete(client)
}

func (r *RespHandler) clientNum() int32 {
	return r.connNum.Load()
}

func (r *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if r.closing.Load() || (config.Properties.MaxClients != 0 && r.clientNum() > int32(config.Properties.MaxClients)) {
		_ = conn.Close()
		return
	}
	r.connNum.Add(1)
	client := connection.NewConn(conn)
	r.activeConn.Store(client, struct{}{})
	ch := parser.ParseStream(conn)
	for payload := range ch {
		// error
		if payload.Err != nil {
			if payload.Err == io.EOF || payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of close network connection") {
				r.closeClient(client)
				logger.Infof("connection closed: %v", client.RemoteAddr().String())
				return
			}
			// protocol error
			errReply := reply.MakeStandardErrReply(payload.Err.Error())
			err := client.Write(errReply.ToBytes())
			if err != nil {
				r.closeClient(client)
				logger.Infof("connection closed: %v", client.RemoteAddr().String())
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
			logger.Error("require multi bulk reply")
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
