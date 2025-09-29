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

func (r *LiteralNegationRule) CanApply(node ast.Node) bool {
	unary, ok := node.(*ast.UnaryExpr)
	if !ok {
		return false
	}

	_, ok = unary.Operand.(*ast.BooleanExpr)

	return unary.Operator == lexer.BlNegation && ok
}

func (r *LiteralNegationRule) Apply(node ast.Node) (ast.Node, error) {
	unary := node.(*ast.UnaryExpr)

	literal := unary.Operand.(*ast.BooleanExpr)

	return ast.NewBooleanExpr(!literal.Value), nil
}
