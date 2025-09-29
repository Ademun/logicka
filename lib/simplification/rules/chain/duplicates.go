package chain

import (
	"logicka/lib/ast"
	"logicka/lib/simplification/rules/base"
)

type DuplicatesRule struct {
	base.BaseRule
}

func NewDuplicatesRule() *DuplicatesRule {
	return &DuplicatesRule{
		BaseRule: *base.NewBaseRule("Сокращение дубликатов"),
	}
}

func (r *DuplicatesRule) CanApply(node ast.Node) bool {
	_, ok := node.(*ast.ChainExpr)

	return ok
}

func (r *DuplicatesRule) Apply(node ast.Node) (ast.Node, error) {
	chain := node.(*ast.ChainExpr)
	operands := collectUniqueOperands(chain.Elements)

	switch len(operands) {
	case 1:
		return operands[0], nil
	case 2:
		return ast.NewBinaryExpr(chain.Operator, operands[0], operands[1]), nil
	default:
		return ast.NewChainExpr(chain.Operator, operands), nil
	}
}

func collectUniqueOperands(operands []ast.Expr) []ast.Expr {
	unique := make(map[uint64]ast.Expr, len(operands)/2)

	for _, operand := range operands {
		hash := operand.Hash()
		if _, ok := unique[hash]; ok {
			continue
		}
		unique[hash] = operand
	}

	result := make([]ast.Expr, 0, len(unique))
	for _, operand := range unique {
		result = append(result, operand)
	}

	return result
}
