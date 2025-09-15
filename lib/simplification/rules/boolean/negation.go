package boolean

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"logicka/lib/simplification/rules/base"
)

type NegationRule struct {
	base.BaseRule
}

func NewNegationRule() *NegationRule {
	return &NegationRule{
		BaseRule: *base.NewBaseRule("Negation law"),
	}
}

func (r *NegationRule) CanApply(node ast.ASTNode) bool {
	binary, ok := node.(*ast.BinaryNode)
	if !ok {
		return false
	}

	return binary.Operator == lexer.CONJ || binary.Operator == lexer.DISJ
}

func (r *NegationRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	binary := node.(*ast.BinaryNode)

	switch binary.Operator {
	case lexer.CONJ:
		return r.applyConjunctionNegation(binary)
	case lexer.DISJ:
		return r.applyDisjunctionNegation(binary)
	default:
		return node, nil
	}
}

func (r *NegationRule) applyConjunctionNegation(node *ast.BinaryNode) (ast.ASTNode, error) {
	if ast.IsNegationOf(node.Left, node.Right) || ast.IsNegationOf(node.Right, node.Left) {
		return ast.NewLiteralNode(false), nil
	}
	return node, nil
}

func (r *NegationRule) applyDisjunctionNegation(node *ast.BinaryNode) (ast.ASTNode, error) {
	if ast.IsNegationOf(node.Left, node.Right) || ast.IsNegationOf(node.Right, node.Left) {
		return ast.NewLiteralNode(true), nil
	}
	return node, nil
}
