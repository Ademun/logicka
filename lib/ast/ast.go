package ast

import (
	"fmt"
	"logicka/lib/lexer"
)

type ASTNode interface {
	Equals(node ASTNode) bool
	String() string
}

type GroupingNode struct {
	Expr ASTNode
}

func (g GroupingNode) Equals(node ASTNode) bool {
	same, ok := node.(GroupingNode)
	return ok && same.Expr.Equals(g.Expr)
}

func (g GroupingNode) String() string {
	return fmt.Sprintf("(%s)", g.Expr.String())
}

type LiteralNode struct {
	Value bool
}

func (l LiteralNode) Equals(node ASTNode) bool {
	same, ok := node.(*LiteralNode)
	return ok && same.Value == l.Value
}

func (l LiteralNode) String() string {
	return fmt.Sprint(l.Value)
}

type VariableNode struct {
	Name string
}

func (v VariableNode) Equals(node ASTNode) bool {
	same, ok := node.(*VariableNode)
	return ok && same.Name == v.Name
}

func (v VariableNode) String() string {
	return v.Name
}

type BinaryNode struct {
	Operator    lexer.TokenType
	Left, Right ASTNode
}

func (b BinaryNode) Equals(node ASTNode) bool {
	same, ok := node.(*BinaryNode)
	return ok && b.Operator == same.Operator && b.Left.Equals(same.Left) && b.Right.Equals(same.Right)
}

func (b BinaryNode) String() string {
	return fmt.Sprint(b.Left.String(), " ", b.Operator.String(), " ", b.Right.String())
}

type UnaryNode struct {
	Operator lexer.TokenType
	Operand  ASTNode
}

func (u UnaryNode) Equals(node ASTNode) bool {
	same, ok := node.(*UnaryNode)
	return ok && u.Operator == same.Operator && u.Operand.Equals(same.Operand)
}

func (u UnaryNode) String() string {
	return fmt.Sprint(u.Operator.String(), u.Operand.String())
}

type PredicateNode struct {
	Name string
	Body any
}

func (p PredicateNode) Equals(node ASTNode) bool {
	return true
}

func (p PredicateNode) String() string {
	return p.Name
}

type QuantifierNode struct {
	Type     lexer.TokenType
	Variable string
	Domain   any
}

func (q QuantifierNode) Equals(node ASTNode) bool {
	return true
}

func (q QuantifierNode) String() string {
	return fmt.Sprint(q.Variable, q.Domain)
}
