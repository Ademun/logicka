package basic

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"logicka/lib/simplification/rules/base"
)

type ImplicationRule struct {
	base.BaseRule
}

func NewImplicationRule() *ImplicationRule {
	return &ImplicationRule{
		BaseRule: *base.NewBaseRule("Упрощение импликации"),
	}
}

func (r *ImplicationRule) CanApply(node ast.Node) bool {
	binary, ok := node.(*ast.BinaryExpr)
	if !ok {
		return false
	}

	return binary.Operator == lexer.BlImplication
}

func (r *ImplicationRule) Apply(node ast.Node) (ast.Node, error) {
	binary := node.(*ast.BinaryExpr)

	return ast.NewBinaryExpr(lexer.BlDisjunction, ast.NewUnaryExpr(lexer.BlNegation, binary.Left), binary.Right), nil
}
