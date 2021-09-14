package cslb

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RangerCD/cslb/service"
	"github.com/RangerCD/cslb/strategy"
)

func TestCSLB(t *testing.T) {
	lb := NewLoadBalancer(
		service.NewRRDNSService(
			[]string{
				"example.com",
			}, true, true,
		),
		strategy.NewRoundRobinStrategy(),
	)

	for i := 0; i < 10; i++ {
		next, err := lb.Next()
		log.Println(next)
		assert.Nil(t, err)
	}
}
