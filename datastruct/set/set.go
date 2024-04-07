package set

import (
	"MyGoRedis/lib/basedatastruct"
)

type Set struct {
	data basedatastruct.BaseSet
}

//向集合添加一个成员

//获取集合的成员数

//返回第一个集合与其他集合之间的差异。

//返回集合的交集

//判断 member元素是否是集合成员

//返回集合中的所有成员

//返回集合中随机个元素

//移除集合中一个成员

//返给定集合的并集
