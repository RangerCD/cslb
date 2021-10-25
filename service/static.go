package service

import (
	"net"
	"sync/atomic"
	"unsafe"
)

// staticService represents a simple static list of net.Addr
type staticService struct {
	staticAddrs []net.Addr
	addrs       unsafe.Pointer // Pointer to []net.Addr
}

func NewStaticService(addrs []net.Addr) *staticService {
	return &staticService{
		staticAddrs: addrs,
		addrs:       nil,
	}
}

func (s *staticService) Nodes() []net.Addr {
	addrs := (*[]net.Addr)(atomic.LoadPointer(&s.addrs))
	result := make([]net.Addr, 0, len(*addrs))
	result = append(result, *addrs...)
	return result
}

func (s *staticService) Refresh() {
	atomic.StorePointer(&s.addrs, (unsafe.Pointer)(&s.staticAddrs))
}

func (s *staticService) NodeFailedCallbackFunc() func(addr net.Addr) {
	return nil
}
