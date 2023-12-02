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
	// realNodeMap is name -> real node
	realNodeMap map[string]realNode[Node]

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
	name              string
	virtualNodeIndexs map[uint32]bool

	node Node
}

type chashOption[Node any] struct {
	nodeNaming func(Node) (string, error)
	indexer    func(data []byte) uint32

	virtualNodeFactor int

	weightSpecify func(node Node) int
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
		nodeNaming: func(n Node) (string, error) {
			bs, err := json.Marshal(n)
			return string(bs), err
		},
		weightSpecify: func(node Node) int {
			return 1
		},
	}
}

type chashOptionFunc[Node any] func(*chashOption[Node])

func CHashOptionNodeNaming[Node any](nodeNaming func(node Node) (string, error)) chashOptionFunc[Node] {
	return func(co *chashOption[Node]) {
		co.nodeNaming = nodeNaming
	}
}

func CHashOptionWeightSpecify[Node any](weightSpecify func(node Node) int) chashOptionFunc[Node] {
	return func(co *chashOption[Node]) {
		co.weightSpecify = weightSpecify
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
		realNodeMap:    make(map[string]realNode[Node], len(nodes)),
		virtualNodeMap: make(map[uint32]*virtualNode[Node], len(nodes)*option.virtualNodeFactor),
		option:         option,
	}

	for _, node := range nodes {
		err := ch.addNode(node, option.weightSpecify(node), false)
		if err != nil {
			return nil, err
		}
	}

	ch.sortVirtualNode()
	return ch, nil
}

func (ch *CHash[Node]) AddNodeWithWeight(node Node, weight int) error {
	ch.lock.Lock()
	err := ch.addNode(node, weight, true)
	ch.lock.Unlock()

	return err
}

func (ch *CHash[Node]) AddNode(node Node) error {
	ch.lock.Lock()
	err := ch.addNode(node, ch.option.weightSpecify(node), true)
	ch.lock.Unlock()

	return err
}

func (ch *CHash[Node]) addNode(node Node, weight int, doSort bool) error {
	if weight < 1 {
		return fmt.Errorf("weight must be greater than zero")
	}

	realNodeName, err := ch.option.nodeNaming(node)
	if err != nil {
		return fmt.Errorf("nodeNaming fail, err : %w", err)
	}
	_, ok := ch.realNodeMap[realNodeName]
	if ok {
		return fmt.Errorf("node existed, name: %s", realNodeName)
	}

	rn := realNode[Node]{
		name:              realNodeName,
		virtualNodeIndexs: make(map[uint32]bool, ch.option.virtualNodeFactor),

		node: node,
	}

	succCount := 0
	for i := 0; succCount < ch.option.virtualNodeFactor*weight; i++ {
		key := virtualNodeKey(realNodeName, i)
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

	if doSort {
		ch.sortVirtualNode()
	}

	ch.realNodeMap[realNodeName] = rn

	return nil
}

func virtualNodeKey(nodeName string, i int) []byte {
	return []byte(fmt.Sprintf("%s#constinctethashing#%d", nodeName, i))
}

func (ch *CHash[Node]) RemoveNode(node Node) error {
	ch.lock.Lock()
	err := ch.removeNode(node)
	ch.lock.Unlock()

	return err
}

func (ch *CHash[Node]) removeNode(node Node) error {
	realNodeName, err := ch.option.nodeNaming(node)
	if err != nil {
		return fmt.Errorf("nodeNaming fail, err : %w", err)
	}
	realNode, ok := ch.realNodeMap[realNodeName]
	if !ok {
		return fmt.Errorf("node not exist, name: %s", realNodeName)
	}
	delete(ch.realNodeMap, realNodeName)

	for index := range realNode.virtualNodeIndexs {
		delete(ch.virtualNodeMap, index)
	}
	ch.sortVirtualNode()

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

func (ch *CHash[Node]) find(index uint32) (node Node) {
	i := sort.Search(len(ch.virtualNodeList), func(i int) bool {
		return ch.virtualNodeList[i].beginIndex > index
	})
	if i == len(ch.virtualNodeList) {
		i = 0
	}
	return ch.virtualNodeList[i].realNode.node
}

func (ch *CHash[Node]) sortVirtualNode() {
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
