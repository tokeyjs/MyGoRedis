package basedatastruct

import (
	"sync"
	"sync/atomic"
)

// 并发安全set
type BaseSet struct {
	data sync.Map //map[mystring]struct{}
	size atomic.Int32
}

func MakeBaseSet() *BaseSet {
	return &BaseSet{}
}

// 插入
func (s *BaseSet) Insert(val string) {
	s.data.Store(val, struct{}{})
	s.size.Add(1)
}

// 删除
func (s *BaseSet) Remove(val string) {
	s.data.Delete(val)
	s.size.Add(-1)
}

// 清空
func (s *BaseSet) Clear() {
	s.data = sync.Map{}
	s.size.Store(0)
}

// 是否存在
func (s *BaseSet) IsExists(val string) bool {
	_, ok := s.data.Load(val)
	return ok
}

// 元素数量
func (s *BaseSet) Size() int32 {
	return s.size.Load()
}

// 随机n个获取元素
func (s *BaseSet) GetRandom(n int32) []string {
	slic := make([]string, 0, n)
	for i := int32(0); i < n; i++ {
		s.data.Range(func(key, value any) bool {
			slic = append(slic, key.(string))
			return false
		})
	}
	return slic
}

// 返回所有成员
func (s *BaseSet) GetAll() []string {
	slic := make([]string, 0, s.Size())
	s.data.Range(func(key, value any) bool {
		slic = append(slic, key.(string))
		return false
	})
	return slic
}

//todo: 集合操作
