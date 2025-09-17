package ast

import (
	"fmt"
	"hash/fnv"
	"logicka/lib/utils"
)

// GroupingNode represents a parenthesized expression
type GroupingNode struct {
	Expr ASTNode
}

func NewGroupingNode(expr ASTNode) *GroupingNode {
	return &GroupingNode{Expr: expr}
}

func (n *GroupingNode) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("grouping"))
	h.Write(utils.Uint64ToBytes(n.Expr.Hash()))
	return h.Sum64()
}

func (n *GroupingNode) Equals(other ASTNode) bool {
	grouping, ok := other.(*GroupingNode)
	return ok && n.Hash() == grouping.Hash()
}

func (n *GroupingNode) String() string {
	return fmt.Sprintf("(%s)", n.Expr.String())
}

func (n *GroupingNode) Contains(node ASTNode) bool {
	container, ok := n.Expr.(Container)
	if ok {
		return container.Contains(node)
	}
	return n.Expr.Equals(node)
}
