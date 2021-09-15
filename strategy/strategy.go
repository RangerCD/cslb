package strategy

import "net"

type Strategy interface {
	// SetNodes update saved nodes
	SetNodes(addrs []net.Addr)
	// Next returns a node address
	Next() (net.Addr, error)
}
