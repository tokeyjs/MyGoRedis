package list

import (
	"MyGoRedis/lib/basedatastruct"
)

type List struct {
	data basedatastruct.BaseDLink
}

//通过索引获取列表中的元素

//在列表的元素前或者后插入元素

//获取列表长度

//移出并获取列表的第一个元素

//将一个值插入到列表头部

//移除列表元素（从前面移除count个或从后面移除count个或者移除全部）

//通过索引设置列表元素的值

//移除列表的最后一个元素，返回值为移除的元素

//在列表中添加一个列表尾部

//为列表添加值
