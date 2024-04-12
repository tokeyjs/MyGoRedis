package database

import (
	"MyGoRedis/aof"
	"MyGoRedis/config"
	_const "MyGoRedis/const"
	"MyGoRedis/interface/resp"
	"MyGoRedis/resp/reply"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type StandaloneDatabase struct {
	dbSet      []*DB
	aofHandler *aof.AofHandler
}

func NewStandaloneDatabase() *StandaloneDatabase {
	database := &StandaloneDatabase{}
	if config.Properties.Databases <= 0 {
		config.Properties.Databases = 16
	}
	database.dbSet = make([]*DB, config.Properties.Databases)
	for i := range database.dbSet {
		db := makeDB()
		db.index = i
		database.dbSet[i] = db
	}
	// 初始化aof
	if config.Properties.AppendOnly {
		aofHandler, err := aof.NewAofHandler(database)
		if err != nil {
			panic(err)
		}
		database.aofHandler = aofHandler
		for _, db := range database.dbSet {
			dbIndex := db.index
			db.aofAdd = func(cmd CmdLine) {
				database.aofHandler.AddAof(dbIndex, cmd)
			}
		}
	}

	return database
}

func (database *StandaloneDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("panic: %v", err)
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	// 判断该连接是否进行认证
	if !client.IsCertification() {
		if cmdName == _const.CMD_CONN_AUTH {
			// 进行认证
			if len(args) != 2 {
				return reply.MakeArgNumErrReply("auth")
			}
			return execAUTH(client, database, args[1:])
		} else {
			// 返回错误
			return reply.MakeStandardErrReply("NOAUTH Authentication required.")
		}
	}
	if cmdName == _const.CMD_CONN_SELECT {
		if len(args) != 2 {
			return reply.MakeArgNumErrReply("select")
		}
		return execSelect(client, database, args[1:])
	}
	return database.dbSet[client.GetDBIndex()].Exec(client, args)
}

func (database *StandaloneDatabase) Close() {

}

func (database *StandaloneDatabase) AfterClientClose(c resp.Connection) {

}

func execSelect(c resp.Connection, database *StandaloneDatabase, args [][]byte) resp.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeStandardErrReply("ERR invalid DB index")
	}
	if dbIndex >= len(database.dbSet) {
		return reply.MakeStandardErrReply("ERR DB index is out of range")
	}
	c.SelectDB(dbIndex)
	return reply.MakeOkReply()
}

func execAUTH(c resp.Connection, database *StandaloneDatabase, args [][]byte) resp.Reply {
	password := string(args[0])
	c.CheckAuth(password)
	if c.IsCertification() {
		return reply.MakeOkReply()
	}
	return reply.MakeStandardErrReply("password error")
}
