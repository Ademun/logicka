package complex

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"logicka/lib/simplification/rules/base"
)

type AbsorptionRule struct {
	base.BaseRule
}

func NewAbsorptionRule() *AbsorptionRule {
	return &AbsorptionRule{
		BaseRule: *base.NewBaseRule("Absorption law"),
	}
}

func (r *AbsorptionRule) CanApply(node ast.ASTNode) bool {
	binary, ok := node.(*ast.BinaryNode)
	if !ok {
		return false
	}

	if binary.Operator != lexer.CONJ && binary.Operator != lexer.DISJ {
		return false
	}

	_, lok := binary.Left.(*ast.GroupingNode)
	_, rok := binary.Right.(*ast.GroupingNode)

	return lok != rok
}

func (r *AbsorptionRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	binary := node.(*ast.BinaryNode)

	if result, absorbed := r.applyAbsorption(binary.Operator, binary.Left, binary.Right); absorbed {
		return result, nil
	}

	if result, absorbed := r.applyAbsorption(binary.Operator, binary.Right, binary.Left); absorbed {
		return result, nil
	}

	return node, nil
}

func (r *AbsorptionRule) applyAbsorption(operator lexer.BooleanTokenType, left, right ast.ASTNode) (ast.ASTNode, bool) {
	flippedOperator := flipOperator(operator)

	grouping, ok := right.(*ast.GroupingNode)
	if !ok {
		return ast.NewBinaryNode(operator, left, right), false
	}

	if binary, ok := grouping.Expr.(*ast.BinaryNode); ok {
		if binary.Operator != flippedOperator {
			return ast.NewBinaryNode(operator, left, right), false
		}

		if left.Equals(binary.Left) || right.Equals(binary.Right) {
			return left, true
		}

		return ast.NewBinaryNode(operator, left, right), false
	}

	if chain, ok := grouping.Expr.(*ast.ChainNode); ok {
		if chain.Operator != flippedOperator {
			return ast.NewBinaryNode(operator, left, right), false
		}

		if chain.Contains(left) {
			return left, true
		}
	}

	return ast.NewBinaryNode(operator, left, right), true
}
