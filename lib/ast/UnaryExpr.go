package ast

import "logicka/lib/lexer"

type UnaryExpr struct {
	Operator lexer.Token
	Operand  Expr
}

func NewUnaryExpr(op lexer.Token, operand Expr) *UnaryExpr {
	return &UnaryExpr{op, operand}
}

func (n *UnaryExpr) expr() {}
