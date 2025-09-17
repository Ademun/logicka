package ast

import (
	"fmt"
	"hash/fnv"
	"logicka/lib/lexer"
	"logicka/lib/utils"
)

// UnaryNode represents a unary logical operation like negation
type UnaryNode struct {
	Operator lexer.BooleanTokenType
	Operand  ASTNode
}

func NewUnaryNode(operator lexer.BooleanTokenType, operand ASTNode) *UnaryNode {
	return &UnaryNode{operator, operand}
}

func (n *UnaryNode) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("unary"))
	h.Write([]byte(n.Operator.String()))
	h.Write(utils.Uint64ToBytes(n.Operand.Hash()))
	return h.Sum64()
}

func (n *UnaryNode) Equals(other ASTNode) bool {
	unary, ok := other.(*UnaryNode)
	return ok && n.Hash() == unary.Hash()
}

func (n *UnaryNode) String() string {
	return fmt.Sprintf("%s%s", n.Operator.String(), n.Operand.String())
}

func (n *UnaryNode) Contains(node ASTNode) bool {
	return n.Operand.Equals(node)
}

func (n *UnaryNode) IsType(unaryType lexer.BooleanTokenType) bool {
	return n.Operator == unaryType
}

// IsNegationOf returns true if the left node is the negation of the right node.
func IsNegationOf(left, right ASTNode) bool {
	unary, ok := left.(*UnaryNode)
	return ok && unary.Operator == lexer.NEG && unary.Operand.Equals(right)
}
