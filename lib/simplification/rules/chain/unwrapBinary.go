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
		BaseRule: *base.NewBaseRule("Объединение в цепочку операторов"),
	}
}

func (r *UnwrapBinaryRule) CanApply(node ast.Node) bool {
	binary, ok := node.(*ast.BinaryExpr)
	return ok && canFlatten(binary)
}

func (r *UnwrapBinaryRule) Apply(node ast.Node) (ast.Node, error) {
	binary := node.(*ast.BinaryExpr)
	operands := collectOperands(binary.Operator, binary)
	if len(operands) == 2 {
		return binary, nil
	}
	return ast.NewChainExpr(binary.Operator, operands), nil
}

func canFlatten(node *ast.BinaryExpr) bool {
	if node.Operator != lexer.BlConjunction && node.Operator != lexer.BlDisjunction {
		return false
	}
	return hasNestedSameOperator(node, node.Operator)
}

func isSameOperatorBinary(node ast.Node, operator lexer.TokenType) bool {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		return n.Operator == operator
	case *ast.ChainExpr:
		return n.Operator == operator
	case *ast.GroupingExpr:
		if bin, ok := n.Body.(*ast.BinaryExpr); ok {
			return bin.Operator == operator
		}
		if chain, ok := n.Body.(*ast.ChainExpr); ok {
			return chain.Operator == operator
		}
	}
	return false
}

func hasNestedSameOperator(node ast.Node, operator lexer.TokenType) bool {
	switch n := node.(type) {
	case *ast.BinaryExpr:
		if n.Operator == operator {
			return isSameOperatorBinary(n.Left, operator) ||
				isSameOperatorBinary(n.Right, operator) ||
				hasNestedSameOperator(n.Left, operator) ||
				hasNestedSameOperator(n.Right, operator)
		}
	case *ast.GroupingExpr:
		return hasNestedSameOperator(n.Body, operator)
	}
	return false
}

func collectOperands(operator lexer.TokenType, node ast.Expr) []ast.Expr {
	operands := make([]ast.Expr, 0)
	switch n := node.(type) {
	case *ast.BinaryExpr:
		if operator == n.Operator {
			operands = append(operands, collectOperands(operator, n.Left)...)
			operands = append(operands, collectOperands(operator, n.Right)...)
		} else {
			operands = append(operands, ast.NewGroupingExpr(n))
		}
	case *ast.ChainExpr:
		if operator == n.Operator {
			operands = append(operands, n.Elements...)
		} else {
			operands = append(operands, ast.NewGroupingExpr(n))
		}
	case *ast.GroupingExpr:
		operands = append(operands, collectOperands(operator, n.Body)...)
	default:
		operands = append(operands, n)
	}
	return operands
}
