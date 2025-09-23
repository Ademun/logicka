package ast

type BooleanExpr struct {
	Value bool
}

func NewBooleanExpr(val bool) *BooleanExpr {
	return &BooleanExpr{val}
}

func (n *BooleanExpr) expr() {}
