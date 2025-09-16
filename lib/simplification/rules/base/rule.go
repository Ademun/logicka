package base

import (
	"fmt"
	"logicka/lib/ast"
	"strings"
)

type Rule interface {
	CanApply(node ast.ASTNode) bool
	Apply(node ast.ASTNode) (ast.ASTNode, error)
	Name() string
	Applications() []RuleApplication
	RecordApplication(description, before, after string)
	ClearApplications()
}

type RuleApplication struct {
	Name        string
	Description string
	Before      string
	After       string
}

func (ra RuleApplication) String() string {
	return fmt.Sprintf("Правило: %s\nОписание: %s\nДо:\n%s\nПосле:\n%s\n",
		ra.Name,
		ra.Description,
		indent(ra.Before),
		indent(ra.After),
	)
}

func indent(s string) string {
	if s == "" {
		return "  <empty>"
	}
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = "  " + line
	}
	return strings.Join(lines, "\n")
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

func (r *BaseRule) RecordApplication(description, before, after string) {
	r.applications = append(r.applications, RuleApplication{
		Name:        r.name,
		Description: description,
		Before:      before,
		After:       after,
	})
}

func (r *BaseRule) ClearApplications() {
	r.applications = make([]RuleApplication, 0)
}
