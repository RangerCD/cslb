package cslb

import (
	"math"
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

	TTLUnlimited = math.MaxInt64 // Never expire
	TTLNone      = 0             // Refresh after every Next()
)

type LoadBalancer struct {
	service    service.Service
	strategy   strategy.Strategy
	ttlSecond  int64
	lastUpdate int64
	sf         *singleflight.Group
}

func NewLoadBalancer(service service.Service, strategy strategy.Strategy, ttlSecond int64) *LoadBalancer {
	lb := &LoadBalancer{
		service:    service,
		strategy:   strategy,
		ttlSecond:  ttlSecond,
		lastUpdate: 0,
		sf:         new(singleflight.Group),
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
	lived := time.Now().Unix() - atomic.LoadInt64(&lb.lastUpdate)
	if lb.ttlSecond != TTLUnlimited && (lived > lb.ttlSecond || lived < 0) {
		// Background refresh
		lb.refresh()
	}
	return next, err
}

func (lb *LoadBalancer) NodeFailed(node net.Addr) {
	lb.sf.Do(NodeFailedKey+node.String(), func() (interface{}, error) {
		lb.service.NodeFailed(node)
		nodes := lb.service.Nodes()
		if len(nodes) <= 0 {
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
		lb.strategy.SetNodes(lb.service.Nodes())
		return nil, nil
	})
}
