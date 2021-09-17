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
	ipv4      bool
	ipv6      bool
	hostnames []string
	addrs     unsafe.Pointer // Pointer to []net.Addr
}

func NewRRDNSService(hostnames []string, ipv4 bool, ipv6 bool) *rrDNSService {
	return &rrDNSService{
		ipv4:      ipv4,
		ipv6:      ipv6,
		hostnames: hostnames,
		addrs:     nil,
	}
}

func (s *rrDNSService) Nodes() []net.Addr {
	addrs := (*[]net.Addr)(atomic.LoadPointer(&s.addrs))
	result := make([]net.Addr, 0, len(*addrs))
	result = append(result, *addrs...)
	return result
}

func (s *rrDNSService) NodeFailedCallbackFunc() func(addr net.Addr) {
	return nil
}

func (s *rrDNSService) Refresh() {
	ips := make([]net.Addr, 0, len(s.hostnames))
	for _, h := range s.hostnames {
		if results, err := net.LookupIP(h); err == nil {
			for _, ip := range results {
				switch {
				case ip.To4() != nil && s.ipv4:
					fallthrough
				case ip.To16() != nil && s.ipv6:
					ips = append(ips, &net.IPAddr{IP: ip})
				}
			}
		}
	}
	atomic.StorePointer(&s.addrs, (unsafe.Pointer)(&ips))
}
