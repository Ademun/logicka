package base

import (
	"logicka/lib/ast"
)

type RuleSet struct {
	Name  string
	Rules []Rule
}

func NewRuleSet(name string, rules []Rule) *RuleSet {
	return &RuleSet{name, rules}
}

func (rs *RuleSet) Apply(node ast.ASTNode, log ApplicationLogger) (ast.ASTNode, error) {
	current := node
	appliedRules := make([]Rule, 0)
	for _, rule := range rs.Rules {
		if !rule.CanApply(current) {
			continue
		}

		simplified, err := rule.Apply(current)
		if err != nil {
			return nil, err
		}

		if simplified.Equals(current) {
			continue
		}

		if rule.Name() != "Объединение в цепочку операторов" {
			log.LogApplication(rule.Name(), "Описание", current.String(), simplified.String())
		}
		current = simplified
		appliedRules = append(appliedRules, rule)
	}

	return current, nil
}
