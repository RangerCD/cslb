package cslb

import (
	"errors"
	"math"
	"math/rand"
	"sync/atomic"
	"unsafe"
)

type roundRobinStrategy struct {
	index         uint64
	internalNodes unsafe.Pointer // pointer to roundRobinInternalNodes
}

type roundRobinInternalNodes struct {
	nodes []Node
	order []int
}

func NewRoundRobinStrategy() *roundRobinStrategy {
	return &roundRobinStrategy{
		index:         math.MaxUint64,
		internalNodes: nil,
	}
}

func (s *roundRobinStrategy) SetNodes(nodes []Node) {
	internalNodes := &roundRobinInternalNodes{
		nodes: nodes,
		order: make([]int, len(nodes)),
	}
	for i := 0; i < len(internalNodes.order); i++ {
		internalNodes.order[i] = i
	}
	atomic.StorePointer(&s.internalNodes, unsafe.Pointer(internalNodes))
}

func (s *roundRobinStrategy) Next() (Node, error) {
	internalNodes := (*roundRobinInternalNodes)(atomic.LoadPointer(&s.internalNodes))
	nodes := internalNodes.nodes
	order := internalNodes.order
	if len(nodes) > 0 {
		index := atomic.AddUint64(&s.index, 1) % uint64(len(order))
		return nodes[order[index]], nil
	} else {
		return nil, errors.New("empty node list")
	}
}

func (s *roundRobinStrategy) NextFor(interface{}) (Node, error) {
	return s.Next()
}
