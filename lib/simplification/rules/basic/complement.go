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
		BaseRule: *base.NewBaseRule("Complement law"),
	}
}

func (r *ComplementRule) CanApply(node ast.ASTNode) bool {
	binary, ok := node.(*ast.BinaryNode)
	if !ok {
		return false
	}

	return binary.Operator == lexer.CONJ || binary.Operator == lexer.DISJ
}

func (r *ComplementRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	binary := node.(*ast.BinaryNode)

	switch binary.Operator {
	case lexer.CONJ:
		return r.applyConjunctionComplement(binary)
	case lexer.DISJ:
		return r.applyDisjunctionComplement(binary)
	default:
		return node, nil
	}
}

func (r *ComplementRule) applyConjunctionComplement(node *ast.BinaryNode) (ast.ASTNode, error) {
	if ast.IsNegationOf(node.Left, node.Right) || ast.IsNegationOf(node.Right, node.Left) {
		return ast.NewLiteralNode(false), nil
	}
	return node, nil
}

func (r *ComplementRule) applyDisjunctionComplement(node *ast.BinaryNode) (ast.ASTNode, error) {
	if ast.IsNegationOf(node.Left, node.Right) || ast.IsNegationOf(node.Right, node.Left) {
		return ast.NewLiteralNode(true), nil
	}
	return node, nil
}
