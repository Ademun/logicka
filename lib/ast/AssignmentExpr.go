package ast

type AssignmentExpr struct {
	Lval Expr
	Rval Expr
}

func NewAssignmentExpr(lval Expr, rval Expr) *AssignmentExpr {
	return &AssignmentExpr{Lval: lval, Rval: rval}
}

func (n *AssignmentExpr) expr() {}
