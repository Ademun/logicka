package base

import "logicka/lib/ast"

type Rule interface {
	CanApply(node ast.ASTNode) bool
	Apply(node ast.ASTNode) (ast.ASTNode, error)
	Name() string
}

type BaseRule struct {
	name string
}

func (r *BaseRule) Name() string {
	return r.name
}
