package ast

import "hash/fnv"

type StringExpr struct {
	Value string
}

func NewStringExpr(value string) *StringExpr {
	return &StringExpr{Value: value}
}

func (n *StringExpr) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("string"))
	h.Write([]byte(n.Value))
	return h.Sum64()
}

func (n *StringExpr) Equals(other Node) bool {
	return n.Hash() == other.Hash()
}

func (n *StringExpr) String() string {
	return "\"" + n.Value + "\""
}

func (n *StringExpr) Children() []Expr {
	return []Expr{}
}

func (n *StringExpr) expr() {}
