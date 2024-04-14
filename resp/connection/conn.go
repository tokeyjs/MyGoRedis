package connection

import (
	"MyGoRedis/config"
	"MyGoRedis/lib/sync/wait"
	"net"
	"sync"
	"time"
)

type Connection struct {
	conn            net.Conn
	waitingReply    wait.Wait
	mutex           sync.Mutex
	selectDB        int   // 目前操作的db
	isCF            bool  //是否进行认证
	connTimeStamp   int64 //连接时的时间戳(s)
	activeTimeStamp int64 //上次收到客户端命令时间戳(s)
	clusterConn     bool  // 标记集群间通信连接
}

// 设置该连接为集群中各个节点通信的连接
func (c *Connection) SetClusterConn() {
	c.clusterConn = true
}

// 是否为集群的内部通信连接
func (c *Connection) IsClusterClient() bool {
	return c.clusterConn
}

func (c *Connection) UpdateConn() {
	c.activeTimeStamp = time.Now().Unix()
}

func (c *Connection) IsCertification() bool {
	return c.isCF
}

func (c *Connection) CheckAuth(password string) {
	if len(config.Properties.RequirePass) <= 0 {
		return
	}
	if config.Properties.RequirePass == password {
		c.isCF = true
	}
}

func (c *Connection) GetAge() int32 {
	return int32(time.Now().Unix() - c.connTimeStamp)
}

func (c *Connection) GetIdle() int32 {
	return int32(time.Now().Unix() - c.activeTimeStamp)
}

func (c *Connection) IsTimeOut() bool {
	if config.Properties.ClientTimeOutSec == 0 {
		return false
	}
	return c.GetIdle() > int32(config.Properties.ClientTimeOutSec)
}

func NewConn(conn net.Conn) *Connection {
	return &Connection{
		conn:            conn,
		isCF:            len(config.Properties.RequirePass) <= 0,
		connTimeStamp:   time.Now().Unix(),
		activeTimeStamp: time.Now().Unix(),
		clusterConn:     false,
	}
}

// 创建一个伪连接
func NewFakeConn() *Connection {
	return &Connection{
		conn:            nil,
		isCF:            true,
		connTimeStamp:   time.Now().Unix(),
		activeTimeStamp: time.Now().Unix(),
		clusterConn:     true,
	}
}

func (c *Connection) Close() error {
	c.waitingReply.WaitWithTimeout(10 * time.Second)
	_ = c.conn.Close()
	return nil
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) Write(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}
	c.mutex.Lock()
	c.waitingReply.Add(1)
	defer func() {
		c.waitingReply.Done()
		c.mutex.Unlock()
	}()
	_, err := c.conn.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) GetDBIndex() int {
	return c.selectDB
}

func (c *Connection) SelectDB(dbIndex int) {
	c.selectDB = dbIndex
}
