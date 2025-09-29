package ast

type Node interface {
	String() string
	Hash() uint64
	Equals(other Node) bool
}

type Stmt interface {
	Node
	stmt()
}

type Expr interface {
	Node
	Children() []Expr
	expr()
}
