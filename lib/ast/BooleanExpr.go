package ast

import (
	"hash/fnv"
	"strconv"
)

type BooleanExpr struct {
	Value bool
}

func NewBooleanExpr(val bool) *BooleanExpr {
	return &BooleanExpr{val}
}

func (n *BooleanExpr) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("boolean"))
	h.Write([]byte(strconv.FormatBool(n.Value)))
	return h.Sum64()
}

func (n *BooleanExpr) Equals(other Node) bool {
	return n.Hash() == other.Hash()
}

func (n *BooleanExpr) String() string {
	return strconv.FormatBool(n.Value)
}

func (n *BooleanExpr) expr() {}
