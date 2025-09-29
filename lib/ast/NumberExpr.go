package ast

import (
	"hash/fnv"
	"strconv"
)

type NumberExpr struct {
	Value float64
}

func NewNumberExpr(val float64) *NumberExpr {
	return &NumberExpr{val}
}

func (n *NumberExpr) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("number"))
	h.Write([]byte(strconv.FormatFloat(n.Value, 'f', -1, 64)))
	return h.Sum64()
}

func (n *NumberExpr) Equals(other Node) bool {
	return n.Hash() == other.Hash()
}

func (n *NumberExpr) String() string {
	return strconv.FormatFloat(n.Value, 'g', -1, 64)
}

func (n *NumberExpr) expr() {}
