package cslb

import (
	"math"
	"net"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/RangerCD/cslb/node"
	"github.com/RangerCD/cslb/service"
	"github.com/RangerCD/cslb/strategy"
)

const (
	NodeFailedKey = "node-failed."
	RefreshKey    = "refresh"
)

type LoadBalancer struct {
	service    service.Service
	strategy   strategy.Strategy
	option     LoadBalancerOption
	lastUpdate int64
	sf         *singleflight.Group
	nodes      *node.Group
}

func NewLoadBalancer(service service.Service, strategy strategy.Strategy, option ...LoadBalancerOption) *LoadBalancer {
	opt := DefaultLoadBalancerOption
	if len(option) > 0 {
		opt = option[0]
	}
	lb := &LoadBalancer{
		service:    service,
		strategy:   strategy,
		option:     opt,
		lastUpdate: 0,
		sf:         new(singleflight.Group),
		nodes:      node.NewGroup(),
	}
	<-lb.refresh()
	return lb
}

func (lb *LoadBalancer) Next() (net.Addr, error) {
	next, err := lb.strategy.Next()
	if err != nil {
		// Refresh and retry
		<-lb.refresh()
		next, err = lb.strategy.Next()
	}
	// TODO: replace with time.Timer, it's too heavy for Next()
	lived := time.Duration(time.Now().Unix()-atomic.LoadInt64(&lb.lastUpdate)) * time.Second
	if lb.option.TTL != TTLUnlimited && (lived > lb.option.TTL || lived < 0) {
		// Background refresh
		lb.refresh()
	}
	return next, err
}

func (lb *LoadBalancer) NodeFailed(node net.Addr) {
	lb.sf.Do(NodeFailedKey+node.String(), func() (interface{}, error) {
		// TODO: allow fail several times before exile
		lb.nodes.Exile(node)
		if fn := lb.service.NodeFailedCallbackFunc(); fn != nil {
			go fn(node)
		}
		nodes := lb.nodes.Get()
		if len(nodes) <= 0 ||
			math.Round(float64(lb.nodes.GetOriginalCount())*lb.option.MinHealthyNodeRatio) > float64(lb.nodes.GetCurrentCount()) {
			<-lb.refresh()
		} else {
			lb.strategy.SetNodes(nodes)
		}
		return nil, nil
	})
}

func (lb *LoadBalancer) refresh() <-chan singleflight.Result {
	return lb.sf.DoChan(RefreshKey, func() (interface{}, error) {
		lb.service.Refresh()
		atomic.StoreInt64(&lb.lastUpdate, time.Now().Unix())
		lb.nodes.Set(lb.service.Nodes())
		lb.strategy.SetNodes(lb.nodes.Get())
		return nil, nil
	})
}
