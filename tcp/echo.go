package tcp

import (
	"MyGoRedis/lib/logger"
	"MyGoRedis/lib/sync/wait"
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// 客户端
type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (e *EchoClient) Close() error {
	e.Waiting.WaitWithTimeout(10 * time.Second)
	_ = e.Conn.Close()
	return nil
}

type EchoHandler struct {
	activeConn sync.Map    // 记录所有连接
	closing    atomic.Bool // 此handler是否关闭
}

func MakeHandler() *EchoHandler {
	return &EchoHandler{}
}

func (handler *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if handler.closing.Load() {
		_ = conn.Close()
		return
	}
	client := &EchoClient{
		Conn: conn,
	}
	// 将新客户端连接塞入
	handler.activeConn.Store(client, struct{}{})

	// 处理
	reader := bufio.NewReader(conn)
	for true {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Infof("client[%v] disconnection :%v", conn.RemoteAddr().String(), err)
			} else {
				logger.Warnf("client[%v] error: %v", conn.RemoteAddr().String(), err)
			}
			_ = client.Close()
			handler.activeConn.Delete(client)
			return
		}
		// 处理业务
		client.Waiting.Add(1)
		logger.Infof("client[%v]:%v", client.Conn.RemoteAddr().String(), msg)
		_, _ = client.Conn.Write([]byte(msg))
		client.Waiting.Done()
	}

}

func (handler *EchoHandler) Close() error {
	logger.Infof("handler shutting down")
	handler.closing.Store(true)
	// 将连接依次关闭
	handler.activeConn.Range(func(key, value any) bool {
		conn, _ := key.(*EchoClient)
		_ = conn.Close()
		return true
	})
	return nil
}
