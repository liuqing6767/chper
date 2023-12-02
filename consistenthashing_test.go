package chper

import (
	"crypto/rand"
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
		CHashOptionNodeNaming[*Node](func(node *Node) (string, error) { return node.Name, nil }),
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
				{name: "case1", data: "1", wantNode: nodeA},
				{name: "case2", data: "2", wantNode: nodeC},
				{name: "case3", data: "3", wantNode: nodeA},
				{name: "case4", data: "4", wantNode: nodeA},
				{name: "case5", data: "5", wantNode: nodeA},
				{name: "case6", data: "6", wantNode: nodeC},
				{name: "case7", data: "7", wantNode: nodeA},
				{name: "case8", data: "8", wantNode: nodeA},
				{name: "case9", data: "9", wantNode: nodeA},
				{name: "case10", data: "10", wantNode: nodeB},
				{name: "case11", data: "11", wantNode: nodeB},
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
			{name: "case1", data: "1", wantNode: nodeA},
			{name: "case2", data: "2", wantNode: nodeC},
			{name: "case3", data: "3", wantNode: nodeA},
			{name: "case4", data: "4", wantNode: nodeA},
			{name: "case5", data: "5", wantNode: nodeA},
			{name: "case6", data: "6", wantNode: nodeC},
			{name: "case7", data: "7", wantNode: nodeA},
			{name: "case8", data: "8", wantNode: nodeA},
			{name: "case9", data: "9", wantNode: nodeA},
			{name: "case10", data: "10", wantNode: nodeB},
			{name: "case11", data: "11", wantNode: nodeB},
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
				{name: "case1", data: "1", wantNode: nodeA},
				{name: "case2", data: "2", wantNode: nodeB},
				{name: "case3", data: "3", wantNode: nodeA},
				{name: "case4", data: "4", wantNode: nodeA},
				{name: "case5", data: "5", wantNode: nodeA},
				{name: "case6", data: "6", wantNode: nodeB},
				{name: "case7", data: "7", wantNode: nodeA},
				{name: "case8", data: "8", wantNode: nodeA},
				{name: "case9", data: "9", wantNode: nodeA},
				{name: "case10", data: "10", wantNode: nodeB},
				{name: "case11", data: "11", wantNode: nodeB},
				{name: "case12", data: "12", wantNode: nodeA},
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
			{name: "case1", data: "1", wantNode: nodeA},
			{name: "case2", data: "2", wantNode: nodeC},
			{name: "case3", data: "3", wantNode: nodeA},
			{name: "case4", data: "4", wantNode: nodeA},
			{name: "case5", data: "5", wantNode: nodeA},
			{name: "case6", data: "6", wantNode: nodeC},
			{name: "case7", data: "7", wantNode: nodeA},
			{name: "case8", data: "8", wantNode: nodeA},
			{name: "case9", data: "9", wantNode: nodeA},
			{name: "case10", data: "10", wantNode: nodeB},
			{name: "case11", data: "11", wantNode: nodeB},
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
			{name: "case1", data: "1", wantNode: nodeA},
			{name: "case2", data: "2", wantNode: nodeC},
			{name: "case3", data: "3", wantNode: nodeA},
			{name: "case4", data: "4", wantNode: nodeA},
			{name: "case5", data: "5", wantNode: nodeA},
			{name: "case6", data: "6", wantNode: nodeC},
			{name: "case7", data: "7", wantNode: nodeA},
			{name: "case8", data: "8", wantNode: nodeA},
			{name: "case9", data: "9", wantNode: nodeA},
			{name: "case10", data: "10", wantNode: nodeB},
			{name: "case11", data: "11", wantNode: nodeB},
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
				{name: "case1", data: "1", wantNode: nodeA},
				{name: "case2", data: "2", wantNode: nodeC},
				{name: "case3", data: "3", wantNode: nodeA},
				{name: "case4", data: "4", wantNode: nodeA},
				{name: "case5", data: "5", wantNode: nodeA},
				{name: "case6", data: "6", wantNode: nodeC},
				{name: "case7", data: "7", wantNode: nodeA},
				{name: "case8", data: "8", wantNode: nodeA},
				{name: "case9", data: "9", wantNode: nodeA},
				{name: "case10", data: "10", wantNode: nodeB},
				{name: "case11", data: "11", wantNode: nodeD},
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

func TestCHashBalance(t *testing.T) {
	ch, err := NewCHash[*Node]([]*Node{nodeA, nodeB, nodeC},
		CHashOptionIndexer[*Node](func(data []byte) uint32 { return crc32.ChecksumIEEE(data) }),
		CHashOptionNodeNaming[*Node](func(node *Node) (string, error) { return node.Name, nil }),
		// CHashOptionVirtualNodeFactor[*Node](500),
	)
	if err != nil {
		t.Errorf("want nil, got: %v", err)
		return
	}

	total := 10000
	{
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
			if count < total*29/100 || count > total*40/100 {
				t.Error(name, count, total)
			}
		}
	}

	{
		err = ch.AddNode(nodeD)
		if err != nil {
			t.Errorf("want nil, got: %v", err)
		}

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
			if count < total*22/100 || count > total*28/100 {
				t.Error(name, count, total)
			}
		}
	}
}

func TestCHashWeightBalance(t *testing.T) {
	ch, err := NewCHash[*Node]([]*Node{nodeA, nodeB, nodeC},
		CHashOptionIndexer[*Node](func(data []byte) uint32 { return crc32.ChecksumIEEE(data) }),
		CHashOptionNodeNaming[*Node](func(node *Node) (string, error) { return node.Name, nil }),
		CHashOptionWeightSpecify[*Node](func(node *Node) int {
			if node.Name == "A" || node.Name == "D" {
				return 98
			}

			return 1
		}),
		// CHashOptionVirtualNodeFactor[*Node](500),
	)
	if err != nil {
		t.Errorf("want nil, got: %v", err)
		return
	}

	total := 10000
	{
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
		{
			name := "A"
			count := frequencies[name]
			if count < 9700 {
				t.Error(name, count, total)
			}
		}
		{
			name := "B"
			count := frequencies[name]
			if count > 200 {
				t.Error(name, count, total)
			}
		}
		{
			name := "C"
			count := frequencies[name]
			if count > 200 {
				t.Error(name, count, total)
			}
		}
	}

	ch.AddNode(nodeD)
	{
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
		{
			name := "A"
			count := frequencies[name]
			if count < 4700 || count > 5300 {
				t.Error(name, count, total)
			}
		}
		{
			name := "B"
			count := frequencies[name]
			if count > 100 {
				t.Error(name, count, total)
			}
		}
		{
			name := "C"
			count := frequencies[name]
			if count > 100 {
				t.Error(name, count, total)
			}
		}
		{
			name := "D"
			count := frequencies[name]
			if count < 4700 || count > 5300 {
				t.Error(name, count, total)
			}
		}
	}
}
