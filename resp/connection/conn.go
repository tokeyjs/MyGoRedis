package connection

import (
	"MyGoRedis/lib/sync/wait"
	"net"
	"sync"
	"time"
)

type Connection struct {
	conn         net.Conn
	waitingReply wait.Wait
	mutex        sync.Mutex
	selectDB     int
}

func NewConn(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
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
