package ast

import "strings"

type FunctionDeclExpr struct {
	Name string
	Args []Expr
	Body Expr
}

func NewFunctionDeclExpr(name string, args []Expr, body Expr) *FunctionDeclExpr {
	return &FunctionDeclExpr{
		Name: name,
		Args: args,
		Body: body,
	}
}

func (n *FunctionDeclExpr) String() string {
	var results []string
	for _, arg := range n.Args {
		results = append(results, arg.String())
	}
	return n.Name + "(" + strings.Join(results, ", ") + ")" + " := " + n.Body.String()
}

func (n *FunctionDeclExpr) Children() []Expr {
	return []Expr{n.Body}
}

func (n *FunctionDeclExpr) expr() {}
