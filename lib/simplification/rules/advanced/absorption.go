package advanced

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
	grouping, ok := right.(*ast.GroupingNode)
	if !ok {
		return ast.NewBinaryNode(operator, left, right), false
	}

	if binary, ok := grouping.Expr.(*ast.BinaryNode); ok {
		return r.applyBinaryAbsorption(operator, left, right, binary)
	}

	if chain, ok := grouping.Expr.(*ast.ChainNode); ok {
		return r.applyChainAbsorption(operator, left, right, chain)
	}

	return ast.NewBinaryNode(operator, left, right), true
}

func (r *AbsorptionRule) applyBinaryAbsorption(operator lexer.BooleanTokenType, left ast.ASTNode, right ast.ASTNode, binary *ast.BinaryNode) (ast.ASTNode, bool) {
	flippedOperator := flipOperator(operator)
	if binary.Operator != flippedOperator {
		return ast.NewBinaryNode(operator, left, right), false
	}

	if left.Equals(binary.Left) || right.Equals(binary.Right) {
		return left, true
	}

	if ast.IsNegationOf(binary.Left, left) {
		return ast.NewBinaryNode(operator, left, ast.NewGroupingNode(binary.Right)), true
	}
	if ast.IsNegationOf(binary.Right, left) {
		return ast.NewBinaryNode(operator, left, ast.NewGroupingNode(binary.Left)), true
	}
	if ast.IsNegationOf(left, binary.Left) {
		return ast.NewBinaryNode(operator, left, binary.Right), true
	}
	if ast.IsNegationOf(left, binary.Right) {
		return ast.NewBinaryNode(operator, left, binary.Left), true
	}

	return ast.NewBinaryNode(operator, left, right), false
}

func (r *AbsorptionRule) applyChainAbsorption(operator lexer.BooleanTokenType, left ast.ASTNode, right ast.ASTNode, chain *ast.ChainNode) (ast.ASTNode, bool) {
	flippedOperator := flipOperator(operator)
	if chain.Operator != flippedOperator {
		return ast.NewBinaryNode(operator, left, right), false
	}

	if chain.Contains(left) {
		return left, true
	}

	if neg := ast.NewUnaryNode(lexer.NEG, left); chain.Contains(neg) {
		chain.Remove(neg)
		return ast.NewBinaryNode(operator, left, ast.NewGroupingNode(chain)), true
	}

	if neg, ok := left.(*ast.UnaryNode); ok && neg.Operator == lexer.NEG {
		if chain.Contains(neg.Operand) {
			chain.Remove(neg.Operand)
			return ast.NewBinaryNode(operator, left, ast.NewGroupingNode(chain)), true
		}
	}

	return nil, false
}
