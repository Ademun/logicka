package chain

import (
	"logicka/lib/ast"
	"logicka/lib/simplification/rules/base"
	"slices"
)

type DuplicatesRule struct {
	base.BaseRule
}

func NewDuplicatesRule() *DuplicatesRule {
	return &DuplicatesRule{
		BaseRule: *base.NewBaseRule("Duplicates law"),
	}
}

func (r *DuplicatesRule) CanApply(node ast.ASTNode) bool {
	_, ok := node.(*ast.ChainNode)

	return ok
}

func (r *DuplicatesRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	chain := node.(*ast.ChainNode)
	operands := collectUniqueOperands(chain.Operands)

	switch len(operands) {
	case 1:
		return operands[0], nil
	case 2:
		return ast.NewBinaryNode(chain.Operator, operands[0], operands[1]), nil
	default:
		return ast.NewChainNode(chain.Operator, operands...)
	}
}

func collectUniqueOperands(operands []ast.ASTNode) []ast.ASTNode {
	return slices.CompactFunc(operands, func(a, b ast.ASTNode) bool {
		return a.Equals(b)
	})
}
