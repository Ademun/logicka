package ast

type AssignmentStmt struct {
	Lval Expr
	Rval Expr
}

func NewAssignmentStmt(lval Expr, rval Expr) *AssignmentStmt {
	return &AssignmentStmt{Lval: lval, Rval: rval}
}

func (n *AssignmentStmt) stmt() {}
