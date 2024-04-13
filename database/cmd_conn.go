package database

import (
	_const "MyGoRedis/const"
	"MyGoRedis/interface/resp"
	"MyGoRedis/resp/reply"
)

// 检查完成

// 实现命令
// 含义：用于检查与服务器的连接是否仍然活动。当收到PING命令时，服务器将返回一个PONG响应。
// 用法：PING
// 返回值：如果服务器正常运行，则返回PONG；如果提供了可选的消息参数，则返回该消息。
func exec_CONN_PING(db *DB, args [][]byte) resp.Reply {
	return reply.MakePongReply()
}

// 含义：用于对Redis服务器进行身份验证。通常在连接Redis服务器后，客户端需要使用AUTH命令提供的密码来进行身份验证。
// 用法：AUTH password
// 返回值：如果密码正确，服务器将返回OK，表示身份验证成功；如果密码错误，服务器将返回一个错误。
//func exec_CONN_AUTH(db *DB, args [][]byte) resp.Reply {
//	// TODO: 其他命令得被拦截
//	password := string(args[0])
//	if config.Properties.RequirePass == "" || password == config.Properties.RequirePass {
//		// 认证成功
//		return reply.MakeOkReply()
//	}
//	return reply.MakeStandardErrReply("password error")
//}

// 含义：用于在不执行任何操作的情况下返回给定的字符串。
// 用法：ECHO message
// 返回值：返回输入的消息字符串。
func exec_CONN_ECHO(db *DB, args [][]byte) resp.Reply {
	return reply.MakeBulkReply(args[0])
}

// 含义：用于选择指定的数据库。Redis服务器支持多个数据库，默认情况下有16个数据库，编号从0到15。
// 用法：SELECT index
// 返回值：如果指定的数据库存在，则返回OK；如果指定的数据库索引超出范围，则返回一个错误。
//func exec_CONN_SELECT(db *DB, args [][]byte) resp.Reply {
//	return reply.MakeUnknownErrReply()
//}
// select以实现===>standalone_database.go

func init() {
	// 注册
	RegisterCommand(_const.CMD_CONN_PING, exec_CONN_PING, 1)
	//RegisterCommand(_const.CMD_CONN_AUTH, exec_CONN_AUTH, 2)
	RegisterCommand(_const.CMD_CONN_ECHO, exec_CONN_ECHO, 2)
	//RegisterCommand(_const.CMD_CONN_SELECT, exec_CONN_SELECT, 2)
}
