package ast

type StringExpr struct {
	Value string
}

func NewStringExpr(value string) *StringExpr {
	return &StringExpr{Value: value}
}

func (n *StringExpr) String() string {
	return "\"" + n.Value + "\""
}

func (n *StringExpr) Children() []Expr {
	return []Expr{}
}

func (n *StringExpr) expr() {}
