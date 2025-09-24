package ast

import "strings"

type FunctionDeclExpr struct {
	Name string
	Args []Expr
	Body Expr
}

func NewFunctionDeclExpr(name string, args []Expr, body Expr) (*FunctionDeclExpr, error) {
	for _, arg := range args {
		if _, ok := arg.(*IdentifierExpr); !ok {
			return nil, NewValidationError("Arguments of function declaration must be identifiers")
		}
	}

	return &FunctionDeclExpr{
		Name: name,
		Args: args,
		Body: body,
	}, nil
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
