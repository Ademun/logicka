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
