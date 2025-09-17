package ast

import (
	"fmt"
	"hash/fnv"
	"logicka/lib/lexer"
	"logicka/lib/utils"
)

// BinaryNode represents a binary logical operation like conjunction, disjunction, etc.
type BinaryNode struct {
	Operator    lexer.BooleanTokenType
	Left, Right ASTNode
}

func NewBinaryNode(operator lexer.BooleanTokenType, left, right ASTNode) *BinaryNode {
	return &BinaryNode{Operator: operator, Left: left, Right: right}
}

func (n *BinaryNode) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("binary"))
	h.Write([]byte(n.Operator.String()))

	lHash := n.Left.Hash()
	rHash := n.Right.Hash()

	// Implication is a non-associative operation
	if n.Operator == lexer.IMPL {
		h.Write(utils.Uint64ToBytes(lHash))
		h.Write(utils.Uint64ToBytes(rHash))
		return h.Sum64()
	}

	// Associative hash order
	if lHash < rHash {
		h.Write(utils.Uint64ToBytes(lHash))
		h.Write(utils.Uint64ToBytes(rHash))
	} else {
		h.Write(utils.Uint64ToBytes(rHash))
		h.Write(utils.Uint64ToBytes(lHash))
	}

	return h.Sum64()
}

func (n *BinaryNode) Equals(other ASTNode) bool {
	binary, ok := other.(*BinaryNode)
	return ok && n.Hash() == binary.Hash()
}

func (n *BinaryNode) String() string {
	return fmt.Sprintf("%s %s %s",
		n.Left.String(),
		n.Operator.String(),
		n.Right.String())
}

func (n *BinaryNode) Children() []ASTNode {
	return []ASTNode{n.Left, n.Right}
}

func (n *BinaryNode) Contains(node ASTNode) bool {
	return n.Left.Equals(node) || n.Right.Equals(node)
}

func (n *BinaryNode) IsType(binType lexer.BooleanTokenType) bool {
	return n.Operator == binType
}
