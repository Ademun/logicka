package base

import "logicka/lib/ast"

type Rule interface {
	CanApply(node ast.ASTNode) bool
	Apply(node ast.ASTNode) (ast.ASTNode, error)
	Name() string
	Applications() []RuleApplication
}

type RuleApplication struct {
	Name        string
	Description string
	Before      ast.ASTNode
	After       ast.ASTNode
}

type BaseRule struct {
	name         string
	applications []RuleApplication
}

func NewBaseRule(name string) *BaseRule {
	return &BaseRule{
		name:         name,
		applications: make([]RuleApplication, 0),
	}
}

func (r *BaseRule) Name() string {
	return r.name
}

func (r *BaseRule) Applications() []RuleApplication {
	return r.applications
}

func (r *BaseRule) RecordApplication(description string, before ast.ASTNode, after ast.ASTNode) {
	r.applications = append(r.applications, RuleApplication{
		Name:        r.name,
		Description: description,
		Before:      before,
		After:       after,
	})
}
