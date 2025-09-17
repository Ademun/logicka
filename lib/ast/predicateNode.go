package ast

import (
	"fmt"
	"hash/fnv"
	"logicka/lib/utils"
	"slices"
	"strings"
)

// PredicateNode represents a predicate logic expression
type PredicateNode struct {
	Name string
	Args []*VariableNode
	Body ASTNode
}

func NewPredicateNode(name string, args []*VariableNode, body ASTNode) *PredicateNode {
	return &PredicateNode{Name: name, Args: args, Body: body}
}

func (n *PredicateNode) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("predicate"))
	h.Write([]byte(n.Name))

	hashes := make([]uint64, 0, len(n.Args))
	for _, arg := range n.Args {
		hashes = append(hashes, arg.Hash())
	}
	slices.Sort(hashes)

	for _, hash := range hashes {
		h.Write(utils.Uint64ToBytes(hash))
	}

	h.Write(utils.Uint64ToBytes(n.Body.Hash()))

	return h.Sum64()
}

func (n *PredicateNode) Equals(other ASTNode) bool {
	predicate, ok := other.(*PredicateNode)
	return ok && n.Hash() == predicate.Hash()
}

func (n *PredicateNode) String() string {
	argsStrings := make([]string, 0, len(n.Args))
	for _, arg := range n.Args {
		argsStrings = append(argsStrings, arg.String())
	}
	return fmt.Sprintf("%s(%s)", n.Name, strings.Join(argsStrings, ", "))
}

func (n *PredicateNode) FullString() string {
	return fmt.Sprintf("%s = %s", n.String(), n.Body.String())
}
