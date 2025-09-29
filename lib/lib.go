package lib

import (
	"fmt"
	"logicka/lib/lexer"
	"logicka/lib/parser"
	"logicka/lib/simplification/rules/advanced"
	"logicka/lib/simplification/rules/basic"
	"logicka/lib/simplification/rules/chain"
	"logicka/lib/visitor"

	"github.com/sanity-io/litter"
)

type Logicka struct {
}

/*
R:= {"h", "k", "d", "j", "a", "i", "b", "g", "c", "f", "e"};
N:={"d", "g", "a", "b", "i", "j"};
L:={"c", "b", "g", "e", "f", "i"};
Q(x) := x element_of N; A(x) := x element_of L;
P1(x) := Q(x) conjunction A(x);
P2(x):=Q(x) equivalence A(x);
P3(x):=Q(x) implication A(x);
P4(x) := Q(x) disjunction A(x);
*/
func (l *Logicka) CalculateTruthTable(expr string, values map[string]bool) (interface{}, error) {
	lex := lexer.NewLexer()
	tokens, err := lex.Tokenize(expr)
	if err != nil {
		return nil, err
	}

	p := parser.NewParser()
	ast, err := p.Parse(tokens)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	litter.Dump(ast)
	fmt.Println(ast.String())

	simplifier := visitor.NewSimplifier()
	simplifier.AddRuleSet(basic.CreateBasicRuleSet())
	simplifier.AddRuleSet(advanced.CreateAdvancedRuleSet())
	simplifier.AddRuleSet(chain.CreateChainRuleSet())
	simplified, err := simplifier.Simplify(ast)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("simplification error: %w", err)
	}

	fmt.Println(simplified.String())

	return nil, nil
}
