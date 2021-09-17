package cslb

import (
	"math"
	"time"
)

const (
	TTLUnlimited time.Duration = math.MaxInt64 // Never expire
	TTLNone      time.Duration = 0             // Refresh after every Next()

	HealthyNodeMustAll float64 = 1.0
	HealthyNodeAny     float64 = 0.0
)

var (
	DefaultLoadBalancerOption = LoadBalancerOption{
		TTL:                 TTLUnlimited,
		MinHealthyNodeRatio: HealthyNodeAny,
	}
)

type LoadBalancerOption struct {
	// Cache TTL
	TTL time.Duration

	// Refresh when healthy node ratio is below MinHealthyNodeRatio
	MinHealthyNodeRatio float64
}
