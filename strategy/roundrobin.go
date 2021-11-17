package strategy

import (
	"errors"
	"math"
	"sync/atomic"
	"unsafe"

	"github.com/RangerCD/cslb/node"
)

type roundRobinStrategy struct {
	index uint64
	nodes unsafe.Pointer // pointer to []node.Node
}

func NewRoundRobinStrategy() *roundRobinStrategy {
	return &roundRobinStrategy{
		index: math.MaxUint64,
		nodes: nil,
	}
}

func (s *roundRobinStrategy) SetNodes(nodes []node.Node) {
	atomic.StorePointer(&s.nodes, unsafe.Pointer(&nodes))
}

func (s *roundRobinStrategy) Next() (node.Node, error) {
	nodes := (*[]node.Node)(atomic.LoadPointer(&s.nodes))
	if len(*nodes) > 0 {
		index := atomic.AddUint64(&s.index, 1) % uint64(len(*nodes))
		return (*nodes)[index], nil
	} else {
		return nil, errors.New("empty node list")
	}
}
