package ast

import (
	"hash/fnv"
	"logicka/lib/utils"
)

type GroupingExpr struct {
	Body Expr
}

func NewGroupingExpr(body Expr) *GroupingExpr {
	return &GroupingExpr{Body: body}
}

func (n *GroupingExpr) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("grouping"))
	h.Write(utils.Uint64ToBytes(n.Body.Hash()))
	return h.Sum64()
}

func (n *GroupingExpr) Equals(other Node) bool {
	return n.Hash() == other.Hash()
}

func (n *GroupingExpr) String() string {
	return "(" + n.Body.String() + ")"
}

func (n *GroupingExpr) expr() {}
