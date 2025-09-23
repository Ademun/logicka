package ast

type IdentifierExpr struct {
	Value string
}

func NewIdentifierExpr(val string) *IdentifierExpr {
	return &IdentifierExpr{val}
}

func (n *IdentifierExpr) expr() {}
