package ast

import (
	"hash/fnv"
	"logicka/lib/utils"
)

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

func (n *AssignmentExpr) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("assignment"))
	h.Write(utils.Uint64ToBytes(n.Lval.Hash()))
	h.Write(utils.Uint64ToBytes(n.Rval.Hash()))
	return h.Sum64()
}

func (n *AssignmentExpr) Equals(other Node) bool {
	return n.Hash() == other.Hash()
}

func (n *AssignmentExpr) String() string {
	return n.Lval.String() + " := " + n.Rval.String()
}

func (n *AssignmentExpr) expr() {}
