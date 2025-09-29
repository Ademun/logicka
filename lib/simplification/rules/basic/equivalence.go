package basic

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"logicka/lib/simplification/rules/base"
)

type EquivalenceRule struct {
	base.BaseRule
}

func NewEquivalenceRule() *EquivalenceRule {
	return &EquivalenceRule{
		BaseRule: *base.NewBaseRule("Упрощение эквивалентности"),
	}
}

func (r *EquivalenceRule) CanApply(node ast.Node) bool {
	binary, ok := node.(*ast.BinaryExpr)
	if !ok {
		return false
	}

	return binary.Operator == lexer.BlEquivalence
}

func (r *EquivalenceRule) Apply(node ast.Node) (ast.Node, error) {
	binary := node.(*ast.BinaryExpr)

	return ast.NewBinaryExpr(
		lexer.BlConjunction,
		ast.NewBinaryExpr(lexer.BlImplication, binary.Left, binary.Right),
		ast.NewBinaryExpr(lexer.BlImplication, binary.Right, binary.Left),
	), nil
}
