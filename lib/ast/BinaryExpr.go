package ast

import "logicka/lib/lexer"

type BinaryExpr struct {
	Operator    lexer.Token
	Left, Right Expr
}

func NewBinaryExpr(op lexer.Token, left, right Expr) *BinaryExpr {
	return &BinaryExpr{op, left, right}
}

func (n *BinaryExpr) expr() {}
