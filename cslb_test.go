package cslb

import (
	"log"
	"testing"

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
