package ast

import (
	"hash/fnv"
	"logicka/lib/lexer"
	"logicka/lib/utils"
	"slices"
	"sort"
	"strings"
)

type ChainExpr struct {
	Operator lexer.TokenType
	Elements []Expr
}

func NewChainExpr(operator lexer.TokenType, elems []Expr) *ChainExpr {
	return &ChainExpr{Operator: operator, Elements: elems}
}

func (n *ChainExpr) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("chain"))

	if n.Operator == lexer.BlConjunction || n.Operator == lexer.BlDisjunction {
		hashes := make([]uint64, len(n.Elements))
		for i, e := range n.Elements {
			hashes[i] = e.Hash()
		}
		sort.Slice(hashes, func(i, j int) bool {
			return hashes[i] < hashes[j]
		})

		for _, hash := range hashes {
			h.Write(utils.Uint64ToBytes(hash))
		}
		return h.Sum64()
	}

	for _, e := range n.Elements {
		h.Write(utils.Uint64ToBytes(e.Hash()))
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

func (n *ChainExpr) Contains(node Expr) bool {
	return slices.ContainsFunc(n.Elements, func(e Expr) bool {
		return e.Equals(node)
	})
}

func (n *ChainExpr) Remove(node Expr) bool {
	originalLen := len(n.Elements)
	n.Elements = slices.DeleteFunc(n.Elements, func(e Expr) bool {
		return e.Equals(node)
	})
	return len(n.Elements) != originalLen
}

func (n *ChainExpr) expr() {}
