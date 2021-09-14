package cslb

import (
	"net"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/RangerCD/cslb/service"
	"github.com/RangerCD/cslb/strategy"
)

const (
	NodeFailedKey = "node-failed."
	RefreshKey    = "refresh"
)

type LoadBalancer struct {
	Service    service.Service
	Strategy   strategy.Strategy
	TTLSecond  int64
	lastUpdate int64
	sf         *singleflight.Group
}

func NewLoadBalancer(service service.Service, strategy strategy.Strategy) *LoadBalancer {
	lb := &LoadBalancer{
		Service:  service,
		Strategy: strategy,
		sf:       new(singleflight.Group),
	}
	<-lb.refresh()
	return lb
}

func (lb *LoadBalancer) Next() (net.Addr, error) {
	next, err := lb.Strategy.Next()
	if err != nil {
		lb.NodeFailed(nil)
		next, err = lb.Strategy.Next()
	}
	lived := time.Now().Unix() - atomic.LoadInt64(&lb.lastUpdate)
	if lived > lb.TTLSecond || lived < 0 {
		lb.refresh()
	}
	return next, err
}

func (lb *LoadBalancer) NodeFailed(node net.Addr) {
	lb.sf.Do(NodeFailedKey+node.String(), func() (interface{}, error) {
		lb.Service.NodeFailed(node)
		lb.Strategy.SetNodes(lb.Service.Nodes())
		return nil, nil
	})
}

func (lb *LoadBalancer) refresh() <-chan singleflight.Result {
	return lb.sf.DoChan(RefreshKey, func() (interface{}, error) {
		lb.Service.Refresh()
		atomic.StoreInt64(&lb.lastUpdate, time.Now().Unix())
		lb.Strategy.SetNodes(lb.Service.Nodes())
		return nil, nil
	})
}
