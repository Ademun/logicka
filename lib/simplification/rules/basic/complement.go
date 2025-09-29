package basic

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"logicka/lib/simplification/rules/base"
)

type ComplementRule struct {
	base.BaseRule
}

func NewComplementRule() *ComplementRule {
	return &ComplementRule{
		BaseRule: *base.NewBaseRule("Закон дополнения"),
	}
}

func (r *ComplementRule) CanApply(node ast.Node) bool {
	binary, ok := node.(*ast.BinaryExpr)
	if !ok {
		return false
	}

	return binary.Operator == lexer.BlConjunction || binary.Operator == lexer.BlDisjunction
}

func (r *ComplementRule) Apply(node ast.Node) (ast.Node, error) {
	binary := node.(*ast.BinaryExpr)

	switch binary.Operator {
	case lexer.BlConjunction:
		return r.applyConjunctionComplement(binary)
	case lexer.BlDisjunction:
		return r.applyDisjunctionComplement(binary)
	default:
		return node, nil
	}
}

func (r *ComplementRule) applyConjunctionComplement(node *ast.BinaryExpr) (ast.Node, error) {
	if ast.IsNegationOf(node.Left, node.Right) || ast.IsNegationOf(node.Right, node.Left) {
		return ast.NewBooleanExpr(false), nil
	}
	return node, nil
}

func (r *ComplementRule) applyDisjunctionComplement(node *ast.BinaryExpr) (ast.Node, error) {
	if ast.IsNegationOf(node.Left, node.Right) || ast.IsNegationOf(node.Right, node.Left) {
		return ast.NewBooleanExpr(true), nil
	}
	return node, nil
}
