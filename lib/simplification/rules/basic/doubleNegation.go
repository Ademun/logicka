package basic

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"logicka/lib/simplification/rules/base"
)

type DoubleNegationRule struct {
	base.BaseRule
}

func NewDoubleNegationRule() *DoubleNegationRule {
	return &DoubleNegationRule{
		BaseRule: *base.NewBaseRule("Закон двойного отрицания"),
	}
}

func (r *DoubleNegationRule) CanApply(node ast.Node) bool {
	unary, ok := node.(*ast.UnaryExpr)
	if !ok {
		return false
	}

	return unary.Operator == lexer.BlNegation
}

func (r *DoubleNegationRule) Apply(node ast.Node) (ast.Node, error) {
	unary := node.(*ast.UnaryExpr)

	if negOp, ok := unary.Operand.(*ast.UnaryExpr); ok && negOp.Operator == lexer.BlNegation {
		return negOp.Operand, nil
	}

	return node, nil
}
