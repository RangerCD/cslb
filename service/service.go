package service

import "net"

type Service interface {
	// Nodes returns a new slice of available node
	Nodes() []net.Addr
	// NodeFailed will be called when certain node failed too many times
	NodeFailed(addr net.Addr)
	// Refresh updates nodes
	Refresh()
}
