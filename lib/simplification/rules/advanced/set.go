package advanced

import "logicka/lib/simplification/rules/base"

func CreateAdvancedRules() []base.Rule {
	return []base.Rule{
		NewAbsorptionRule(),
		NewDeMorganRule(),
		NewResolutionRule(),
	}
}

func CreateAdvancedRuleSet() *base.RuleSet {
	ruleSet := &base.RuleSet{
		Rules: CreateAdvancedRules(),
	}
	return ruleSet
}
