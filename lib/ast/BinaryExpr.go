package ast

import "logicka/lib/lexer"

type BinaryExpr struct {
	Operator    lexer.Token
	Left, Right Expr
}

func NewBinaryExpr(op lexer.Token, left, right Expr) *BinaryExpr {
	return &BinaryExpr{op, left, right}
}

func (n *BinaryExpr) String() string {
	return n.Left.String() + " " + n.Operator.Type.String() + " " + n.Right.String()
}

func (n *BinaryExpr) Children() []Expr {
	return []Expr{n.Left, n.Right}
}

func (n *BinaryExpr) expr() {}
