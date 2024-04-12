package tcp

import (
	"MyGoRedis/interface/tcp"
	"context"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Address string
}

func ListenAndServeWithSignal(
	cfg *Config,
	handler tcp.Handler) error {
	// 关闭tcp
	closeChan := make(chan struct{})
	signalChan := make(chan os.Signal)
	// 退出
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-signalChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()
	// 创建TCP监听
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logrus.Infof("start listen %v ...", cfg.Address)
	ListenAndServe(listener, handler, closeChan)
	return nil
}

func ListenAndServe(
	listener net.Listener,
	handler tcp.Handler,
	closeChan <-chan struct{}) {
	go func() {
		<-closeChan
		_ = listener.Close()
		_ = handler.Close()
	}()
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()
	// 创建上下文
	ctx := context.Background()
	var waitDone sync.WaitGroup
	for true {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Errorf("accept error: %v", err)
			break
		}
		logrus.Infof("new client[%v] connection.", conn.RemoteAddr().String())
		// 处理新客户端
		waitDone.Add(1)
		go func() {
			defer func() { waitDone.Done() }()
			handler.Handle(ctx, conn)
		}()
	}
	// 等待所有业务处理完成
	waitDone.Wait()
}
