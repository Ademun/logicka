package ast

import "strings"

type FunctionCallExpr struct {
	Name string
	Args []Expr
}

func NewFunctionCallExpr(name string, args []Expr) *FunctionCallExpr {
	return &FunctionCallExpr{
		Name: name,
		Args: args,
	}
}

func (n *FunctionCallExpr) String() string {
	var results []string
	for _, arg := range n.Args {
		results = append(results, arg.String())
	}
	return n.Name + "(" + strings.Join(results, ", ") + ")"
}

func (n *FunctionCallExpr) Children() []Expr {
	return n.Args
}

func (n *FunctionCallExpr) expr() {}
