package chper

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
)

/*
layout example:

Time1: there are 3 nodes：
			N2#1		N3#1
		------2----->-----3--------------
		|								|
		|								|
N1#1	1								4 	N2#1
		|								|
		|								|
		------6-----------5--------------
			N3#1		N1#1

index range:	1			2			3			4			5			6			1
virtual Node:	(	N2#1	](	N3#1	](	N2#1	](	N1#1	](	N3#1	](	N1#1	]
Real Node:		(	N2		](	N3		](	N2		](	N1		](	N3		](	N1		]


Time2: one node be deleted, there are 2 nodes：
						N3#1
		------+----->-----3--------------
		|								|
		|								|
N1#1	1								+
		|								|
		|								|
		------6-----------5--------------
			N3#1		N1#1

index range:	1			2			3			4			5			6			1
virtual Node:	(				N3#1	](				N1#1	](	N3#1	](	N1#1	]
Real Node:		(				N3		](				N1		](	N3		](	N1		]
*/

// CHash is a consistent hashing
// more information see http://en.wikipedia.org/wiki/Consistent_hashing
type CHash[Node any] struct {
	// realNodesIndex2Nodes is index -> real node id -> real node
	// two nodes may be has the same index
	realNodesIndex2Nodes map[uint32]map[string]realNode[Node]

	virtualNodeMap  map[uint32]*virtualNode[Node]
	virtualNodeList []*virtualNode[Node]

	option *chashOption[Node]

	lock sync.RWMutex
}

type virtualNode[Node any] struct {
	beginIndex uint32

	realNode realNode[Node]
}

type realNode[Node any] struct {
	id                string
	index             uint32
	virtualNodeIndexs map[uint32]bool

	node Node
}

type chashOption[Node any] struct {
	nodeIDer func(Node) ([]byte, error)
	indexer  func(data []byte) uint32

	virtualNodeFactor int
}

func (cho *chashOption[Node]) adaptVirtualNodeFactor(nodeSize int) {
	if cho.virtualNodeFactor > 0 {
		return
	}

	factor := 1500 / nodeSize
	if factor < 5 {
		factor = 5
	}

	cho.virtualNodeFactor = factor
}

func defaultCHashOption[Node any]() *chashOption[Node] {
	return &chashOption[Node]{
		indexer: crc32.ChecksumIEEE,
		nodeIDer: func(n Node) ([]byte, error) {
			return json.Marshal(n)
		},
	}
}

type chashOptionFunc[Node any] func(*chashOption[Node])

func CHashOptionNodeIDer[Node any](nodeIDer func(node Node) ([]byte, error)) chashOptionFunc[Node] {
	return func(co *chashOption[Node]) {
		co.nodeIDer = nodeIDer
	}
}

func CHashOptionVirtualNodeFactor[Node any](virtualNodeFactor int) chashOptionFunc[Node] {
	return func(co *chashOption[Node]) {
		co.virtualNodeFactor = virtualNodeFactor
	}
}

func CHashOptionIndexer[Node any](indexer func(data []byte) uint32) chashOptionFunc[Node] {
	return func(co *chashOption[Node]) {
		co.indexer = indexer
	}
}

func NewCHash[Node any](nodes []Node, options ...chashOptionFunc[Node]) (*CHash[Node], error) {
	if len(nodes) == 0 {
		return nil, fmt.Errorf("want at least one node")
	}

	option := defaultCHashOption[Node]()
	for _, f := range options {
		f(option)
	}

	option.adaptVirtualNodeFactor(len(nodes))

	ch := &CHash[Node]{
		realNodesIndex2Nodes: make(map[uint32]map[string]realNode[Node], len(nodes)),
		virtualNodeMap:       make(map[uint32]*virtualNode[Node], len(nodes)*option.virtualNodeFactor),
		option:               option,
	}

	for _, node := range nodes {
		err := ch.addNode(node)
		if err != nil {
			return nil, err
		}
	}

	return ch, nil
}

func (ch *CHash[Node]) AddNode(node Node) error {
	ch.lock.Lock()
	err := ch.addNode(node)
	ch.lock.Unlock()

	return err
}

func (ch *CHash[Node]) addNode(node Node) error {
	realNodeID, err := ch.option.nodeIDer(node)
	if err != nil {
		return fmt.Errorf("nodeIDer fail, err : %w", err)
	}
	realNodeIDStr := fmt.Sprintf("%x", realNodeID)
	realNodeIndex := ch.option.indexer(realNodeID)
	realNodes, ok := ch.realNodesIndex2Nodes[realNodeIndex]
	if !ok {
		realNodes = map[string]realNode[Node]{}
		ch.realNodesIndex2Nodes[realNodeIndex] = realNodes
	}
	if _, ok := realNodes[realNodeIDStr]; ok {
		return fmt.Errorf("node existed, id: %s", realNodeIDStr)
	}

	rn := realNode[Node]{
		id:                realNodeIDStr,
		index:             realNodeIndex,
		virtualNodeIndexs: make(map[uint32]bool, ch.option.virtualNodeFactor),

		node: node,
	}

	succCount := 0
	for i := ch.option.virtualNodeFactor * len(realNodes); succCount < ch.option.virtualNodeFactor; i++ {
		key := virtualNodeKey(realNodeIDStr, i)
		index := ch.option.indexer(key)
		if _, ok := ch.virtualNodeMap[index]; ok {
			continue
		}

		succCount++

		ch.virtualNodeMap[index] = &virtualNode[Node]{
			beginIndex: index,
			realNode:   rn,
		}

		rn.virtualNodeIndexs[index] = true
	}

	ch.sortvirtualNode()

	realNodes[realNodeIDStr] = rn

	return nil
}

func virtualNodeKey(nodeID string, i int) []byte {
	return []byte(fmt.Sprintf("%s#constinctethashing#%d", nodeID, i))
}

func (ch *CHash[Node]) RemoveNode(node Node) error {
	ch.lock.Lock()
	err := ch.removeNode(node)
	ch.lock.Unlock()

	return err
}

func (ch *CHash[Node]) removeNode(node Node) error {
	realNodeID, err := ch.option.nodeIDer(node)
	if err != nil {
		return fmt.Errorf("nodeIDer fail, err : %w", err)
	}
	realNodeIDStr := fmt.Sprintf("%x", realNodeID)
	realNodeIndex := ch.option.indexer(realNodeID)
	realNodes, ok := ch.realNodesIndex2Nodes[realNodeIndex]
	if !ok {
		return fmt.Errorf("node not exist, id: %s", realNodeIDStr)
	}
	realNode, ok := realNodes[realNodeIDStr]
	if !ok {
		return fmt.Errorf("node not exist, id: %s", realNodeIDStr)
	}
	delete(realNodes, realNodeIDStr)

	for index := range realNode.virtualNodeIndexs {
		delete(ch.virtualNodeMap, index)
	}
	ch.sortvirtualNode()

	return nil
}

func (ch *CHash[Node]) Hash(data []byte) (Node, error) {
	ch.lock.RLock()
	node, err := ch.hash(data)
	ch.lock.RUnlock()

	return node, err
}

func (ch *CHash[Node]) hash(data []byte) (node Node, err error) {
	if len(ch.virtualNodeList) == 0 {
		err = fmt.Errorf("zero node")
		return
	}

	return ch.find(ch.option.indexer(data)), nil
}

type IDer interface {
	ID() ([]byte, error)
}

func (ch *CHash[Node]) HashIDer(ider IDer) (node Node, err error) {
	data, err := ider.ID()
	if err != nil {
		err = fmt.Errorf("get ID fail, err: %w", err)
		return
	}

	return ch.Hash(data)
}

func (ch *CHash[Node]) find(index uint32) (node Node) {
	i := sort.Search(len(ch.virtualNodeList), func(i int) bool {
		return ch.virtualNodeList[i].beginIndex > index
	})
	if i == len(ch.virtualNodeList) {
		i = 0
	}
	return ch.virtualNodeList[i].realNode.node
}

func (ch *CHash[Node]) sortvirtualNode() {
	list := make([]*virtualNode[Node], 0, len(ch.virtualNodeMap))
	for _, vn := range ch.virtualNodeMap {
		list = append(list, vn)
	}
	sort.Sort(virtualNodeSlice[Node](list))

	ch.virtualNodeList = list
}

type virtualNodeSlice[Node any] []*virtualNode[Node]

func (x virtualNodeSlice[Node]) Len() int           { return len(x) }
func (x virtualNodeSlice[Node]) Less(i, j int) bool { return x[i].beginIndex < x[j].beginIndex }
func (x virtualNodeSlice[Node]) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
