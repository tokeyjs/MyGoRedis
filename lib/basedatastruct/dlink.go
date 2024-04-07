package basedatastruct

import (
	"MyGoRedis/lib/logger"
	"errors"
)

// 双向链表
type BaseDLink struct {
	size int32  // 链表大小
	head *dNode // 头节点
	tail *dNode // 尾节点
}
type dNode struct {
	data     string // 元素值
	next     *dNode // 上一个元素
	previous *dNode //下一个元素
}

// 创建一个双链表
func MakeDLink() *BaseDLink {
	return &BaseDLink{
		size: 0,
		head: nil,
		tail: nil,
	}
}

// 链表是否为空
func (link *BaseDLink) IsEmpty() bool {
	return link.size == 0
}

// 获取长度
func (link *BaseDLink) GetLen() int32 {
	return link.size
}

// 获取头元素
func (link *BaseDLink) GetHeadVal() (string, error) {
	if link.IsEmpty() {
		return "", errors.New("the link is empty.")
	}
	return link.head.data, nil
}

// 获取尾元素
func (link *BaseDLink) GetTailVal() (string, error) {
	if link.IsEmpty() {
		return "", errors.New("the link is empty.")
	}
	return link.tail.data, nil
}

// 添加元素 头插
func (link *BaseDLink) AddElemHead(val string) error {
	if link.IsEmpty() {
		link.head = &dNode{
			data:     val,
			next:     nil,
			previous: nil,
		}
		link.tail = link.head
	} else {
		oldHead := link.head
		link.head = &dNode{
			data:     val,
			next:     oldHead,
			previous: nil,
		}
		oldHead.previous = link.head
	}
	link.size++
	return nil
}

// 添加元素 尾插
func (link *BaseDLink) AddElemTail(val string) error {
	if link.IsEmpty() {
		link.head = &dNode{
			data:     val,
			next:     nil,
			previous: nil,
		}
		link.tail = link.head
	} else {
		oldTail := link.tail
		link.tail = &dNode{
			data:     val,
			next:     nil,
			previous: oldTail,
		}
		oldTail.next = link.tail
	}
	link.size++
	return nil
}

// 删除元素 头
func (link *BaseDLink) DelElemHead() (string, error) {
	var val string
	if link.IsEmpty() {
		return "", errors.New("the link is empty")
	} else {
		// 获取下一个元素
		next := link.head.next
		if next == nil {
			// link只有一个元素
			val = link.head.data
			link.head = nil
			link.tail = nil
		} else {
			val = link.head.data
			next.previous = nil
			link.head = next
		}
	}
	link.size--
	return val, nil
}

// 删除元素 尾
func (link *BaseDLink) DelElemTail() (string, error) {
	var val string
	if link.IsEmpty() {
		return "", errors.New("the link is empty")
	} else {
		// 获取上一个元素
		last := link.tail.previous
		if last == nil {
			// link只有一个元素
			val = link.tail.data
			link.head = nil
			link.tail = nil
		} else {
			val = link.tail.data
			last.next = nil
			link.tail = last
		}
	}
	link.size--
	return val, nil
}

// 插入元素 指定元素(第一个)前
func (link *BaseDLink) InsertElemBefore(val string, secVal string) error {
	node := link.head
	for i := int32(0); i < link.GetLen() && node != nil; i++ {
		if node.data == secVal {
			// 前一个元素
			pre := node.previous
			if pre == nil {
				// 当前元素已经是head --> 进行头插
				return link.AddElemHead(val)
			} else {
				pre.next = &dNode{
					data:     val,
					next:     node,
					previous: pre,
				}
				node.previous = pre.next
				link.size++
				return nil
			}
		} else {
			node = node.next
		}
	}
	return errors.New("not found secVal in link")
}

// 插入元素 指定元素（第一个）后
func (link *BaseDLink) InsertElemAfter(val string, secVal string) error {
	node := link.head
	for i := int32(0); i < link.GetLen() && node != nil; i++ {
		if node.data == secVal {
			// 后一个元素
			last := node.next
			if last == nil {
				// 当前元素已经是tail --> 进行尾插
				return link.AddElemTail(val)
			} else {
				node.next = &dNode{
					data:     val,
					next:     last,
					previous: node,
				}
				last.previous = node.next
				link.size++
				return nil
			}
		} else {
			node = node.next
		}
	}
	return errors.New("not found secVal in link")
}

// 删除指定元素 val 只删除从头搜索到的第一个元素
func (link *BaseDLink) DelElemByVal(val string) error {
	if link.IsEmpty() {
		return errors.New("the link is empty")
	}
	node := link.head
	for i := int32(0); i < link.GetLen() && node != nil; i++ {
		if node.data == val {
			if node.previous == nil {
				// 该元素是head --> 头删除
				_, err := link.DelElemHead()
				return err
			} else if node.next == nil {
				// 该元素是tail --> 尾删除
				_, err := link.DelElemTail()
				return err
			} else {
				// 该元素在中间
				preNode := node.previous
				lastNode := node.next
				preNode.next = lastNode
				lastNode.previous = preNode
				link.size--
				return nil
			}
		} else {
			node = node.next
		}
	}
	return errors.New("not found secVal in link")
}

// 删除指定元素 index
func (link *BaseDLink) DelElemByIndex(index int32) (string, error) {
	if index >= link.size {
		return "", errors.New("Index exceeds array maximum length")
	}
	if index == 0 {
		// 头删
		return link.DelElemHead()
	} else if index == link.size-1 {
		// 尾删
		return link.DelElemTail()
	}
	node := link.head
	for i := int32(0); i < index; i++ {
		node = node.next
	}
	preNode := node.previous
	lastNode := node.next
	preNode.next = lastNode
	lastNode.previous = preNode
	link.size--
	return "", nil
}

// 修改指定索引元素值
func (link *BaseDLink) ModElemByVal(index int32, newVal string) error {
	if index >= link.size {
		return errors.New("Index exceeds array maximum length")
	}
	node := link.head
	for i := int32(0); i < index; i++ {
		node = node.next
	}
	node.data = newVal
	return nil
}

// 获取指定索引元素值
func (link *BaseDLink) GetElemByIndex(index int32) (string, error) {
	if index >= link.size {
		return "", errors.New("Index exceeds array maximum length")
	}
	node := link.head
	for i := int32(0); i < index; i++ {
		node = node.next
	}
	return node.data, nil
}

// -- 测试功能
func showDLink(node *dNode) {
	for node != nil {
		logger.Infof("val[%v]", node.data)
		node = node.next
	}
	logger.Info("")
}

func Test_DLink() {
	link := MakeDLink()
	link.AddElemHead("1")
	link.AddElemHead("2")
	link.AddElemHead("3")
	link.AddElemHead("4")
	link.AddElemHead("5")
	link.AddElemTail("0")
	showDLink(link.head) //5 4 3 2 1 0
	link.DelElemHead()
	showDLink(link.head) // 4 3 2 1 0
	link.DelElemByIndex(2)
	showDLink(link.head) //4 3 1 0
	link.DelElemByVal("3")
	showDLink(link.head) // 4 1 0
}
