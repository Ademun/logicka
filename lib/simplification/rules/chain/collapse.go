package chain

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"logicka/lib/simplification/rules/base"
)

type CollapseRule struct {
	base.BaseRule
}

func NewCollapseRule() *CollapseRule {
	return &CollapseRule{
		BaseRule: *base.NewBaseRule("Закон тождественности"),
	}
}

func (r *CollapseRule) CanApply(node ast.ASTNode) bool {
	_, ok := node.(*ast.ChainNode)

	return ok
}

func (r *CollapseRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	chain := node.(*ast.ChainNode)
	if chain.Contains(ast.NewLiteralNode(true)) && chain.IsType(lexer.DISJ) {
		return ast.NewLiteralNode(true), nil
	}
	if chain.Contains(ast.NewLiteralNode(false)) && chain.IsType(lexer.CONJ) {
		return ast.NewLiteralNode(false), nil
	}
	return node, nil
}
