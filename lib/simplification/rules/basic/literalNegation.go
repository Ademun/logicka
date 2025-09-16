package basic

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"logicka/lib/simplification/rules/base"
)

type LiteralNegationRule struct {
	base.BaseRule
}

func NewLiteralNegationRule() *LiteralNegationRule {
	return &LiteralNegationRule{
		BaseRule: *base.NewBaseRule("Отрицание константы"),
	}
}

func (r *LiteralNegationRule) CanApply(node ast.ASTNode) bool {
	unary, ok := node.(*ast.UnaryNode)
	if !ok {
		return false
	}

	_, ok = unary.Operand.(*ast.LiteralNode)

	return unary.Operator == lexer.NEG && ok
}

func (r *LiteralNegationRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	unary := node.(*ast.UnaryNode)

	literal := unary.Operand.(*ast.LiteralNode)

	return ast.NewLiteralNode(!literal.Value), nil
}
