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

func (r *IdempotencyRule) CanApply(node ast.Node) bool {
	binary, ok := node.(*ast.BinaryExpr)
	if !ok {
		return false
	}

	return binary.Operator == lexer.BlConjunction || binary.Operator == lexer.BlDisjunction
}

func (r *IdempotencyRule) Apply(node ast.Node) (ast.Node, error) {
	binary := node.(*ast.BinaryExpr)

	if binary.Left.Equals(binary.Right) {
		return binary.Left, nil
	}

	return node, nil
}
