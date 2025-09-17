package ast

import (
	"hash/fnv"
	"logicka/lib/lexer"
	"logicka/lib/utils"
	"slices"
	"strings"
)

// ChainNode represent a chain of binary operations
// It is required for proper simplification, as it reduces the tree depth
type ChainNode struct {
	Operator lexer.BooleanTokenType
	Operands []ASTNode
}

func NewChainNode(operator lexer.BooleanTokenType, operands ...ASTNode) *ChainNode {
	return &ChainNode{Operator: operator, Operands: slices.Clone(operands)}
}

func (n *ChainNode) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("chain"))
	h.Write([]byte(n.Operator.String()))

	hashes := make([]uint64, 0, len(n.Operands))
	for _, op := range n.Operands {
		hashes = append(hashes, op.Hash())
	}

	// Preserve order for non-associative implication
	if !n.IsType(lexer.IMPL) {
		slices.Sort(hashes)
	}

	for _, hash := range hashes {
		h.Write(utils.Uint64ToBytes(hash))
	}

	return h.Sum64()
}

func (n *ChainNode) Equals(other ASTNode) bool {
	chain, ok := other.(*ChainNode)
	return ok && n.Hash() == chain.Hash()
}

func (n *ChainNode) String() string {
	parts := make([]string, 0, len(n.Operands))
	for _, operand := range n.Operands {
		if operand != nil {
			parts = append(parts, operand.String())
		}
	}

	return strings.Join(parts, " "+n.Operator.String()+" ")
}

func (n *ChainNode) Children() []ASTNode {
	return n.Operands
}

func (n *ChainNode) Contains(node ASTNode) bool {
	return slices.ContainsFunc(n.Operands, node.Equals)
}

func (n *ChainNode) Add(nodes ...ASTNode) {
	n.Operands = append(n.Operands, nodes...)
}

func (n *ChainNode) Remove(node ASTNode) {
	n.Operands = slices.DeleteFunc(n.Operands, node.Equals)
}

func (n *ChainNode) IsType(chainType lexer.BooleanTokenType) bool {
	return n.Operator == chainType
}
