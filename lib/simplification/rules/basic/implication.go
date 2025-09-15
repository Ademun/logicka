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
		BaseRule: *base.NewBaseRule("Implication law"),
	}
}

func (r *ImplicationRule) CanApply(node ast.ASTNode) bool {
	binary, ok := node.(*ast.BinaryNode)
	if !ok {
		return false
	}

	return binary.Operator == lexer.IMPL
}

func (r *ImplicationRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	binary := node.(*ast.BinaryNode)

	return ast.NewBinaryNode(lexer.DISJ, ast.NewUnaryNode(lexer.NEG, binary.Left), binary.Right), nil
}
