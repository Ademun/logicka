package basic

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"logicka/lib/simplification/rules/base"
)

type IdentityRule struct {
	base.BaseRule
}

func NewIdentityRule() *IdentityRule {
	return &IdentityRule{
		BaseRule: *base.NewBaseRule("Закон тождественности"),
	}
}

func (r *IdentityRule) CanApply(node ast.Node) bool {
	binary, ok := node.(*ast.BinaryExpr)
	if !ok {
		return false
	}

	return binary.Operator == lexer.BlConjunction || binary.Operator == lexer.BlDisjunction
}

func (r *IdentityRule) Apply(node ast.Node) (ast.Node, error) {
	binary := node.(*ast.BinaryExpr)

	switch binary.Operator {
	case lexer.BlConjunction:
		return r.applyConjunctionIdentity(binary)
	case lexer.BlDisjunction:
		return r.applyDisjunctionIdentity(binary)
	default:
		return node, nil
	}
}

func (r *IdentityRule) applyConjunctionIdentity(node *ast.BinaryExpr) (ast.Node, error) {
	if ast.IsTrue(node.Left) {
		return node.Right, nil
	}
	if ast.IsTrue(node.Right) {
		return node.Left, nil
	}
	return node, nil
}

func (r *IdentityRule) applyDisjunctionIdentity(node *ast.BinaryExpr) (ast.Node, error) {
	if ast.IsFalse(node.Left) {
		return node.Right, nil
	}
	if ast.IsFalse(node.Right) {
		return node.Left, nil
	}
	return node, nil
}
