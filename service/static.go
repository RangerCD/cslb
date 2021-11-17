package service

import (
	"sync/atomic"
	"unsafe"

	"github.com/RangerCD/cslb/node"
)

// staticService represents a simple static list of net.Addr
// Node type: node.Node
type staticService struct {
	staticNodes []node.Node
	nodes       unsafe.Pointer // Pointer to []node.Node
}

func NewStaticService(nodes []node.Node) *staticService {
	return &staticService{
		staticNodes: nodes,
		nodes:       nil,
	}
}

func (s *staticService) Nodes() []node.Node {
	nodes := (*[]node.Node)(atomic.LoadPointer(&s.nodes))
	result := make([]node.Node, 0, len(*nodes))
	result = append(result, *nodes...)
	return result
}

func (s *staticService) Refresh() {
	atomic.StorePointer(&s.nodes, (unsafe.Pointer)(&s.staticNodes))
}

func (s *staticService) NodeFailedCallbackFunc() func(node node.Node) {
	return nil
}
