package ast

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

func (n *FunctionDeclExpr) expr() {}
