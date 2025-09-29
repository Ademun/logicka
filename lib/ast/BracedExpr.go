package ast

import (
	"hash/fnv"
	"logicka/lib/utils"
	"sort"
	"strings"
)

type BracedExpr struct {
	Elements []Expr
}

func NewBracedExpr(elems []Expr) *BracedExpr {
	return &BracedExpr{Elements: elems}
}

func (n *BracedExpr) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("braced"))

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

func (n *BracedExpr) Equals(other Node) bool {
	return n.Hash() == other.Hash()
}

func (n *BracedExpr) String() string {
	var results []string
	for _, stmt := range n.Elements {
		results = append(results, stmt.String())
	}
	return "{" + strings.Join(results, ", ") + "}"
}

func (n *BracedExpr) Children() []Expr {
	return n.Elements
}

func (n *BracedExpr) expr() {}
