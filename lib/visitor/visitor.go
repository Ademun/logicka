package visitor

import (
	"fmt"
	"logicka/lib/ast"
)

type UnknownNodeType struct {
	NodeType string
}

func (u UnknownNodeType) Error() string {
	return "unknown AST node type: " + u.NodeType
}

type UnknownTokenType struct {
	TokenType string
}

func (u UnknownTokenType) Error() string {
	return "unknown token type: " + u.TokenType
}

type Visitor[T any] interface {
	VisitGrouping(node *ast.GroupingNode) (T, error)
	VisitLiteral(node *ast.LiteralNode) (T, error)
	VisitVariable(node *ast.VariableNode) (T, error)
	VisitBinary(node *ast.BinaryNode) (T, error)
	VisitUnary(node *ast.UnaryNode) (T, error)
	VisitPredicate(node *ast.PredicateNode) (T, error)
	VisitQuantifier(node *ast.QuantifierNode) (T, error)
}

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
	case *ast.UnaryNode:
		return visitor.VisitUnary(n)
	case *ast.PredicateNode:
		return visitor.VisitPredicate(n)
	case *ast.QuantifierNode:
		return visitor.VisitQuantifier(n)
	default:
		var zero T
		return zero, UnknownNodeType{NodeType: fmt.Sprintf("%T", n)}
	}
}

type EvaluationContext struct {
	Variables map[string]bool
}
