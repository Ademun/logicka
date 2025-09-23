package ast

type BlockStmt struct {
	Body []Stmt
}

func NewBlockStmt(body []Stmt) *BlockStmt {
	return &BlockStmt{Body: body}
}

func (n *BlockStmt) stmt() {}
