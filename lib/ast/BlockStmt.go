package ast

import "strings"

type BlockStmt struct {
	Body []Stmt
}

func NewBlockStmt(body []Stmt) *BlockStmt {
	return &BlockStmt{Body: body}
}

func (n *BlockStmt) String() string {
	var results []string
	for _, stmt := range n.Body {
		results = append(results, stmt.String())
	}
	return strings.Join(results, "\n")
}

func (n *BlockStmt) stmt() {}
