package strategy

import (
	"log"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRoundRobin(t *testing.T) {
	s := roundRobinStrategy{}
	s.SetNodes([]net.Addr{
		&net.IPAddr{IP: net.IPv4(1, 2, 3, 4)},
		&net.IPAddr{IP: net.IPv4(2, 3, 4, 5)},
		&net.IPAddr{IP: net.IPv4(3, 4, 5, 6)},
		&net.IPAddr{IP: net.IPv4(4, 5, 6, 7)},
	})

	for i := 0; i < 10; i++ {
		next, err := s.Next()
		log.Println(next)
		assert.Nil(t, err)
	}
}

func TestConcurrentRoundRobin(t *testing.T) {
	rand.Seed(0)
	s := roundRobinStrategy{}
	s.SetNodes([]net.Addr{
		&net.IPAddr{IP: net.IPv4(1, 2, 3, 4)},
		&net.IPAddr{IP: net.IPv4(2, 3, 4, 5)},
		&net.IPAddr{IP: net.IPv4(3, 4, 5, 6)},
		&net.IPAddr{IP: net.IPv4(4, 5, 6, 7)},
	})

	// Concurrent read
	for i := 0; i < 4; i++ {
		go func() {
			for {
				_, err := s.Next()
				assert.Nil(t, err)
			}
		}()
	}

	// Concurrent write
	for i := 0; i < 4; i++ {
		go func() {
			for {
				s.SetNodes([]net.Addr{
					&net.IPAddr{IP: net.IPv4(1, 2, 3, byte(rand.Intn(256)))},
					&net.IPAddr{IP: net.IPv4(2, 3, 4, byte(rand.Intn(256)))},
					&net.IPAddr{IP: net.IPv4(3, 4, 5, byte(rand.Intn(256)))},
					&net.IPAddr{IP: net.IPv4(4, 5, 6, byte(rand.Intn(256)))},
				})
			}
		}()
	}

	time.Sleep(time.Second * 1)
}

func BenchmarkRoundRobin(b *testing.B) {
	s := roundRobinStrategy{}
	s.SetNodes([]net.Addr{
		&net.IPAddr{IP: net.IPv4(1, 2, 3, 4)},
		&net.IPAddr{IP: net.IPv4(2, 3, 4, 5)},
		&net.IPAddr{IP: net.IPv4(3, 4, 5, 6)},
		&net.IPAddr{IP: net.IPv4(4, 5, 6, 7)},
	})
	for i := 0; i < b.N; i++ {
		s.Next()
	}
}
