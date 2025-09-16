package complex

import "logicka/lib/simplification/rules/base"

func CreateComplexRules() []base.Rule {
	return []base.Rule{
		NewAbsorptionRule(),
		NewDeMorganRule(),
	}
}
