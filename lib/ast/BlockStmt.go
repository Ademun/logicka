package ast

import (
	"hash/fnv"
	"logicka/lib/utils"
	"strings"
)

type BlockStmt struct {
	Body []Stmt
}

func NewBlockStmt(body []Stmt) *BlockStmt {
	return &BlockStmt{Body: body}
}

func (n *BlockStmt) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("block_stmt"))
	for _, stmt := range n.Body {
		h.Write(utils.Uint64ToBytes(stmt.Hash()))
	}
	return h.Sum64()
}

func (n *BlockStmt) Equals(other Node) bool {
	return n.Hash() == other.Hash()
}

func (n *BlockStmt) String() string {
	var results []string
	for _, stmt := range n.Body {
		results = append(results, stmt.String())
	}
	return strings.Join(results, "\n")
}

func (n *BlockStmt) stmt() {}
