package base

import (
	"fmt"
	"logicka/lib/ast"
)

type Rule interface {
	CanApply(node ast.ASTNode) bool
	Apply(node ast.ASTNode) (ast.ASTNode, error)
	Name() string
}

type RuleApplication struct {
	Order       int
	Name        string
	Description string
	Before      string
	After       string
}

func (ra RuleApplication) String() string {
	return fmt.Sprintf(
		"%d. %s\n%s => %s\n",
		ra.Order,
		ra.Name,
		ra.Before,
		ra.After,
	)
}

type BaseRule struct {
	name string
}

func NewBaseRule(name string) *BaseRule {
	return &BaseRule{
		name: name,
	}
}

func (r *BaseRule) Name() string {
	return r.name
}
