package ast

type NumberExpr struct {
	Value float64
}

func NewNumberExpr(val float64) *NumberExpr {
	return &NumberExpr{val}
}

func (n *NumberExpr) expr() {}
