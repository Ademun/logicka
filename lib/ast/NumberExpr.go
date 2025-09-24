package ast

import "strconv"

type NumberExpr struct {
	Value float64
}

func NewNumberExpr(val float64) *NumberExpr {
	return &NumberExpr{val}
}

func (n *NumberExpr) String() string {
	return strconv.FormatFloat(n.Value, 'g', -1, 64)
}

func (n *NumberExpr) Children() []Expr {
	return []Expr{}
}

func (n *NumberExpr) expr() {}
