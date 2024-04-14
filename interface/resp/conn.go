package resp

type Connection interface {
	Write([]byte) error
	GetDBIndex() int
	SelectDB(int)
	IsCertification() bool     //-->是否认证
	CheckAuth(password string) //-->输入密码进行认证，改变认证状态
	GetAge() int32             //-->获取连接时长s
	GetIdle() int32            //-->获取空闲时长s
	IsTimeOut() bool           //-->判断客户端是否超时（空闲连接清理）
	UpdateConn()               // 更新连接最新活动时间
	IsClusterClient() bool     // 是否为集群间的内部连接
	SetClusterConn()           // 标记该连接为集群内部连接
}
