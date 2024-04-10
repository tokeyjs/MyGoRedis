package myset

import (
	"MyGoRedis/lib/basedatastruct"
)

type Set struct {
	data *basedatastruct.BaseSet
}

func MakeSet() *Set {
	return &Set{
		data: basedatastruct.MakeBaseSet(),
	}
}

// 向集合添加一个成员
func (s *Set) Add(data string) {
	s.data.Insert(data)
}

// 获取集合的成员数
func (s *Set) Size() int32 {
	return s.data.Size()
}

// 判断 member元素是否是集合成员
func (s *Set) IsExists(data string) bool {
	return s.data.IsExists(data)
}

// 返回集合中的所有成员
func (s *Set) GetAll() []string {
	return s.data.GetAll()
}

// 返回集合中随机个元素
func (s *Set) GetRandom(n int32) []string {
	if n <= 0 {
		return nil
	}
	return s.data.GetRandom(n)
}

// 移除集合中一个成员
func (s *Set) Remove(data string) {
	s.data.Remove(data)
}

// todo :集合操作
//返回第一个集合与其他集合之间的差异

//返回集合的交集

// 返给定集合的并集
