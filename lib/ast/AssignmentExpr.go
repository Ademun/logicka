package ast

type AssignmentExpr struct {
	Lval Expr
	Rval Expr
}

func NewAssignmentExpr(lval Expr, rval Expr) (*AssignmentExpr, error) {
	switch lval.(type) {
	case *IdentifierExpr, *FunctionDeclExpr:
		break
	default:
		return nil, NewValidationError("Left side of assignment expression must be an identifier or a function declaration")
	}
	return &AssignmentExpr{Lval: lval, Rval: rval}, nil
}

func (n *AssignmentExpr) String() string {
	return n.Lval.String() + " := " + n.Rval.String()
}

func (n *AssignmentExpr) Children() []Expr {
	return []Expr{n.Lval, n.Rval}
}

func (n *AssignmentExpr) expr() {}
