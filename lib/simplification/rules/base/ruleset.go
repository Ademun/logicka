package base

import "logicka/lib/ast"

type RuleSet struct {
	Name  string
	Rules []Rule
}

func NewRuleSet(name string, rules []Rule) *RuleSet {
	return &RuleSet{name, rules}
}

func (rs *RuleSet) Apply(node ast.ASTNode) (ast.ASTNode, error) {
	current := node

	for _, rule := range rs.Rules {
		if !rule.CanApply(node) {
			continue
		}

		simplified, err := rule.Apply(current)
		if err != nil {
			return nil, err
		}

		current = simplified
	}

	return current, nil
}
