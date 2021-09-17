package node

import (
	"math"
	"net"
	"sync"
	"sync/atomic"
)

const (
	NodeCountUnlimited = math.MaxInt64
)

type Group struct {
	m             sync.Map // string => *Node
	originalCount int64
	currentCount  int64

	maxNodeCount int
}

func NewGroup(maxNodeCount int) *Group {
	return &Group{
		m:             sync.Map{},
		originalCount: 0,
		currentCount:  0,
		maxNodeCount:  maxNodeCount,
	}
}

func (g *Group) Set(addrs []net.Addr) {
	nodesCheck := make(map[string]struct{}, len(addrs))
	for _, addr := range addrs {
		node := NewNode(addr)
		key := node.Key()
		nodesCheck[key] = struct{}{}
		g.m.Store(key, node)
	}

	g.m.Range(func(key, value interface{}) bool {
		if _, ok := nodesCheck[key.(string)]; !ok {
			g.m.Delete(key)
		}
		return true
	})

	atomic.StoreInt64(&g.originalCount, int64(len(addrs)))
	atomic.StoreInt64(&g.currentCount, int64(len(addrs)))
}

func (g *Group) Get() []net.Addr {
	result := make([]net.Addr, 0)
	g.m.Range(func(key, value interface{}) bool {
		result = append(result, value.(*Node).Addr())
		if len(result) >= g.maxNodeCount {
			return false
		}
		return true
	})
	return result
}

func (g *Group) GetNode(key string) *Node {
	if val, loaded := g.m.Load(key); loaded {
		return (val).(*Node)
	}
	return nil
}

func (g *Group) GetOriginalCount() int64 {
	return atomic.LoadInt64(&g.originalCount)
}

func (g *Group) GetCurrentCount() int64 {
	return atomic.LoadInt64(&g.currentCount)
}

func (g *Group) Exile(addr net.Addr) bool {
	_, loaded := g.m.LoadAndDelete(addr.String())
	if loaded {
		atomic.AddInt64(&g.currentCount, -1)
	}
	return loaded
}
