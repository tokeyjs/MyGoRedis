package basedatastruct

import (
	"MyGoRedis/lib/logger"
	"errors"
	"sync"
	"sync/atomic"
)

// 并发安全双向链表
type BaseDLink struct {
	size   atomic.Int32 // 链表大小
	head   *dNode       // 头节点
	tail   *dNode       // 尾节点
	rmutex sync.RWMutex
}
type dNode struct {
	data string // 元素值
	next *dNode // 上一个元素
	prev *dNode //下一个元素
}

// 创建一个双链表
func MakeDLink() *BaseDLink {
	return &BaseDLink{
		head: nil,
		tail: nil,
	}
}

// 链表是否为空
func (link *BaseDLink) IsEmpty() bool {
	return link.size.Load() == 0
}

// 获取长度
func (link *BaseDLink) Size() int32 {
	return link.size.Load()
}

// 获取头元素
func (link *BaseDLink) GetHeadVal() (string, error) {
	link.rmutex.RLock()
	defer link.rmutex.RUnlock()
	return link.getHeadVal()
}
func (link *BaseDLink) getHeadVal() (string, error) {
	if link.IsEmpty() {
		return "", errors.New("the link is empty.")
	}
	return link.head.data, nil
}

// 获取尾元素
func (link *BaseDLink) GetTailVal() (string, error) {
	link.rmutex.RLock()
	defer link.rmutex.RUnlock()
	return link.getTailVal()
}
func (link *BaseDLink) getTailVal() (string, error) {
	if link.IsEmpty() {
		return "", errors.New("the link is empty.")
	}
	return link.tail.data, nil
}

// 添加元素 头插
func (link *BaseDLink) AddElemHead(val string) error {
	link.rmutex.Lock()
	defer link.rmutex.Unlock()
	return link.addElemHead(val)
}
func (link *BaseDLink) addElemHead(val string) error {
	if link.IsEmpty() {
		link.head = &dNode{
			data: val,
			next: nil,
			prev: nil,
		}
		link.tail = link.head
	} else {
		oldHead := link.head
		link.head = &dNode{
			data: val,
			next: oldHead,
			prev: nil,
		}
		oldHead.prev = link.head
	}
	link.size.Add(1)
	return nil
}

// 添加元素 尾插
func (link *BaseDLink) AddElemTail(val string) error {
	link.rmutex.Lock()
	defer link.rmutex.Unlock()
	return link.addElemTail(val)
}
func (link *BaseDLink) addElemTail(val string) error {
	if link.IsEmpty() {
		link.head = &dNode{
			data: val,
			next: nil,
			prev: nil,
		}
		link.tail = link.head
	} else {
		oldTail := link.tail
		link.tail = &dNode{
			data: val,
			next: nil,
			prev: oldTail,
		}
		oldTail.next = link.tail
	}
	link.size.Add(1)
	return nil
}

// 删除元素 头
func (link *BaseDLink) DelElemHead() (string, error) {
	link.rmutex.Lock()
	defer link.rmutex.Unlock()
	return link.delElemHead()
}
func (link *BaseDLink) delElemHead() (string, error) {
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
			next.prev = nil
			link.head = next
		}
	}
	link.size.Add(-1)
	return val, nil
}

// 删除元素 尾
func (link *BaseDLink) DelElemTail() (string, error) {
	link.rmutex.Lock()
	defer link.rmutex.Unlock()
	return link.delElemTail()
}
func (link *BaseDLink) delElemTail() (string, error) {
	var val string
	if link.IsEmpty() {
		return "", errors.New("the link is empty")
	} else {
		// 获取上一个元素
		last := link.tail.prev
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
	link.size.Add(-1)
	return val, nil
}

// 插入元素 指定元素(第一个)前
func (link *BaseDLink) InsertElemBefore(val string, secVal string) error {
	link.rmutex.Lock()
	defer link.rmutex.Unlock()
	node := link.head
	for i := int32(0); i < link.Size() && node != nil; i++ {
		if node.data == secVal {
			// 前一个元素
			pre := node.prev
			if pre == nil {
				// 当前元素已经是head --> 进行头插
				return link.addElemHead(val)
			} else {
				pre.next = &dNode{
					data: val,
					next: node,
					prev: pre,
				}
				node.prev = pre.next
				link.size.Add(1)
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
	link.rmutex.Lock()
	defer link.rmutex.Unlock()
	node := link.head
	for i := int32(0); i < link.Size() && node != nil; i++ {
		if node.data == secVal {
			// 后一个元素
			last := node.next
			if last == nil {
				// 当前元素已经是tail --> 进行尾插
				return link.addElemTail(val)
			} else {
				node.next = &dNode{
					data: val,
					next: last,
					prev: node,
				}
				last.prev = node.next
				link.size.Add(1)
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
	link.rmutex.Lock()
	defer link.rmutex.Unlock()
	if link.IsEmpty() {
		return errors.New("the link is empty")
	}
	node := link.head
	for i := int32(0); i < link.Size() && node != nil; i++ {
		if node.data == val {
			if node.prev == nil {
				// 该元素是head --> 头删除
				_, err := link.delElemHead()
				return err
			} else if node.next == nil {
				// 该元素是tail --> 尾删除
				_, err := link.delElemTail()
				return err
			} else {
				// 该元素在中间
				preNode := node.prev
				lastNode := node.next
				preNode.next = lastNode
				lastNode.prev = preNode
				link.size.Add(-1)
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
	link.rmutex.Lock()
	defer link.rmutex.Unlock()
	if index >= link.Size() {
		return "", errors.New("Index exceeds array maximum length")
	}
	if index == 0 {
		// 头删
		return link.delElemHead()
	} else if index == link.Size()-1 {
		// 尾删
		return link.delElemTail()
	}
	node := link.head
	for i := int32(0); i < index; i++ {
		node = node.next
	}
	preNode := node.prev
	lastNode := node.next
	preNode.next = lastNode
	lastNode.prev = preNode
	link.size.Add(-1)
	return "", nil
}

// 修改指定索引元素值
func (link *BaseDLink) ModElemByVal(index int32, newVal string) error {
	link.rmutex.Lock()
	defer link.rmutex.Unlock()
	if index >= link.Size() {
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
	link.rmutex.RLock()
	defer link.rmutex.RUnlock()
	if index >= link.Size() {
		return "", errors.New("Index exceeds array maximum length")
	}
	node := link.head
	for i := int32(0); i < index; i++ {
		node = node.next
	}
	return node.data, nil
}

// 从头删除 count个 value元素
func (link *BaseDLink) DelFromBegin(count int32, value string) int32 {
	link.rmutex.Lock()
	defer link.rmutex.Unlock()
	argCount := count
	node := link.head
	for node != nil && count > 0 {
		if node.data == value {
			if node.next == node.prev && node.next == nil {
				// 只有一个元素了
				link.head = nil
				link.tail = nil
				count--
				break
			} else if node.prev == nil {
				// 处在头部
				link.head = node.next
				node = node.next
				node.prev = nil
				count--
				continue
			} else if node.next == nil {
				// 处在尾部
				link.tail = node.prev
				link.tail.next = nil
				count--
				break
			} else {
				node = node.next
				node.prev = node.prev.prev
				node.prev.next = node
				count--

			}
		}
		node = node.next
	}
	link.size.Add(count - argCount)
	return argCount - count
}

// 从尾部删除count个 value元素
func (link *BaseDLink) DelFromBack(count int32, value string) int32 {
	link.rmutex.Lock()
	defer link.rmutex.Unlock()
	argCount := count
	node := link.tail
	for node != nil && count > 0 {
		if node.data == value {
			// 只有一个元素了
			if node.prev == node.next && node.next == nil {
				count--
				link.head = nil
				link.tail = nil
				break
			} else if node.next == nil {
				// 是尾
				node = node.prev
				link.tail = node
				node.next = nil
				count--
				continue
			} else if node.prev == nil {
				// 是头
				count--
				link.head = node.next
				link.head.prev = nil
				break
			} else {
				// 既不是头也不是尾
				node = node.prev
				node.next = node.next.next
				node.next.prev = node
				count--
			}
		}
		node = node.prev
	}
	link.size.Add(count - argCount)
	return argCount - count
}

// 删除所有元素
func (link *BaseDLink) Clear() {
	link.rmutex.Lock()
	defer link.rmutex.Unlock()
	link.clear()
}
func (link *BaseDLink) clear() {
	link.head = nil
	link.tail = nil
	link.size.Store(0)
}

// 获取从start 到 stop的元素
func (link *BaseDLink) GetRange(start, stop int32) []string {
	if stop < 0 {
		stop += link.Size()
	}
	if start > stop {
		return nil
	}
	link.rmutex.RLock()
	defer link.rmutex.RUnlock()
	slic := make([]string, 0)
	index := int32(0)
	node := link.head
	for index <= stop {
		if index >= start {
			slic = append(slic, node.data)
		}
		node = node.next
		index++
	}
	return slic
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
