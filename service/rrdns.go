package service

import (
	"net"
	"sync/atomic"
	"unsafe"
)

// rrDNSService is for Round-robin DNS load balancing solution.
// Usually multiple A or AAAA records are associated with single hostname.
//
// For example:
//   Hostname www.a.com
//     |- A 1.2.3.4
//     |- A 2.3.4.5
//     |- A 3.4.5.6
//     ...
// Everytime a client wants to send a request, one of these A records will be chosen to establish connection.
type rrDNSService struct {
	IPv4      bool
	IPv6      bool
	Hostnames []string
	addrs     unsafe.Pointer // pointer to []net.Addr
}

func NewRRDNSService(hostnames []string, ipv4 bool, ipv6 bool) *rrDNSService {
	return &rrDNSService{
		IPv4:      ipv4,
		IPv6:      ipv6,
		Hostnames: hostnames,
		addrs:     nil,
	}
}

func (s *rrDNSService) Nodes() []net.Addr {
	return *(*[]net.Addr)(atomic.LoadPointer(&s.addrs))
}

func (s *rrDNSService) NodeFailed(node net.Addr) {
	// TODO: add rate limited refresh operation
}

func (s *rrDNSService) Refresh() {
	ips := make([]net.Addr, 0, len(s.Hostnames))
	for _, h := range s.Hostnames {
		if results, err := net.LookupIP(h); err == nil {
			for _, ip := range results {
				switch {
				case ip.To4() != nil && s.IPv4:
					fallthrough
				case ip.To16() != nil && s.IPv6:
					ips = append(ips, &net.IPAddr{IP: ip})
				}
			}
		}
	}

	atomic.StorePointer(&s.addrs, unsafe.Pointer(&ips))
}
