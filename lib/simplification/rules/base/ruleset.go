package base

import (
	"logicka/lib/ast"
	"slices"
	"strings"
)

type RuleSet struct {
	Name  string
	Rules []Rule
}

func NewRuleSet(name string, rules []Rule) *RuleSet {
	return &RuleSet{name, rules}
}

func (rs *RuleSet) Apply(node ast.ASTNode) (ast.ASTNode, error) {
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
			rule.RecordApplication("Описание", current.String(), simplified.String())
		}
		current = simplified
		appliedRules = append(appliedRules, rule)
	}

	return current, nil
}

func (rs *RuleSet) String(verbose bool, reset bool) string {
	result := strings.Builder{}
	for _, rule := range rs.Rules {
		records := slices.Clone(rule.Applications())
		if len(records) == 0 {
			continue
		}
		for _, record := range records {
			if verbose {
				result.WriteString(record.VerboseString())
			} else {
				result.WriteString(record.String())
			}
		}
		if reset {
			rule.ClearApplications()
		}
	}
	return result.String()
}
