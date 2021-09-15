package cslb

import (
	"context"
	"log"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/RangerCD/cslb/service"
	"github.com/RangerCD/cslb/strategy"
)

func TestCSLB(t *testing.T) {
	srv := service.NewRRDNSService(
		[]string{
			"example.com",
		}, true, true,
	)
	stg := strategy.NewRoundRobinStrategy()
	lb := NewLoadBalancer(
		srv,
		stg,
		TTLUnlimited,
	)

	nodes := srv.Nodes()

	for i := 0; i < 10; i++ {
		next, err := lb.Next()
		log.Println(next)
		assert.Contains(t, nodes, next)
		assert.Nil(t, err)
	}
}

func Test100RCSLBRandomFail(t *testing.T) {
	var counter uint64 = 0
	var failedCounter uint64 = 0
	srv := service.NewRRDNSService(
		[]string{
			"example.com",
		}, true, true,
	)
	stg := strategy.NewRoundRobinStrategy()
	lb := NewLoadBalancer(
		srv,
		stg,
		TTLUnlimited,
	)

	ctx, cancel := context.WithCancel(context.Background())
	// 100 concurrent read & 10% random fail
	for i := 0; i < 100; i++ {
		go func() {
			done := ctx.Done()
			for {
				select {
				case <-done:
					return
				default:
				}
				n, err := lb.Next()
				assert.Nil(t, err)
				atomic.AddUint64(&counter, 1)
				if rand.Intn(10) < 1 {
					lb.NodeFailed(n)
					atomic.AddUint64(&failedCounter, 1)
				}
			}
		}()
	}

	time.Sleep(time.Second * 1)
	cancel()

	log.Println("Next() called", counter, "times")
	log.Println("NodeFailed() called", failedCounter, "times")
}
