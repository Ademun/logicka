package ast

import "logicka/lib/lexer"

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

func IsTrue(node Expr) bool {
	boolean, ok := node.(*BooleanExpr)
	return ok && boolean.Value
}

func IsFalse(node Expr) bool {
	boolean, ok := node.(*BooleanExpr)
	return ok && !boolean.Value
}

func IsNegationOf(left, right Expr) bool {
	unary, ok := left.(*UnaryExpr)
	return ok && unary.Operator == lexer.BlNegation && unary.Operand.Equals(right)
}
