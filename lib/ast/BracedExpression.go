package ast

type BracedExpression struct {
	Elements []Expr
}

func NewBracedExpression(elems []Expr) *BracedExpression {
	return &BracedExpression{Elements: elems}
}

func (n *BracedExpression) expr() {}
