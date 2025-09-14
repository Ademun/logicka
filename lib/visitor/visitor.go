// Package visitor provides the visitor pattern implementation for AST traversal.
package visitor

import (
	"fmt"
	"logicka/lib/ast"
)

// Custom error types for better error handling
type NodeTypeError struct {
	NodeType string
}

func (e NodeTypeError) Error() string {
	return fmt.Sprintf("unknown AST node type: %s", e.NodeType)
}

type OperatorError struct {
	Operator string
}

func (e OperatorError) Error() string {
	return fmt.Sprintf("unknown operator: %s", e.Operator)
}

// Visitor defines the interface for AST node visitors.
type Visitor[T any] interface {
	VisitGrouping(node *ast.GroupingNode) (T, error)
	VisitLiteral(node *ast.LiteralNode) (T, error)
	VisitVariable(node *ast.VariableNode) (T, error)
	VisitBinary(node *ast.BinaryNode) (T, error)
	VisitChain(node *ast.ChainNode) (T, error)
	VisitUnary(node *ast.UnaryNode) (T, error)
	VisitPredicate(node *ast.PredicateNode) (T, error)
	VisitQuantifier(node *ast.QuantifierNode) (T, error)
}

// Accept dispatches the appropriate visitor method based on the node type.
func Accept[T any](node ast.ASTNode, visitor Visitor[T]) (T, error) {
	switch n := node.(type) {
	case *ast.GroupingNode:
		return visitor.VisitGrouping(n)
	case *ast.LiteralNode:
		return visitor.VisitLiteral(n)
	case *ast.VariableNode:
		return visitor.VisitVariable(n)
	case *ast.BinaryNode:
		return visitor.VisitBinary(n)
	case *ast.ChainNode:
		return visitor.VisitChain(n)
	case *ast.UnaryNode:
		return visitor.VisitUnary(n)
	case *ast.PredicateNode:
		return visitor.VisitPredicate(n)
	case *ast.QuantifierNode:
		return visitor.VisitQuantifier(n)
	default:
		var zero T
		return zero, NodeTypeError{NodeType: fmt.Sprintf("%T", n)}
	}
}

// EvaluationContext holds variable assignments for expression evaluation.
type EvaluationContext struct {
	Variables map[string]bool
}

func NewEvaluationContext() *EvaluationContext {
	return &EvaluationContext{
		Variables: make(map[string]bool),
	}
}

func (ctx *EvaluationContext) SetVariable(name string, value bool) {
	ctx.Variables[name] = value
}

func (ctx *EvaluationContext) GetVariable(name string) (bool, bool) {
	value, exists := ctx.Variables[name]
	return value, exists
}
