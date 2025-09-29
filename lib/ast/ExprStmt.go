package ast

import (
	"hash/fnv"
	"logicka/lib/utils"
)

type ExprStmt struct {
	Expr Expr
}

func NewExprStmt(expr Expr) *ExprStmt {
	return &ExprStmt{Expr: expr}
}

func (n *ExprStmt) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("expr_stmt"))
	h.Write(utils.Uint64ToBytes(n.Expr.Hash()))
	return h.Sum64()
}

func (n *ExprStmt) Equals(other Node) bool {
	return n.Hash() == other.Hash()
}

func (n *ExprStmt) String() string {
	return n.Expr.String()
}

func (n *ExprStmt) stmt() {}
