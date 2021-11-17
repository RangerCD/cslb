package service

import (
	"log"
	"net"
	"testing"

	"github.com/RangerCD/cslb/node"
)

func TestStaticService(t *testing.T) {
	s := NewStaticService(
		[]node.Node{
			&net.TCPAddr{
				IP:   net.IPv4(1, 2, 3, 4),
				Port: 1234,
			},
			&net.TCPAddr{
				IP:   net.IPv4(2, 3, 4, 5),
				Port: 2345,
			},
		})
	s.Refresh()
	log.Println(s.Nodes())
}
