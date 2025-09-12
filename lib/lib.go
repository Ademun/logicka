package lib

import (
	"fmt"
	"logicka/lib/lexer"
	"logicka/lib/parser"
	"logicka/lib/visitor"
	"regexp"
	"slices"
	"strings"
)

type Logicka struct {
}

func (l *Logicka) CalculateTruthTable(expr string, values map[string]bool) ([]visitor.TruthTableEntry, error) {
	lex := lexer.NewBooleanLexer(expr)
	tokens, err := lex.Lex()
	if err != nil {
		return nil, err
	}

	p := &parser.Parser{Tokens: tokens}
	ast, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}

	ctx := &visitor.EvaluationContext{Variables: values}
	solver := visitor.NewBooleanSolver(ctx)
	simplifier := visitor.NewSimplifier()
	simplified, err := simplifier.Simplify(ast)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("simplification error: %w", err)
	}

	fmt.Println(simplified.String())

	table, err := solver.Solve(simplified)
	if err != nil {
		return nil, fmt.Errorf("solving error: %w", err)
	}

	for _, entry := range table {
		slices.SortFunc(entry.Variables, sortVariables)
	}

	return table, nil
}

func (l *Logicka) SimplifyExpression(expr string) (string, error) {
	lex := lexer.NewBooleanLexer(expr)
	tokens, err := lex.Lex()
	if err != nil {
		return "", fmt.Errorf("lexing error: %w", err)
	}

	p := &parser.Parser{Tokens: tokens}
	ast, err := p.ParseExpression()
	if err != nil {
		return "", err
	}

	simplifier := visitor.NewSimplifier()
	simplified, err := simplifier.Simplify(ast)
	if err != nil {
		return "", fmt.Errorf("simplification error: %w", err)
	}

	return simplified.String(), nil
}

func sortVariables(a, b visitor.TruthTableVariable) int {
	return strings.Compare(a.Name, b.Name)
}

func (l *Logicka) ExtractVariables(expr string) []string {
	re := regexp.MustCompile(`\b[a-z]+\b`)
	words := re.FindAllString(expr, -1)

	keywords := map[string]struct{}{
		"true": {}, "false": {}, "nil": {},
		"and": {}, "or": {}, "not": {},
		"if": {}, "else": {}, "for": {},
	}

	seen := make(map[string]struct{})
	var variables []string

	for _, word := range words {
		if _, ok := keywords[word]; ok {
			continue
		}
		if _, ok := seen[word]; ok {
			continue
		}
		seen[word] = struct{}{}
		variables = append(variables, word)
	}

	return variables
}
