package consistenthash

import (
	"github.com/serialx/hashring"
)

type NodeMap struct {
	size int32              //记录所有节点 map[string]struct{}
	har  *hashring.HashRing // 一致性hash环
}

func NewNodeMap() *NodeMap {
	return &NodeMap{
		size: 0,
		har:  hashring.New(nil),
	}
}

func (m *NodeMap) IsEmpty() bool {
	return m.size == 0
}

// 添加节点
func (m *NodeMap) AddNode(nodes ...string) {
	for _, node := range nodes {
		m.har = m.har.AddNode(node)
		m.size++
	}
}

// 删除节点
func (m *NodeMap) RemoveNode(node string) {
	m.har = m.har.RemoveNode(node)
	m.size--
}

// 查找该key的数据应该去哪个node中
func (m *NodeMap) PickNode(key string) (string, bool) {
	return m.har.GetNode(key)
}
