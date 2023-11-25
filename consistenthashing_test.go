package chper

import (
	"crypto/rand"
	"fmt"
	"hash/crc32"
	mrand "math/rand"
	"reflect"
	"testing"
)

type Node struct {
	Name string
}

var (
	nodeA = &Node{Name: "A"}
	nodeB = &Node{Name: "B"}
	nodeC = &Node{Name: "C"}
	nodeD = &Node{Name: "D"}
)

func TestCHash(t *testing.T) {
	ch, err := NewCHash[*Node]([]*Node{nodeA, nodeB, nodeC},
		CHashOptionIndexer[*Node](func(data []byte) uint32 { return crc32.ChecksumIEEE(data) }),
		CHashOptionNodeIDer[*Node](func(node *Node) ([]byte, error) { return []byte(node.Name), nil }),
		CHashOptionVirtualNodeFactor[*Node](5),
	)

	if err != nil {
		t.Errorf("want nil, got: %v", err)
		return
	}

	// ABC
	{
		for i := 0; i < 5; i++ {
			for _, cas := range []struct {
				name     string
				data     string
				wantNode *Node
			}{
				{name: "case1", data: "1", wantNode: nodeC},
				{name: "case2", data: "2", wantNode: nodeA},
				{name: "case3", data: "3", wantNode: nodeB},
				{name: "case4", data: "4", wantNode: nodeC},
				{name: "case5", data: "5", wantNode: nodeC},
				{name: "case6", data: "6", wantNode: nodeA},
				{name: "case7", data: "7", wantNode: nodeB},
				{name: "case8", data: "8", wantNode: nodeB},
				{name: "case9", data: "9", wantNode: nodeC},
				{name: "case10", data: "10", wantNode: nodeA},
				{name: "case11", data: "11", wantNode: nodeA},
				{name: "case12", data: "12", wantNode: nodeC},
			} {
				got, err := ch.Hash([]byte(cas.data))
				if err != nil {
					t.Errorf("%s, want nil, got: %v", cas.name, err)
				}
				if !reflect.DeepEqual(got, cas.wantNode) {
					t.Errorf("%s, want: %v, got: %v", cas.name, cas.wantNode, got)
				}
			}
		}
	}

	// AB
	{
		err = ch.RemoveNode(nodeD)
		if err == nil {
			t.Errorf("want err, got nil")
		}
		for _, cas := range []struct {
			name     string
			data     string
			wantNode *Node
		}{
			{name: "case1", data: "1", wantNode: nodeC},
			{name: "case2", data: "2", wantNode: nodeA},
			{name: "case3", data: "3", wantNode: nodeB},
			{name: "case4", data: "4", wantNode: nodeC},
			{name: "case5", data: "5", wantNode: nodeC},
			{name: "case6", data: "6", wantNode: nodeA},
			{name: "case7", data: "7", wantNode: nodeB},
			{name: "case8", data: "8", wantNode: nodeB},
			{name: "case9", data: "9", wantNode: nodeC},
			{name: "case10", data: "10", wantNode: nodeA},
			{name: "case11", data: "11", wantNode: nodeA},
			{name: "case12", data: "12", wantNode: nodeC},
		} {
			got, err := ch.Hash([]byte(cas.data))
			if err != nil {
				t.Errorf("%s, want nil, got: %v", cas.name, err)
			}
			if !reflect.DeepEqual(got, cas.wantNode) {
				t.Errorf("%s, want: %v, got: %v", cas.name, cas.wantNode, got)
			}
		}

		err = ch.RemoveNode(nodeC)
		if err != nil {
			t.Errorf("want nil, got: %v", err)
		}
		for i := 0; i < 5; i++ {
			for _, cas := range []struct {
				name     string
				data     string
				wantNode *Node
			}{
				{name: "case1", data: "1", wantNode: nodeB},
				{name: "case2", data: "2", wantNode: nodeA},
				{name: "case3", data: "3", wantNode: nodeB},
				{name: "case4", data: "4", wantNode: nodeB},
				{name: "case5", data: "5", wantNode: nodeB},
				{name: "case6", data: "6", wantNode: nodeA},
				{name: "case7", data: "7", wantNode: nodeB},
				{name: "case8", data: "8", wantNode: nodeB},
				{name: "case9", data: "9", wantNode: nodeB},
				{name: "case10", data: "10", wantNode: nodeA},
				{name: "case11", data: "11", wantNode: nodeA},
				{name: "case12", data: "12", wantNode: nodeB},
			} {
				got, err := ch.Hash([]byte(cas.data))
				if err != nil {
					t.Errorf("%s, want nil, got: %v", cas.name, err)
				}
				if !reflect.DeepEqual(got, cas.wantNode) {
					t.Errorf("%s, want: %v, got: %v", cas.name, cas.wantNode, got)
				}
			}
		}
	}

	// ABC
	{
		err = ch.AddNode(nodeC)
		if err != nil {
			t.Errorf("want nil, got: %v", err)
		}
		for _, cas := range []struct {
			name     string
			data     string
			wantNode *Node
		}{
			{name: "case1", data: "1", wantNode: nodeC},
			{name: "case2", data: "2", wantNode: nodeA},
			{name: "case3", data: "3", wantNode: nodeB},
			{name: "case4", data: "4", wantNode: nodeC},
			{name: "case5", data: "5", wantNode: nodeC},
			{name: "case6", data: "6", wantNode: nodeA},
			{name: "case7", data: "7", wantNode: nodeB},
			{name: "case8", data: "8", wantNode: nodeB},
			{name: "case9", data: "9", wantNode: nodeC},
			{name: "case10", data: "10", wantNode: nodeA},
			{name: "case11", data: "11", wantNode: nodeA},
			{name: "case12", data: "12", wantNode: nodeC},
		} {
			got, err := ch.Hash([]byte(cas.data))
			if err != nil {
				t.Errorf("%s, want nil, got: %v", cas.name, err)
			}
			if !reflect.DeepEqual(got, cas.wantNode) {
				t.Errorf("%s, want: %v, got: %v", cas.name, cas.wantNode, got)
			}
		}
	}

	// ABCD
	{
		err = ch.AddNode(nodeC)
		if err == nil {
			t.Errorf("want err, got: %v", err)
		}
		for _, cas := range []struct {
			name     string
			data     string
			wantNode *Node
		}{
			{name: "case1", data: "1", wantNode: nodeC},
			{name: "case2", data: "2", wantNode: nodeA},
			{name: "case3", data: "3", wantNode: nodeB},
			{name: "case4", data: "4", wantNode: nodeC},
			{name: "case5", data: "5", wantNode: nodeC},
			{name: "case6", data: "6", wantNode: nodeA},
			{name: "case7", data: "7", wantNode: nodeB},
			{name: "case8", data: "8", wantNode: nodeB},
			{name: "case9", data: "9", wantNode: nodeC},
			{name: "case10", data: "10", wantNode: nodeA},
			{name: "case11", data: "11", wantNode: nodeA},
			{name: "case12", data: "12", wantNode: nodeC},
		} {
			got, err := ch.Hash([]byte(cas.data))
			if err != nil {
				t.Errorf("%s, want nil, got: %v", cas.name, err)
			}
			if !reflect.DeepEqual(got, cas.wantNode) {
				t.Errorf("%s, want: %v, got: %v", cas.name, cas.wantNode, got)
			}
		}

		err = ch.AddNode(nodeD)
		if err != nil {
			t.Errorf("want nil, got: %v", err)
		}
		for i := 0; i < 5; i++ {
			for _, cas := range []struct {
				name     string
				data     string
				wantNode *Node
			}{
				{name: "case1", data: "1", wantNode: nodeC},
				{name: "case2", data: "2", wantNode: nodeA},
				{name: "case3", data: "3", wantNode: nodeD},
				{name: "case4", data: "4", wantNode: nodeC},
				{name: "case5", data: "5", wantNode: nodeC},
				{name: "case6", data: "6", wantNode: nodeA},
				{name: "case7", data: "7", wantNode: nodeD},
				{name: "case8", data: "8", wantNode: nodeB},
				{name: "case9", data: "9", wantNode: nodeC},
				{name: "case10", data: "10", wantNode: nodeA},
				{name: "case11", data: "11", wantNode: nodeA},
				{name: "case12", data: "12", wantNode: nodeC},
			} {
				got, err := ch.Hash([]byte(cas.data))
				if err != nil {
					t.Errorf("%s, want nil, got: %v", cas.name, err)
				}
				if !reflect.DeepEqual(got, cas.wantNode) {
					t.Errorf("%s, want: %v, got: %v", cas.name, cas.wantNode, got)
				}
			}
		}
	}
}

func TestCHashBlance(t *testing.T) {
	ch, err := NewCHash[*Node]([]*Node{nodeA, nodeB, nodeC},
		CHashOptionIndexer[*Node](func(data []byte) uint32 { return crc32.ChecksumIEEE(data) }),
		CHashOptionNodeIDer[*Node](func(node *Node) ([]byte, error) { return []byte(node.Name), nil }),
		// CHashOptionVirtualNodeFactor[*Node](500),
	)
	if err != nil {
		t.Errorf("want nil, got: %v", err)
		return
	}

	total := 10000
	frequencies := map[string]int{}
	for i := 0; i < total; i++ {
		bs := make([]byte, mrand.Int31()%1000)
		_, err = rand.Read(bs)
		if err != nil {
			t.Errorf("want nil, got: %v", err)
			return
		}

		node, err := ch.Hash(bs)
		if err != nil {
			t.Errorf("want nil, got: %v", err)
			return
		}

		frequencies[node.Name]++
	}

	for name, count := range frequencies {
		if count < total*30/100 || count > total*36/100 {
			t.Error(name, count, total)
		}
	}
}

func TestCHashfind(t *testing.T) {
	ch := &CHash[int]{
		virtualNodeList: []*virtualNode[int]{
			{
				beginIndex: 5,
				realNode:   realNode[int]{node: 5},
			},

			{
				beginIndex: 10,
				realNode:   realNode[int]{node: 10},
			},
			{
				beginIndex: 15,
				realNode:   realNode[int]{node: 15},
			},
			{
				beginIndex: 20,
				realNode:   realNode[int]{node: 20},
			},
		},
	}
	for _, one := range []struct {
		index    uint32
		wantNode int
	}{
		{index: 0, wantNode: 5},
		{index: 1, wantNode: 5},

		{index: 5, wantNode: 10},
		{index: 7, wantNode: 10},

		{index: 15, wantNode: 20},
		{index: 18, wantNode: 20},

		{index: 20, wantNode: 5},
		{index: 50, wantNode: 5},
	} {
		if got := ch.find(one.index); got != one.wantNode {
			t.Errorf("index: %d, want: %d, got: %d", one.index, one.wantNode, got)
		}
	}
}

func ExampleCHash() {

	ch, err := NewCHash[*Node]([]*Node{{Name: "A"}, {Name: "B"}, {Name: "C"}},
		CHashOptionIndexer[*Node](func(data []byte) uint32 { return crc32.ChecksumIEEE(data) }),
		CHashOptionNodeIDer[*Node](func(node *Node) ([]byte, error) { return []byte(node.Name), nil }),
		CHashOptionVirtualNodeFactor[*Node](5),
	)
	fmt.Println(err)

	node, err := ch.Hash([]byte("1"))
	fmt.Println(node.Name, err)
	node, err = ch.Hash([]byte("1"))
	fmt.Println(node.Name, err)

	node, err = ch.Hash([]byte("2"))
	fmt.Println(node.Name, err)
	node, err = ch.Hash([]byte("2"))
	fmt.Println(node.Name, err)

	node, err = ch.Hash([]byte("3"))
	fmt.Println(node.Name, err)
	node, err = ch.Hash([]byte("3"))
	fmt.Println(node.Name, err)

	node, err = ch.Hash([]byte("4"))
	fmt.Println(node.Name, err)
	node, err = ch.Hash([]byte("5"))
	fmt.Println(node.Name, err)
	node, err = ch.Hash([]byte("6"))

	node, err = ch.Hash([]byte("7"))
	fmt.Println(node.Name, err)
	node, err = ch.Hash([]byte("8"))
	fmt.Println(node.Name, err)
	node, err = ch.Hash([]byte("9"))

}
