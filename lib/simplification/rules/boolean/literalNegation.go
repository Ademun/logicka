package boolean

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
		BaseRule: *base.NewBaseRule("LiteralNegation law"),
	}
}

func (r *LiteralNegationRule) CanApply(node ast.ASTNode) bool {
	unary, ok := node.(*ast.UnaryNode)
	if !ok {
		return false
	}

	return unary.Operator == lexer.NEG
}

func (r *LiteralNegationRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	unary := node.(*ast.UnaryNode)

	if literalOp, ok := unary.Operand.(*ast.LiteralNode); ok {
		return ast.NewLiteralNode(!literalOp.Value), nil
	}

	return node, nil
}
