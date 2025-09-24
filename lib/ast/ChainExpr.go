package ast

import (
	"logicka/lib/lexer"
	"strings"
)

type ChainExpr struct {
	Operator lexer.TokenType
	Elements []Expr
}

func NewChainExpr(elems []Expr) *ChainExpr {
	return &ChainExpr{Elements: elems}
}

func (n *ChainExpr) String() string {
	var results []string
	for _, stmt := range n.Elements {
		results = append(results, stmt.String())
	}
	return "[" + strings.Join(results, ", ") + "]"
}

func (n *ChainExpr) Children() []Expr {
	return n.Elements
}

func (n *ChainExpr) expr() {}
