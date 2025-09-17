package ast

import (
	"hash/fnv"
	"strconv"
)

// LiteralNode represents a boolean literal like 1 or 0
type LiteralNode struct {
	Value bool
}

func NewLiteralNode(value bool) *LiteralNode {
	return &LiteralNode{value}
}

func (n *LiteralNode) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("literal"))
	h.Write([]byte(strconv.FormatBool(n.Value)))
	return h.Sum64()
}

func (n *LiteralNode) Equals(other ASTNode) bool {
	literal, ok := other.(*VariableNode)
	return ok && n.Hash() == literal.Hash()
}

func (n *LiteralNode) String() string {
	if n.Value {
		return "1"
	}
	return "0"
}

func IsTrue(node ASTNode) bool {
	literal, ok := node.(*LiteralNode)
	return ok && literal.Value
}

func IsFalse(node ASTNode) bool {
	literal, ok := node.(*LiteralNode)
	return ok && !literal.Value
}
