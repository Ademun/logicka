package ast

type AssignmentExpr struct {
	Lval Expr
	Rval Expr
}

func NewAssignmentExpr(lval Expr, rval Expr) *AssignmentExpr {
	return &AssignmentExpr{Lval: lval, Rval: rval}
}

func (n *AssignmentExpr) String() string {
	return n.Lval.String() + " := " + n.Rval.String()
}

func (n *AssignmentExpr) Children() []Expr {
	return []Expr{n.Lval, n.Rval}
}

func (n *AssignmentExpr) expr() {}
