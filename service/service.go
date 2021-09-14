package service

import "net"

type Service interface {
	Nodes() []net.Addr
	NodeFailed(addr net.Addr)
	Refresh()
}
