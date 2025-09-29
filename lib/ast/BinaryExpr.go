package ast

import (
	"hash/fnv"
	"logicka/lib/lexer"
	"logicka/lib/utils"
)

type BinaryExpr struct {
	Operator    lexer.TokenType
	Left, Right Expr
}

func NewBinaryExpr(op lexer.TokenType, left, right Expr) *BinaryExpr {
	return &BinaryExpr{op, left, right}
}

func (n *BinaryExpr) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("binary"))
	h.Write([]byte(n.Operator.String()))

	if n.Operator == lexer.BlImplication {
		h.Write(utils.Uint64ToBytes(n.Left.Hash()))
		h.Write(utils.Uint64ToBytes(n.Right.Hash()))
		return h.Sum64()
	}

	leftHash := n.Left.Hash()
	rightHash := n.Right.Hash()

	if leftHash < rightHash {
		h.Write(utils.Uint64ToBytes(leftHash))
		h.Write(utils.Uint64ToBytes(rightHash))
	} else {
		h.Write(utils.Uint64ToBytes(rightHash))
		h.Write(utils.Uint64ToBytes(leftHash))
	}
	return h.Sum64()
}

func (n *BinaryExpr) Equals(other Node) bool {
	return n.Hash() == other.Hash()
}

func (n *BinaryExpr) String() string {
	return n.Left.String() + " " + n.Operator.String() + " " + n.Right.String()
}

func (n *BinaryExpr) Children() []Expr {
	return []Expr{n.Left, n.Right}
}

func (n *BinaryExpr) expr() {}
