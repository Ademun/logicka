package chain

import "logicka/lib/simplification/rules/base"

func CreateChainRules() []base.Rule {
	return []base.Rule{
		NewDuplicatesRule(),
	}
}

func CreateChainRuleSet() *base.RuleSet {
	ruleSet := &base.RuleSet{
		Rules: CreateChainRules(),
	}
	return ruleSet
}
