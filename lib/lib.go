package lib

import (
	"logicka/lib/lexer"
	"logicka/lib/parser"
	"logicka/lib/visitor"
	"regexp"
	"slices"
	"strings"
)

func GenerateTruthTable(expr string, values map[string]bool) ([]visitor.TruthTableEntry, error) {
	lex := lexer.NewLexer(expr)
	tokens, err := lex.Lex()
	if err != nil {
		return nil, err
	}
	prsr := &parser.Parser{Tokens: tokens}
	ast, err := prsr.ParseExpression()
	if err != nil {
		return nil, err
	}
	ctx := &visitor.EvaluationContext{Variables: values}
	solver := visitor.NewBooleanSolver(ctx)
	table := solver.Visit(ast)

	for _, entry := range table {
		slices.SortFunc(entry.Variables, sortVariables)
	}
	return table, nil
}

func sortVariables(a, b visitor.TruthTableVariable) int {
	return strings.Compare(a.Name, b.Name)
}

func ExtractVariables(expr string) []string {
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
