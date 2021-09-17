package node

import (
	"net"
)

type Node struct {
	addr net.Addr
}

func NewNode(addr net.Addr) *Node {
	return &Node{
		addr: addr,
	}
}

func (n Node) Key() string {
	return n.addr.String()
}

func (n Node) Addr() net.Addr {
	return n.addr
}
