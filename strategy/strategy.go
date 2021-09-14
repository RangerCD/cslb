package strategy

import "net"

type Strategy interface {
	SetNodes(addrs []net.Addr)
	Next() (net.Addr, error)
}
