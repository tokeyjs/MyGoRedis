package basedatastruct

type BaseSet struct {
	data map[string]struct{}
	size int32
}

func MakeSet() *BaseSet {
	return &BaseSet{
		data: make(map[string]struct{}),
		size: 0,
	}
}

// 插入
func (s *BaseSet) Insert(val string) {
	s.data[val] = struct{}{}
	s.size++
}

// 删除
func (s *BaseSet) Remove(val string) {
	delete(s.data, val)
	s.size--
}

// 清空
func (s *BaseSet) Clear() {
	s.data = make(map[string]struct{})
	s.size = 0
}

// 是否存在
func (s *BaseSet) IsExists(val string) bool {
	_, ok := s.data[val]
	return ok
}

// 元素数量
func (s *BaseSet) Size() int32 {
	return s.size
}
