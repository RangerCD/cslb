package strategy

import (
	"errors"
	"net"
	"sync/atomic"
	"unsafe"
)

type roundRobinStrategy struct {
	nodesChan unsafe.Pointer // pointer to chan net.Addr
}

func NewRoundRobinStrategy() *roundRobinStrategy {
	return &roundRobinStrategy{}
}

func (s *roundRobinStrategy) SetNodes(addrs []net.Addr) {
	nodes := make(chan net.Addr, len(addrs))
	for _, addr := range addrs {
		nodes <- addr
	}
	atomic.StorePointer(&s.nodesChan, unsafe.Pointer(&nodes))
}

func (s *roundRobinStrategy) Next() (net.Addr, error) {
	nodesChan := (*chan net.Addr)(atomic.LoadPointer(&s.nodesChan))
	select {
	case addr := <-*nodesChan:
		*nodesChan <- addr
		return addr, nil
	default:
		return nil, errors.New("empty node list")
	}
}
