package ast

type ExprStmt struct {
	Expr Expr
}

func NewExprStmt(expr Expr) *ExprStmt {
	return &ExprStmt{Expr: expr}
}

func (n *ExprStmt) stmt() {}
