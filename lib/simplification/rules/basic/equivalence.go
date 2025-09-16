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

func (r *EquivalenceRule) CanApply(node ast.ASTNode) bool {
	binary, ok := node.(*ast.BinaryNode)
	if !ok {
		return false
	}

	return binary.Operator == lexer.EQUIV
}

func (r *EquivalenceRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	binary := node.(*ast.BinaryNode)

	return ast.NewBinaryNode(
		lexer.CONJ,
		ast.NewBinaryNode(lexer.IMPL, binary.Left, binary.Right),
		ast.NewBinaryNode(lexer.IMPL, binary.Right, binary.Left),
	), nil
}
