package basic

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"logicka/lib/simplification/rules/base"
)

type DominationRule struct {
	base.BaseRule
}

func NewDominationRule() *DominationRule {
	return &DominationRule{
		BaseRule: *base.NewBaseRule("Закон исключённого третьего"),
	}
}

func (r *DominationRule) CanApply(node ast.Node) bool {
	binary, ok := node.(*ast.BinaryExpr)
	if !ok {
		return false
	}

	return binary.Operator == lexer.BlConjunction || binary.Operator == lexer.BlDisjunction
}

func (r *DominationRule) Apply(node ast.Node) (ast.Node, error) {
	binary := node.(*ast.BinaryExpr)

	switch binary.Operator {
	case lexer.BlConjunction:
		return r.applyConjunctionDomination(binary)
	case lexer.BlDisjunction:
		return r.applyDisjunctionDomination(binary)
	default:
		return node, nil
	}
}

func (r *DominationRule) applyConjunctionDomination(node *ast.BinaryExpr) (ast.Node, error) {
	if ast.IsFalse(node.Left) || ast.IsFalse(node.Right) {
		return ast.NewBooleanExpr(false), nil
	}
	return node, nil
}

func (r *DominationRule) applyDisjunctionDomination(node *ast.BinaryExpr) (ast.Node, error) {
	if ast.IsTrue(node.Left) || ast.IsTrue(node.Right) {
		return ast.NewBooleanExpr(true), nil
	}
	return node, nil
}
