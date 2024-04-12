package mystring

import (
	"strconv"
	"sync"
)

type String struct {
	rmutex sync.RWMutex //读写锁
	data   string       //存储数据
}

func MakeString() *String {
	return &String{data: ""}
}

// 设置值
func (s *String) Set(val string) {
	s.rmutex.Lock()
	defer s.rmutex.Unlock()
	s.data = val
}

// 获取值
func (s *String) Get() string {
	s.rmutex.RLock()
	defer s.rmutex.RUnlock()
	return s.data
}

// 转换成数字
func (s *String) ToNumber() (float64, error) {
	s.rmutex.RLock()
	defer s.rmutex.RUnlock()
	return s.toNumber()
}
func (s *String) toNumber() (float64, error) {
	parseFloat, err := strconv.ParseFloat(s.data, 64)
	if err != nil {
		return 0, err
	}
	return parseFloat, nil
}

// 替换值并返回旧值
func (s *String) Modify(val string) string {
	s.rmutex.Lock()
	defer s.rmutex.Unlock()
	ret := s.data
	s.data = val
	return ret
}

// 返回储存的字符串值的长度
func (s *String) Len() int32 {
	s.rmutex.RLock()
	defer s.rmutex.RUnlock()
	return int32(len(s.data))
}

// 储存的数字值增一（如果不能转换成数字就返回错误）
func (s *String) Incr() (float64, error) {
	s.rmutex.Lock()
	defer s.rmutex.Unlock()
	n, err := s.toNumber()
	if err != nil {
		return 0, err
	}
	n++
	s.data = strconv.FormatFloat(n, 'g', -1, 64)
	return n, nil
}

// 所储存的值加上给定的增量值（如果不能转换成数字就返回错误）
func (s *String) IncrNum(num float64) (float64, error) {
	s.rmutex.Lock()
	defer s.rmutex.Unlock()
	n, err := s.toNumber()
	if err != nil {
		return 0, err
	}
	n += num
	s.data = strconv.FormatFloat(n, 'g', -1, 64)
	return n, nil
}

// 储存的数字值减一（如果不能转换成数字就返回错误）
func (s *String) Decr() (float64, error) {
	s.rmutex.Lock()
	defer s.rmutex.Unlock()
	n, err := s.toNumber()
	if err != nil {
		return 0, err
	}
	n--
	s.data = strconv.FormatFloat(n, 'g', -1, 64)
	return n, nil
}

// 储存的值减去给定的减量值（如果不能转换成数字就返回错误）
func (s *String) DecrNum(num float64) (float64, error) {
	s.rmutex.Lock()
	defer s.rmutex.Unlock()
	n, err := s.toNumber()
	if err != nil {
		return 0, err
	}
	n -= num
	s.data = strconv.FormatFloat(n, 'g', -1, 64)
	return n, nil
}

// 将字符串追加到末尾
func (s *String) AppendStr(str string) int32 {
	s.rmutex.Lock()
	defer s.rmutex.Unlock()
	s.data = s.data + str
	return int32(len(s.data))
}
