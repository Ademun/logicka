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

func (r *DoubleNegationRule) CanApply(node ast.ASTNode) bool {
	unary, ok := node.(*ast.UnaryNode)
	if !ok {
		return false
	}

	return unary.Operator == lexer.NEG
}

func (r *DoubleNegationRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	unary := node.(*ast.UnaryNode)

	if negOp, ok := unary.Operand.(*ast.UnaryNode); ok && negOp.Operator == lexer.NEG {
		return negOp.Operand, nil
	}

	return node, nil
}
