package ast

import "hash/fnv"

type IdentifierExpr struct {
	Name string
}

func NewIdentifierExpr(val string) *IdentifierExpr {
	return &IdentifierExpr{val}
}

func (n *IdentifierExpr) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("identifier"))
	h.Write([]byte(n.Name))
	return h.Sum64()
}

func (n *IdentifierExpr) Equals(other Node) bool {
	return n.Hash() == other.Hash()
}

func (n *IdentifierExpr) String() string {
	return n.Name
}

func (n *IdentifierExpr) Children() []Expr {
	return []Expr{}
}

func (n *IdentifierExpr) expr() {}
