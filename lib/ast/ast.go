package ast

import (
	"logicka/lib/lexer"
)

type ASTNode interface {
	Equals(node ASTNode) bool
}

type GroupingNode struct {
	Expr ASTNode
}

func (g GroupingNode) Equals(node ASTNode) bool {
	same, ok := node.(GroupingNode)
	return ok && same.Expr.Equals(g.Expr)
}

type LiteralNode struct {
	Value bool
}

func (l LiteralNode) Equals(node ASTNode) bool {
	same, ok := node.(LiteralNode)
	return ok && same.Value == l.Value
}

type VariableNode struct {
	Name string
}

func (v VariableNode) Equals(node ASTNode) bool {
	same, ok := node.(VariableNode)
	return ok && same.Name == v.Name
}

type BinaryNode struct {
	Operator    lexer.TokenType
	Left, Right ASTNode
}

func (b BinaryNode) Equals(node ASTNode) bool {
	same, ok := node.(BinaryNode)
	return ok && b.Operator == same.Operator && b.Left.Equals(same.Left) && b.Right.Equals(same.Right)
}

type UnaryNode struct {
	Operator lexer.TokenType
	Operand  ASTNode
}

func (u UnaryNode) Equals(node ASTNode) bool {
	same, ok := node.(UnaryNode)
	return ok && u.Operator == same.Operator && u.Operand.Equals(same.Operand)
}

type PredicateNode struct {
	Name string
	Body any
}

func (p PredicateNode) Equals(node ASTNode) bool {
	return true
}

type QuantifierNode struct {
	Type     lexer.TokenType
	Variable string
	Domain   any
}

func (q QuantifierNode) Equals(node ASTNode) bool {
	return true
}
