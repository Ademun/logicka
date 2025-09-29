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

type Visitor[T any] interface {
	VisitAssignmentExpr(node *ast.AssignmentExpr) (T, error)
	VisitBinaryExpr(node *ast.BinaryExpr) (T, error)
	VisitBlockStmt(node *ast.BlockStmt) (T, error)
	VisitBooleanExpr(node *ast.BooleanExpr) (T, error)
	VisitBracedExpr(node *ast.BracedExpr) (T, error)
	VisitChainExpr(node *ast.ChainExpr) (T, error)
	VisitExprStmt(node *ast.ExprStmt) (T, error)
	VisitFunctionCallExpr(node *ast.FunctionCallExpr) (T, error)
	VisitFunctionDeclExpr(node *ast.FunctionDeclExpr) (T, error)
	VisitGroupingExpr(node *ast.GroupingExpr) (T, error)
	VisitIdentifierExpr(node *ast.IdentifierExpr) (T, error)
	VisitNumberExpr(node *ast.NumberExpr) (T, error)
	VisitQuantifierExpr(node *ast.QuantifierExpr) (T, error)
	VisitStringExpr(node *ast.StringExpr) (T, error)
	VisitUnaryExpr(node *ast.UnaryExpr) (T, error)
}

func Accept[T any](node ast.Node, visitor Visitor[T]) (T, error) {
	switch n := node.(type) {
	case *ast.AssignmentExpr:
		return visitor.VisitAssignmentExpr(n)
	case *ast.BinaryExpr:
		return visitor.VisitBinaryExpr(n)
	case *ast.BlockStmt:
		return visitor.VisitBlockStmt(n)
	case *ast.BooleanExpr:
		return visitor.VisitBooleanExpr(n)
	case *ast.BracedExpr:
		return visitor.VisitBracedExpr(n)
	case *ast.ChainExpr:
		return visitor.VisitChainExpr(n)
	case *ast.ExprStmt:
		return visitor.VisitExprStmt(n)
	case *ast.FunctionCallExpr:
		return visitor.VisitFunctionCallExpr(n)
	case *ast.FunctionDeclExpr:
		return visitor.VisitFunctionDeclExpr(n)
	case *ast.GroupingExpr:
		return visitor.VisitGroupingExpr(n)
	case *ast.IdentifierExpr:
		return visitor.VisitIdentifierExpr(n)
	case *ast.NumberExpr:
		return visitor.VisitNumberExpr(n)
	case *ast.QuantifierExpr:
		return visitor.VisitQuantifierExpr(n)
	case *ast.StringExpr:
		return visitor.VisitStringExpr(n)
	case *ast.UnaryExpr:
		return visitor.VisitUnaryExpr(n)
	default:
		var zero T
		return zero, NodeTypeError{NodeType: fmt.Sprintf("%T", n)}
	}
}
