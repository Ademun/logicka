package ast

import "strings"

type BracedExpression struct {
	Elements []Expr
}

func NewBracedExpression(elems []Expr) *BracedExpression {
	return &BracedExpression{Elements: elems}
}

func (n *BracedExpression) String() string {
	var results []string
	for _, stmt := range n.Elements {
		results = append(results, stmt.String())
	}
	return "{" + strings.Join(results, ", ") + "}"
}

func (n *BracedExpression) Children() []Expr {
	return n.Elements
}

func (n *BracedExpression) expr() {}
