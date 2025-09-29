package ast

import (
	"hash/fnv"
	"logicka/lib/lexer"
	"logicka/lib/utils"
)

type UnaryExpr struct {
	Operator lexer.TokenType
	Operand  Expr
}

func NewUnaryExpr(op lexer.TokenType, operand Expr) *UnaryExpr {
	return &UnaryExpr{op, operand}
}

func (n *UnaryExpr) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("unary"))
	h.Write([]byte(n.Operator.String()))
	h.Write(utils.Uint64ToBytes(n.Operand.Hash()))
	return h.Sum64()
}

func (n *UnaryExpr) Equals(other Node) bool {
	return n.Hash() == other.Hash()
}

func (n *UnaryExpr) String() string {
	return n.Operator.String() + n.Operand.String()
}

func (n *UnaryExpr) expr() {}
