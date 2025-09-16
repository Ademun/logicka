package basic

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"logicka/lib/simplification/rules/base"
)

type IdempotencyRule struct {
	base.BaseRule
}

func NewIdempotencyRule() *IdempotencyRule {
	return &IdempotencyRule{
		BaseRule: *base.NewBaseRule("Закон идемпотентности"),
	}
}

func (r *IdempotencyRule) CanApply(node ast.ASTNode) bool {
	binary, ok := node.(*ast.BinaryNode)
	if !ok {
		return false
	}

	return binary.Operator == lexer.CONJ || binary.Operator == lexer.DISJ
}

func (r *IdempotencyRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	binary := node.(*ast.BinaryNode)

	if binary.Left.Equals(binary.Right) {
		return binary.Left, nil
	}

	return node, nil
}
