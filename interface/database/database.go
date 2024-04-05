package database

import "MyGoRedis/interface/resp"

type CmdLine = [][]byte

type DataBase interface {
	Exec(client resp.Connection, args [][]byte) resp.Reply
	Close()
	AfterClientClose(c resp.Connection)
}

type DataEntity struct {
	Data interface{}
}
