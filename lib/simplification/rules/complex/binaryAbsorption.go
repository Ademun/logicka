package complex

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"logicka/lib/simplification/rules/base"
)

type BinaryAbsorptionRule struct {
	base.BaseRule
}

func NewAbsorptionRule() *BinaryAbsorptionRule {
	return &BinaryAbsorptionRule{
		BaseRule: *base.NewBaseRule("Absorption law"),
	}
}

func (r *BinaryAbsorptionRule) CanApply(node ast.ASTNode) bool {
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

func (r *BinaryAbsorptionRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	binary := node.(*ast.BinaryNode)

	if result, absorbed := r.applyAbsorption(binary.Operator, binary.Left, binary.Right); absorbed {
		return result, nil
	}

	if result, absorbed := r.applyAbsorption(binary.Operator, binary.Right, binary.Left); absorbed {
		return result, nil
	}

	return node, nil
}

func (r *BinaryAbsorptionRule) applyAbsorption(operator lexer.BooleanTokenType, left, right ast.ASTNode) (ast.ASTNode, bool) {
	grouping, ok := right.(*ast.GroupingNode)
	if !ok {
		return ast.NewBinaryNode(operator, left, right), false
	}

	binary, ok := grouping.Expr.(*ast.BinaryNode)
	if !ok {
		return ast.NewBinaryNode(operator, left, right), false
	}

	flippedOperator := flipOperator(operator)
	if binary.Operator != flippedOperator {
		return ast.NewBinaryNode(operator, left, right), false
	}

	if left.Equals(binary.Left) || right.Equals(binary.Right) {
		return left, true
	}

	return ast.NewBinaryNode(operator, left, right), false
}
