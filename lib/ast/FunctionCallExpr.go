package ast

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

func (n *FunctionCallExpr) expr() {}
