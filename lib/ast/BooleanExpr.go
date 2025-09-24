package ast

import "strconv"

type BooleanExpr struct {
	Value bool
}

func NewBooleanExpr(val bool) *BooleanExpr {
	return &BooleanExpr{val}
}

func (n *BooleanExpr) String() string {
	return strconv.FormatBool(n.Value)
}

func (n *BooleanExpr) Children() []Expr {
	return []Expr{}
}

func (n *BooleanExpr) expr() {}
