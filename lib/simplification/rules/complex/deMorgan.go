package complex

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
		BaseRule: *base.NewBaseRule("DeMorgan law"),
	}
}

func (r *DeMorganRule) CanApply(node ast.ASTNode) bool {
	unary, ok := node.(*ast.UnaryNode)
	if !ok || unary.Operator != lexer.NEG {
		return false
	}

	_, ok = unary.Operand.(*ast.GroupingNode)

	return ok
}

func (r *DeMorganRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	unary := node.(*ast.UnaryNode)
	grouping := unary.Operand.(*ast.GroupingNode)

	switch v := grouping.Expr.(type) {
	case *ast.BinaryNode:
		if v.Operator != lexer.CONJ && v.Operator != lexer.DISJ {
			return node, nil
		}
		return r.applyDeMorganBinary(v)
	case *ast.ChainNode:
		if v.Operator != lexer.CONJ && v.Operator != lexer.DISJ {
			return node, nil
		}
		return r.applyDeMorganChain(v)
	default:
		return node, nil
	}
}

func (r *DeMorganRule) applyDeMorganBinary(node *ast.BinaryNode) (ast.ASTNode, error) {
	newOperator := flipOperator(node.Operator)

	return ast.NewGroupingNode(
		ast.NewBinaryNode(
			newOperator,
			ast.NewUnaryNode(lexer.NEG, node.Left),
			ast.NewUnaryNode(lexer.NEG, node.Right),
		),
	), nil
}

func (r *DeMorganRule) applyDeMorganChain(node *ast.ChainNode) (ast.ASTNode, error) {
	newOperator := flipOperator(node.Operator)

	negatedOps := make([]ast.ASTNode, len(node.Operands))
	for i, op := range node.Operands {
		negatedOps[i] = ast.NewUnaryNode(lexer.NEG, op)
	}

	return ast.NewGroupingNode(
		&ast.ChainNode{
			Operator: newOperator,
			Operands: negatedOps,
		},
	), nil
}

func flipOperator(operator lexer.BooleanTokenType) lexer.BooleanTokenType {
	var newOperator lexer.BooleanTokenType

	switch {
	case operator == lexer.CONJ:
		newOperator = lexer.DISJ
	case operator == lexer.DISJ:
		newOperator = lexer.CONJ
	default:
		return operator
	}

	return newOperator
}
