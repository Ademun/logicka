package boolean

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
		BaseRule: *base.NewBaseRule("Identity law"),
	}
}

func (r *IdentityRule) CanApply(node ast.ASTNode) bool {
	binary, ok := node.(*ast.BinaryNode)
	if !ok {
		return false
	}

	return binary.Operator == lexer.CONJ || binary.Operator == lexer.DISJ
}

func (r *IdentityRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	binary := node.(*ast.BinaryNode)

	switch binary.Operator {
	case lexer.CONJ:
		return r.applyConjunctionIdentity(binary)
	case lexer.DISJ:
		return r.applyDisjunctionIdentity(binary)
	default:
		return node, nil
	}
}

func (r *IdentityRule) applyConjunctionIdentity(node *ast.BinaryNode) (ast.ASTNode, error) {
	if ast.IsTrue(node.Left) {
		return node.Right, nil
	}
	if ast.IsTrue(node.Right) {
		return node.Left, nil
	}
	return node, nil
}

func (r *IdentityRule) applyDisjunctionIdentity(node *ast.BinaryNode) (ast.ASTNode, error) {
	if ast.IsFalse(node.Left) {
		return node.Right, nil
	}
	if ast.IsFalse(node.Right) {
		return node.Left, nil
	}
	return node, nil
}
