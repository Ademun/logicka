package basic

import "logicka/lib/simplification/rules/base"

func CreateBasicRules() []base.Rule {
	return []base.Rule{
		NewIdentityRule(),
		NewDominationRule(),
		NewIdempotencyRule(),
		NewDoubleNegationRule(),
		NewComplementRule(),
		NewLiteralNegationRule(),
		NewImplicationRule(),
		NewEquivalenceRule(),
	}
}

func CreateBasicRuleSet() *base.RuleSet {
	ruleSet := &base.RuleSet{
		Rules: CreateBasicRules(),
	}
	return ruleSet
}
