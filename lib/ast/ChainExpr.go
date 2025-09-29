package ast

import (
	"hash/fnv"
	"logicka/lib/lexer"
	"logicka/lib/utils"
	"strings"
)

type ChainExpr struct {
	Operator lexer.TokenType
	Elements []Expr
}

func NewChainExpr(elems []Expr) *ChainExpr {
	return &ChainExpr{Elements: elems}
}

func (n *ChainExpr) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("chain"))
	for _, elem := range n.Elements {
		h.Write(utils.Uint64ToBytes(elem.Hash()))
	}
	return h.Sum64()
}

func (n *ChainExpr) Equals(other Node) bool {
	return n.Hash() == other.Hash()
}

func (n *ChainExpr) String() string {
	var results []string
	for _, stmt := range n.Elements {
		results = append(results, stmt.String())
	}
	return "[" + strings.Join(results, ", ") + "]"
}

func (n *ChainExpr) Children() []Expr {
	return n.Elements
}

func (n *ChainExpr) expr() {}
