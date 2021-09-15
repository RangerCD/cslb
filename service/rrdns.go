package service

import (
	"net"
	"sync"
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
	addrs     sync.Map // string => net.Addr
}

func NewRRDNSService(hostnames []string, ipv4 bool, ipv6 bool) *rrDNSService {
	return &rrDNSService{
		ipv4:      ipv4,
		ipv6:      ipv6,
		hostnames: hostnames,
		addrs:     sync.Map{},
	}
}

func (s *rrDNSService) Nodes() []net.Addr {
	result := make([]net.Addr, 0)
	s.addrs.Range(func(key, value interface{}) bool {
		result = append(result, value.(net.Addr))
		return true
	})
	return result
}

func (s *rrDNSService) NodeFailed(node net.Addr) {
	s.addrs.Delete(node.String())
}

func (s *rrDNSService) Refresh() {
	ips := make(map[string]net.Addr, len(s.hostnames))
	for _, h := range s.hostnames {
		if results, err := net.LookupIP(h); err == nil {
			for _, ip := range results {
				switch {
				case ip.To4() != nil && s.ipv4:
					fallthrough
				case ip.To16() != nil && s.ipv6:
					ip := &net.IPAddr{IP: ip}
					ips[ip.String()] = ip
				}
			}
		}
	}

	for k, v := range ips {
		s.addrs.Store(k, v)
	}

	s.addrs.Range(func(key, value interface{}) bool {
		if _, ok := ips[key.(string)]; !ok {
			s.addrs.Delete(key)
		}
		return true
	})
}
