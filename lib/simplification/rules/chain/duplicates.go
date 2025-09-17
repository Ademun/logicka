package chain

import (
	"fmt"
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

func (r *DuplicatesRule) CanApply(node ast.ASTNode) bool {
	_, ok := node.(*ast.ChainNode)

	return ok
}

func (r *DuplicatesRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	chain := node.(*ast.ChainNode)
	fmt.Println("Operands", chain.Operands)
	operands := collectUniqueOperands(chain.Operands)

	switch len(operands) {
	case 1:
		return operands[0], nil
	case 2:
		return ast.NewBinaryNode(chain.Operator, operands[0], operands[1]), nil
	default:
		return ast.NewChainNode(chain.Operator, operands...), nil
	}
}

func collectUniqueOperands(operands []ast.ASTNode) []ast.ASTNode {
	unique := make(map[uint64]ast.ASTNode, len(operands)/2)

	for _, operand := range operands {
		hash := operand.Hash()
		if _, ok := unique[hash]; ok {
			continue
		}
		unique[hash] = operand
	}

	result := make([]ast.ASTNode, 0, len(unique))
	for _, operand := range unique {
		result = append(result, operand)
	}

	return result
}
