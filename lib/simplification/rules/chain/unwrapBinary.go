package chain

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"logicka/lib/simplification/rules/base"
)

type UnwrapBinaryRule struct {
	base.BaseRule
}

func NewUnwrapBinaryRule() *UnwrapBinaryRule {
	return &UnwrapBinaryRule{
		BaseRule: *base.NewBaseRule("UnwrapBinary law"),
	}
}

func (r *UnwrapBinaryRule) CanApply(node ast.ASTNode) bool {
	binary, ok := node.(*ast.BinaryNode)
	return ok && canFlatten(binary)
}

func (r *UnwrapBinaryRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	binary := node.(*ast.BinaryNode)
	operands := collectOperands(binary.Operator, binary)
	if len(operands) == 2 {
		return binary, nil
	}
	return ast.NewChainNode(binary.Operator, operands...)
}

func canFlatten(node *ast.BinaryNode) bool {
	if node.Operator != lexer.CONJ && node.Operator != lexer.DISJ {
		return false
	}
	return hasNestedSameOperator(node, node.Operator)
}

func isSameOperatorBinary(node ast.ASTNode, operator lexer.BooleanTokenType) bool {
	switch n := node.(type) {
	case *ast.BinaryNode:
		return n.Operator == operator
	case *ast.ChainNode:
		return n.Operator == operator
	case *ast.GroupingNode:
		if bin, ok := n.Expr.(*ast.BinaryNode); ok {
			return bin.Operator == operator
		}
		if chain, ok := n.Expr.(*ast.ChainNode); ok {
			return chain.Operator == operator
		}
	}
	return false
}

func hasNestedSameOperator(node ast.ASTNode, operator lexer.BooleanTokenType) bool {
	switch n := node.(type) {
	case *ast.BinaryNode:
		if n.Operator == operator {
			return isSameOperatorBinary(n.Left, operator) ||
				isSameOperatorBinary(n.Right, operator) ||
				hasNestedSameOperator(n.Left, operator) ||
				hasNestedSameOperator(n.Right, operator)
		}
	case *ast.GroupingNode:
		return hasNestedSameOperator(n.Expr, operator)
	}
	return false
}

func collectOperands(operator lexer.BooleanTokenType, node ast.ASTNode) []ast.ASTNode {
	operands := make([]ast.ASTNode, 0)
	switch n := node.(type) {
	case *ast.BinaryNode:
		if operator == n.Operator {
			operands = append(operands, collectOperands(operator, n.Left)...)
			operands = append(operands, collectOperands(operator, n.Right)...)
		} else {
			operands = append(operands, ast.NewGroupingNode(n))
		}
	case *ast.ChainNode:
		if operator == n.Operator {
			operands = append(operands, n.Operands...)
		} else {
			operands = append(operands, ast.NewGroupingNode(n))
		}
	case *ast.GroupingNode:
		operands = append(operands, collectOperands(operator, n.Expr)...)
	default:
		operands = append(operands, n)
	}
	return operands
}
