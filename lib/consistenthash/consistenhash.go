package consistenthash

import "github.com/serialx/hashring"

type NodeMap struct {
	nodes []string           //记录所有节点
	har   *hashring.HashRing // 一致性hash环
}

func NewNodeMap() *NodeMap {
	return &NodeMap{
		nodes: make([]string, 0),
		har:   hashring.New(nil),
	}
}

func (m *NodeMap) IsEmpty() bool {
	return len(m.nodes) == 0
}

// 添加节点
func (m *NodeMap) AddNode(nodes ...string) {
	for _, node := range nodes {
		m.har = m.har.AddNode(node)
		m.nodes = append(m.nodes, node)
	}
}

// 查找该key的数据应该去哪个node中
func (m *NodeMap) PickNode(key string) (string, bool) {
	return m.har.GetNode(key)
}
