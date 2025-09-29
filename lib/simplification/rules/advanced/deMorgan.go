package advanced

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"logicka/lib/simplification/rules/base"
)

type DeMorganRule struct {
	base.BaseRule
}

func NewDeMorganRule() *DeMorganRule {
	return &DeMorganRule{
		BaseRule: *base.NewBaseRule("Закон Де Моргана"),
	}
}

func (r *DeMorganRule) CanApply(node ast.Node) bool {
	unary, ok := node.(*ast.UnaryExpr)
	if !ok {
		return false
	}

	_, ok = unary.Operand.(*ast.GroupingExpr)

	return unary.Operator == lexer.BlNegation && ok
}

func (r *DeMorganRule) Apply(node ast.Node) (ast.Node, error) {
	unary := node.(*ast.UnaryExpr)
	grouping := unary.Operand.(*ast.GroupingExpr)

	switch v := grouping.Body.(type) {
	case *ast.BinaryExpr:
		if v.Operator != lexer.BlConjunction && v.Operator != lexer.BlDisjunction {
			return node, nil
		}
		return r.applyDeMorganBinary(v)
	case *ast.ChainExpr:
		if v.Operator != lexer.BlConjunction && v.Operator != lexer.BlDisjunction {
			return node, nil
		}
		return r.applyDeMorganChain(v)
	default:
		return node, nil
	}
}

func (r *DeMorganRule) applyDeMorganBinary(node *ast.BinaryExpr) (ast.Node, error) {
	newOperator := flipOperator(node.Operator)

	return ast.NewGroupingExpr(
		ast.NewBinaryExpr(
			newOperator,
			ast.NewUnaryExpr(lexer.BlNegation, node.Left),
			ast.NewUnaryExpr(lexer.BlNegation, node.Right),
		),
	), nil
}

func (r *DeMorganRule) applyDeMorganChain(node *ast.ChainExpr) (ast.Node, error) {
	newOperator := flipOperator(node.Operator)

	negatedOps := make([]ast.Expr, len(node.Elements))
	for i, op := range node.Elements {
		negatedOps[i] = ast.NewUnaryExpr(lexer.BlNegation, op)
	}

	return ast.NewGroupingExpr(
		&ast.ChainExpr{
			Operator: newOperator,
			Elements: negatedOps,
		},
	), nil
}

func flipOperator(operator lexer.TokenType) lexer.TokenType {
	var newOperator lexer.TokenType

	switch {
	case operator == lexer.BlConjunction:
		newOperator = lexer.BlDisjunction
	case operator == lexer.BlDisjunction:
		newOperator = lexer.BlConjunction
	default:
		return operator
	}

	return newOperator
}
