package mylist

import (
	"MyGoRedis/lib/basedatastruct"
)

type List struct {
	data *basedatastruct.BaseDLink
}

func MakeList() *List {
	return &List{
		data: basedatastruct.MakeDLink(),
	}
}

// 通过索引获取列表中的元素
func (list *List) GetByIndex(index int32) (string, error) {
	return list.data.GetElemByIndex(index)
}

// 在列表的元素前或者后插入元素
func (list *List) InsertByValue(value string, secVal string, isFront bool) error {
	if isFront {
		return list.data.InsertElemBefore(value, secVal)
	}
	return list.data.InsertElemAfter(value, secVal)
}

// 获取列表长度
func (list *List) Size() int32 {
	return list.data.Size()
}

// 移出并获取列表的第一个元素
func (list *List) PopBegin() (string, error) {
	return list.data.DelElemHead()
}

// 将一个值插入到列表头部
func (list *List) PushBegin(value string) error {
	return list.data.AddElemHead(value)
}

// 移除列表元素（从前面移除count个或从后面移除-count个或者移除全部）
func (list *List) Remove(count int32, value string) int32 {
	if count == 0 {
		// 移除全部value
		return list.data.DelFromBegin(list.Size(), value)
	} else if count > 0 {
		// 从前面开始移除
		return list.data.DelFromBegin(count, value)
	} else {
		// 从后面开始移除
		return list.data.DelFromBack(-count, value)
	}
}

// 通过索引设置列表元素的值
func (list *List) SetByIndex(index int32, value string) error {
	return list.data.ModElemByVal(index, value)
}

// 移除列表的最后一个元素，返回值为移除的元素
func (list *List) PopBack() (string, error) {
	return list.data.DelElemTail()
}

// 在列表中添加一个列表尾部
func (list *List) PushBack(value string) {
	_ = list.data.AddElemTail(value)
}

// 获取区间元素
func (list *List) GetRange(start, stop int32) []string {
	return list.data.GetRange(start, stop)
}
