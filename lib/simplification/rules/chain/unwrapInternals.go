package chain

import (
	"fmt"
	"logicka/lib/ast"
	"logicka/lib/simplification/rules/base"
	"slices"
)

type UnwrapInternalsRule struct {
	base.BaseRule
}

func NewUnwrapInternalsRule() *UnwrapInternalsRule {
	return &UnwrapInternalsRule{
		BaseRule: *base.NewBaseRule("Развёртка внутренних операторов"),
	}
}

func (r *UnwrapInternalsRule) CanApply(node ast.ASTNode) bool {
	_, ok := node.(*ast.ChainNode)
	return ok
}

func (r *UnwrapInternalsRule) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	chain := node.(*ast.ChainNode)
	for _, op := range slices.Clone(chain.Children()) {
		if opChain, ok := op.(*ast.ChainNode); ok && opChain.IsType(chain.Operator) {
			chain.Remove(opChain)
			chain.Add(opChain.Children()...)
		}
		if opBinary, ok := op.(*ast.BinaryNode); ok && opBinary.IsType(chain.Operator) {
			chain.Remove(opBinary)
			chain.Add(opBinary.Children()...)
		}
	}
	fmt.Println("exited", chain.Operands)
	return chain, nil
}
