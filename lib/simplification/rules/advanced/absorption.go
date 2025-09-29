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
		BaseRule: *base.NewBaseRule("Закон поглощения"),
	}
}

func (r *AbsorptionRule) CanApply(node ast.Node) bool {
	binary, ok := node.(*ast.BinaryExpr)
	if !ok {
		return false
	}

	if binary.Operator != lexer.BlConjunction && binary.Operator != lexer.BlDisjunction {
		return false
	}

	_, lok := binary.Left.(*ast.GroupingExpr)
	_, rok := binary.Right.(*ast.GroupingExpr)

	return lok != rok
}

func (r *AbsorptionRule) Apply(node ast.Node) (ast.Node, error) {
	binary := node.(*ast.BinaryExpr)

	if result, absorbed := r.applyAbsorption(binary.Operator, binary.Left, binary.Right); absorbed {
		return result, nil
	}

	if result, absorbed := r.applyAbsorption(binary.Operator, binary.Right, binary.Left); absorbed {
		return result, nil
	}

	return node, nil
}

func (r *AbsorptionRule) applyAbsorption(operator lexer.TokenType, left, right ast.Expr) (ast.Node, bool) {
	grouping, ok := right.(*ast.GroupingExpr)
	if !ok {
		return ast.NewBinaryExpr(operator, left, right), false
	}

	if binary, ok := grouping.Body.(*ast.BinaryExpr); ok {
		return r.applyBinaryAbsorption(operator, left, right, binary)
	}

	if chain, ok := grouping.Body.(*ast.ChainExpr); ok {
		return r.applyChainAbsorption(operator, left, right, chain)
	}

	return ast.NewBinaryExpr(operator, left, right), true
}

func (r *AbsorptionRule) applyBinaryAbsorption(operator lexer.TokenType, left ast.Expr, right ast.Expr, binary *ast.BinaryExpr) (ast.Node, bool) {
	flippedOperator := flipOperator(operator)
	if binary.Operator != flippedOperator {
		return ast.NewBinaryExpr(operator, left, right), false
	}

	if left.Equals(binary.Left) || right.Equals(binary.Right) {
		return left, true
	}

	if ast.IsNegationOf(binary.Left, left) {
		return ast.NewBinaryExpr(operator, left, ast.NewGroupingExpr(binary.Right)), true
	}
	if ast.IsNegationOf(binary.Right, left) {
		return ast.NewBinaryExpr(operator, left, ast.NewGroupingExpr(binary.Left)), true
	}
	if ast.IsNegationOf(left, binary.Left) {
		return ast.NewBinaryExpr(operator, left, binary.Right), true
	}
	if ast.IsNegationOf(left, binary.Right) {
		return ast.NewBinaryExpr(operator, left, binary.Left), true
	}

	return ast.NewBinaryExpr(operator, left, right), false
}

func (r *AbsorptionRule) applyChainAbsorption(operator lexer.TokenType, left ast.Expr, right ast.Expr, chain *ast.ChainExpr) (ast.Node, bool) {
	flippedOperator := flipOperator(operator)
	if chain.Operator != flippedOperator {
		return ast.NewBinaryExpr(operator, left, right), false
	}

	if chain.Contains(left) {
		return left, true
	}

	if neg := ast.NewUnaryExpr(lexer.BlNegation, left); chain.Contains(neg) {
		chain.Remove(neg)
		return ast.NewBinaryExpr(operator, left, ast.NewGroupingExpr(chain)), true
	}

	if neg, ok := left.(*ast.UnaryExpr); ok && neg.Operator == lexer.BlNegation {
		if chain.Contains(neg.Operand) {
			chain.Remove(neg.Operand)
			return ast.NewBinaryExpr(operator, left, ast.NewGroupingExpr(chain)), true
		}
	}

	return nil, false
}
