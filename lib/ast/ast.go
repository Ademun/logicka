package ast

import "fmt"

type Stmt interface {
	fmt.Stringer
	stmt()
}

type Expr interface {
	fmt.Stringer
	Children() []Expr
	expr()
}
