package ast

type IdentifierExpr struct {
	Value string
}

func NewIdentifierExpr(val string) *IdentifierExpr {
	return &IdentifierExpr{val}
}

func (n *IdentifierExpr) String() string {
	return n.Value
}

func (n *IdentifierExpr) Children() []Expr {
	return []Expr{}
}

func (n *IdentifierExpr) expr() {}
