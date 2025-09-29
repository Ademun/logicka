package ast

import (
	"hash/fnv"
	"logicka/lib/utils"
	"strings"
)

type FunctionCallExpr struct {
	Name string
	Args []Expr
}

func NewFunctionCallExpr(name string, args []Expr) *FunctionCallExpr {
	return &FunctionCallExpr{
		Name: name,
		Args: args,
	}
}

func (n *FunctionCallExpr) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("function_call"))
	h.Write([]byte(n.Name))
	for _, arg := range n.Args {
		h.Write(utils.Uint64ToBytes(arg.Hash()))
	}
	return h.Sum64()
}

func (n *FunctionCallExpr) Equals(other Node) bool {
	return n.Hash() == other.Hash()
}

func (n *FunctionCallExpr) String() string {
	var results []string
	for _, arg := range n.Args {
		results = append(results, arg.String())
	}
	return n.Name + "(" + strings.Join(results, ", ") + ")"
}

func (n *FunctionCallExpr) expr() {}
