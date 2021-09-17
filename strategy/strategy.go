package strategy

import "net"

// Strategy controls how the nodes are chosen
// This type should be thread safe
type Strategy interface {
	// SetNodes update saved nodes
	SetNodes(addrs []net.Addr)
	// Next returns a node address
	Next() (net.Addr, error)
}
