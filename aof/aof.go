package aof

import (
	"MyGoRedis/config"
	databaseface "MyGoRedis/interface/database"
	"MyGoRedis/lib/utils"
	"MyGoRedis/resp/connection"
	"MyGoRedis/resp/parser"
	"MyGoRedis/resp/reply"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
)

var aofBufferSize = 65535

type payload struct {
	cmdLine databaseface.CmdLine
	dbIndex int
}

type AofHandler struct {
	database    databaseface.DataBase
	aofChan     chan *payload
	aofFile     *os.File
	aofFileName string
	currentDB   int // 记录上一次操作的db index
}

// NewAofHandler
func NewAofHandler(database databaseface.DataBase) (*AofHandler, error) {
	handler := &AofHandler{
		database:    database,
		currentDB:   0,
		aofFileName: config.Properties.AppendFilename,
	}
	handler.loadAof()
	aofFile, err := os.OpenFile(handler.aofFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handler.aofFile = aofFile
	handler.aofChan = make(chan *payload, aofBufferSize)
	go handler.handleAof()
	return handler, nil
}

// add
func (handler *AofHandler) AddAof(dbIndex int, cmd databaseface.CmdLine) {
	if config.Properties.AppendOnly && handler.aofChan != nil {
		handler.aofChan <- &payload{
			cmdLine: cmd,
			dbIndex: dbIndex,
		}
	}
}

// handler aof
func (handler *AofHandler) handleAof() {
	for p := range handler.aofChan {
		if p.dbIndex != handler.currentDB {
			selectDBCmdData := reply.MakeMultiBulkReply(utils.ToCmdLine("select", strconv.Itoa(p.dbIndex))).ToBytes()
			if _, err := handler.aofFile.Write(selectDBCmdData); err != nil {
				logrus.Errorf("aof insert selectdbcmd [%v] error:%v", p.dbIndex, err)
				continue
			}
			handler.currentDB = p.dbIndex
		}
		if _, err := handler.aofFile.Write(reply.MakeMultiBulkReply(p.cmdLine).ToBytes()); err != nil {
			logrus.Errorf("aof write cmd error; srcCmd:[%v]  error:[%v]", string(reply.MakeMultiBulkReply(p.cmdLine).ToBytes()), err)
			continue
		}

	}
}

// 加载aof文件中的数据并加载到数据库中
func (handler *AofHandler) loadAof() {
	file, err := os.Open(handler.aofFileName)
	if err != nil {
		logrus.Errorf("open aof error:%v", err)
		return
	}
	defer file.Close()
	fakeConn := &connection.Connection{}
	ch := parser.ParseStream(file)
	for p := range ch {
		if p.Err != nil {
			if p.Err == io.EOF {
				break
			}
			logrus.Errorln(p.Err)
			continue
		}
		if p.Data == nil {
			logrus.Errorf("data is nil\n")
			continue
		}
		r, ok := p.Data.(*reply.MultiBulkReply)
		if !ok {
			logrus.Errorf("data err: %v\n", string(p.Data.ToBytes()))
			continue
		}
		rep := handler.database.Exec(fakeConn, r.Args)
		if reply.IsErrReply(rep) {
			logrus.Errorf("exec err: %v\n", string(rep.ToBytes()))
		}
	}
}
