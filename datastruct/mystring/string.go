package mystring

import "strconv"

type String struct {
	data string //存储数据
}

func MakeString() *String {
	return &String{data: ""}
}

// 设置值
func (s *String) Set(val string) {
	s.data = val
}

// 获取值
func (s *String) Get() string {
	return s.data
}

// 转换成数字
func (s *String) ToNumber() (float64, error) {
	parseFloat, err := strconv.ParseFloat(s.data, 64)
	if err != nil {
		return 0, err
	}
	return parseFloat, nil
}

// 替换值并返回旧值
func (s *String) Modify(val string) string {
	ret := s.data
	s.data = val
	return ret
}

// 返回储存的字符串值的长度
func (s *String) Len() int32 {
	return int32(len(s.data))
}

// 储存的数字值增一（如果不能转换成数字就返回错误）
func (s *String) Incr() error {
	n, err := s.ToNumber()
	if err != nil {
		return err
	}
	n++
	s.Set(strconv.FormatFloat(n, 'g', -1, 64))
	return nil
}

// 所储存的值加上给定的增量值（如果不能转换成数字就返回错误）
func (s *String) IncrNum(num float64) error {
	n, err := s.ToNumber()
	if err != nil {
		return err
	}
	n += num
	s.Set(strconv.FormatFloat(n, 'g', -1, 64))
	return nil
}

// 储存的数字值减一（如果不能转换成数字就返回错误）
func (s *String) Decr() error {
	n, err := s.ToNumber()
	if err != nil {
		return err
	}
	n--
	s.Set(strconv.FormatFloat(n, 'g', -1, 64))
	return nil
}

// 储存的值减去给定的减量值（如果不能转换成数字就返回错误）
func (s *String) DecrNum(num float64) error {
	n, err := s.ToNumber()
	if err != nil {
		return err
	}
	n -= num
	s.Set(strconv.FormatFloat(n, 'g', -1, 64))
	return nil
}

// 将字符串追加到末尾
func (s *String) AppendStr(str string) {
	s.data = s.data + str
}
