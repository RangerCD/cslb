package cslb

import (
	"math/big"
	"sync/atomic"
	"unsafe"
)

type weightedRoundRobinStrategy struct {
	*roundRobinStrategy
	weightFunc func(Node) int
}

// NewWeightedRoundRobinStrategy creates a new strategy with weighted round-robin method.
// This strategy will store a slice of node order in memory as cache, the length of the slice will be the sum of the weight
// of all nodes divided by the greatest common divisor of all weights.
// For example, a nodes configuration like this:
// | Node | Weight |
// | ---- | ------ |
// | A    | 10     |
// | B    | 10     |
// | C    | 15     |
// | D    | 20     |
// The greatest common divisor of [10, 10, 15, 20] is 5, so equivalent weights are [2, 2, 3, 4], which means the length
// of cache slice will be 2 + 2 + 3 + 4 = 11. Actually, the cache slice will be [D, C, A, B, D, C, D, A, B, C, D].
// Be careful when nodes have large weights and co-prime with each other, cache size might be very large.
func NewWeightedRoundRobinStrategy(weightFunc func(Node) int) *weightedRoundRobinStrategy {
	return &weightedRoundRobinStrategy{
		roundRobinStrategy: NewRoundRobinStrategy(),
		weightFunc:         weightFunc,
	}
}

func (s *weightedRoundRobinStrategy) SetNodes(nodes []Node) {
	nodes = s.generateWeightedNodes(nodes)
	atomic.StorePointer(&s.nodes, unsafe.Pointer(&nodes))
}

func (s *weightedRoundRobinStrategy) generateWeightedNodes(nodes []Node) []Node {
	if len(nodes) <= 1 {
		return nodes
	}

	weights := make([]int, len(nodes))
	for i, node := range nodes {
		weights[i] = s.weightFunc(node)
		if weights[i] < 0 {
			weights[i] = 0
		}
	}
	gcd := big.NewInt(int64(weights[0]))
	for i := 1; i < len(weights); i++ {
		gcd = gcd.GCD(nil, nil, gcd, big.NewInt(int64(weights[i])))
	}
	sum := 0
	gcdi := int(gcd.Int64())
	for i := 0; i < len(weights); i++ {
		weights[i] /= gcdi
		sum += weights[i]
	}

	result := make([]Node, 0, sum)
	curr := make([]int, len(weights))
	for i := 0; i < sum; i++ {
		for j := 0; j < len(weights); j++ {
			curr[j] += weights[j]
		}
		max := 0
		maxIndex := 0
		for j := 0; j < len(weights); j++ {
			if curr[j] > max {
				max = curr[j]
				maxIndex = j
			}
		}
		result = append(result, nodes[maxIndex])
		curr[maxIndex] -= sum
	}
	return result
}
