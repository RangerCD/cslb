package cslb

import (
	"log"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWeightedRoundRobin(t *testing.T) {
	nodes := []Node{
		&net.IPAddr{IP: net.IPv4(1, 2, 3, 4)},
		&net.IPAddr{IP: net.IPv4(2, 3, 4, 5)},
		&net.IPAddr{IP: net.IPv4(3, 4, 5, 6)},
		&net.IPAddr{IP: net.IPv4(4, 5, 6, 7)},
	}
	weightMap := make(map[string]int, len(nodes))
	for i, node := range nodes {
		weightMap[node.String()] = i + 1
	}
	s := NewWeightedRoundRobinStrategy(
		func(node Node) int {
			return weightMap[node.String()]
		},
	)
	s.SetNodes(nodes)

	sum := 0
	for _, weight := range weightMap {
		sum += weight
	}

	hitMap := make(map[string]int, len(nodes))
	for _, node := range nodes {
		hitMap[node.String()] = 0
	}
	for i := 0; i < sum; i++ {
		next, err := s.Next()
		log.Println(next)
		assert.Nil(t, err)
		hitMap[next.String()]++
	}
	for _, node := range nodes {
		assert.Equal(t, weightMap[node.String()], hitMap[node.String()])
	}
}

func Test4R4WWeightedRoundRobin(t *testing.T) {
	nodes := []Node{
		&net.IPAddr{IP: net.IPv4(1, 2, 3, 4)},
		&net.IPAddr{IP: net.IPv4(2, 3, 4, 5)},
		&net.IPAddr{IP: net.IPv4(3, 4, 5, 6)},
		&net.IPAddr{IP: net.IPv4(4, 5, 6, 7)},
	}
	weightMap := make(map[string]int, len(nodes))
	for i, node := range nodes {
		weightMap[node.String()] = i + 1
	}
	s := NewWeightedRoundRobinStrategy(
		func(node Node) int {
			return weightMap[node.String()]
		},
	)
	s.SetNodes(nodes)

	// 4 concurrent read
	for i := 0; i < 4; i++ {
		go func() {
			for {
				_, err := s.Next()
				assert.Nil(t, err)
			}
		}()
	}

	// 4 concurrent write
	for i := 0; i < 4; i++ {
		go func() {
			for {
				s.SetNodes([]Node{
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

func BenchmarkWeightedRoundRobin(b *testing.B) {
	nodes := []Node{
		&net.IPAddr{IP: net.IPv4(1, 2, 3, 4)},
		&net.IPAddr{IP: net.IPv4(2, 3, 4, 5)},
		&net.IPAddr{IP: net.IPv4(3, 4, 5, 6)},
		&net.IPAddr{IP: net.IPv4(4, 5, 6, 7)},
	}
	weightMap := make(map[string]int, len(nodes))
	for i, node := range nodes {
		weightMap[node.String()] = i + 1
	}
	s := NewWeightedRoundRobinStrategy(
		func(node Node) int {
			return weightMap[node.String()]
		},
	)
	s.SetNodes(nodes)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Next()
	}
	b.StopTimer()
}
