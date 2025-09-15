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
		BaseRule: *base.NewBaseRule("Domination law"),
	}
}

func (r *DominationRule) CanApply(node ast.ASTNode) bool {
	binary, ok := node.(*ast.BinaryNode)
	if !ok {
		return false
	}

	return binary.Operator == lexer.CONJ || binary.Operator == lexer.DISJ
}

func (r *DominationRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	binary := node.(*ast.BinaryNode)

	switch binary.Operator {
	case lexer.CONJ:
		return r.applyConjunctionDomination(binary)
	case lexer.DISJ:
		return r.applyDisjunctionDomination(binary)
	default:
		return node, nil
	}
}

func (r *DominationRule) applyConjunctionDomination(node *ast.BinaryNode) (ast.ASTNode, error) {
	if ast.IsFalse(node.Left) || ast.IsFalse(node.Right) {
		return ast.NewLiteralNode(false), nil
	}
	return node, nil
}

func (r *DominationRule) applyDisjunctionDomination(node *ast.BinaryNode) (ast.ASTNode, error) {
	if ast.IsTrue(node.Left) || ast.IsTrue(node.Right) {
		return ast.NewLiteralNode(true), nil
	}
	return node, nil
}
