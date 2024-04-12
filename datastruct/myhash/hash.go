package myhash

import (
	"errors"
	"strconv"
)

type Hash struct {
	data map[string]string // map[mystring]mystring
	size int32             // 字段数量
}

func MakeHash() *Hash {
	return &Hash{
		data: make(map[string]string),
		size: 0,
	}
}

// 删除一个字段
func (h *Hash) DelFiled(field string) int32 {
	if !h.IsExists(field) {
		return 0
	}
	delete(h.data, field)
	h.size--
	return 1
}

// 查看指定的字段是否存在
func (h *Hash) IsExists(field string) bool {
	_, ok := h.data[field]
	return ok
}

// 获取指定字段的值
func (h *Hash) Get(field string) (string, bool) {
	d, ok := h.data[field]
	return d, ok
}

// 获取所有字段和值
func (h *Hash) GetAllKV() []struct {
	Field string
	Value string
} {
	slic := make([]struct {
		Field string
		Value string
	}, 0, h.size)
	for k, v := range h.data {
		slic = append(slic, struct {
			Field string
			Value string
		}{Field: k, Value: v})
	}
	return slic
}

// 为指定字段的整数值加上增量
func (h *Hash) Incr(field string, incr float64) (float64, error) {
	val, ok := h.data[field]
	if !ok {
		return 0, errors.New("not field is " + field)
	}
	// 转成float64
	parseFloat, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}
	// 加上增量
	parseFloat += incr
	// 转为string
	h.data[field] = strconv.FormatFloat(parseFloat, 'g', -1, 64)
	return parseFloat, nil
}

// 获取所有字段
func (h *Hash) GetAllField() []string {
	slic := make([]string, h.size)
	for k, _ := range h.data {
		slic = append(slic, k)
	}
	return slic
}

// 获取字段的数量
func (h *Hash) Len() int32 {
	return h.size
}

// 将字段的值设为 value
func (h *Hash) Modify(field string, value string) {
	h.data[field] = value
}

// 新增字段值及value
func (h *Hash) Set(field string, value string) int32 {
	if h.IsExists(field) {
		h.data[field] = value
		return 0
	} else {
		h.data[field] = value
		h.size++
		return 1
	}

}

// 获取所有值
func (h *Hash) GetAllValue() []string {
	slic := make([]string, h.size)
	for _, v := range h.data {
		slic = append(slic, v)
	}
	return slic
}
