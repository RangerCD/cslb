package strategy

import (
	"errors"
	"math"
	"net"
	"sync/atomic"
	"unsafe"
)

type roundRobinStrategy struct {
	index uint64
	nodes unsafe.Pointer // pointer to []net.Addr
}

func NewRoundRobinStrategy() *roundRobinStrategy {
	return &roundRobinStrategy{
		index: math.MaxUint64,
		nodes: nil,
	}
}

func (s *roundRobinStrategy) SetNodes(addrs []net.Addr) {
	atomic.StorePointer(&s.nodes, unsafe.Pointer(&addrs))
}

func (s *roundRobinStrategy) Next() (net.Addr, error) {
	nodes := (*[]net.Addr)(atomic.LoadPointer(&s.nodes))
	if len(*nodes) > 0 {
		index := atomic.AddUint64(&s.index, 1) % uint64(len(*nodes))
		return (*nodes)[index], nil
	} else {
		return nil, errors.New("empty node list")
	}
}
