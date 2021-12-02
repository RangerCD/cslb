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
		MaxNodeCount:        NodeCountUnlimited,
		TTL:                 TTLUnlimited,
		MinHealthyNodeRatio: HealthyNodeAny,
	}
)

type LoadBalancerOption struct {
	// LoadBalancer will keep MaxNodeCount nodes in Next() result set
	// Please notice that refresh or exile will change the result set
	// Number of connections might be greater than this value, if any pre-connected node has been excluded in newer
	// result set, but no new connection will be established to these nodes
	MaxNodeCount int

	// Cache TTL
	TTL time.Duration

	// Refresh when healthy node ratio is below MinHealthyNodeRatio
	MinHealthyNodeRatio float64
}
