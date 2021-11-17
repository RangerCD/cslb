package strategy

import (
	"github.com/RangerCD/cslb/node"
)

// Strategy controls how the nodes are chosen
// This type should be thread safe
type Strategy interface {
	// SetNodes update saved nodes
	SetNodes(nodes []node.Node)
	// Next returns a node address
	Next() (node.Node, error)
}
