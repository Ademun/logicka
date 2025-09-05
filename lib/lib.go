package lib

import (
	"regexp"
	"slices"
	"strings"
)

func GenerateTruthTable(expr string, values map[string]bool) ([]TruthTableEntry, error) {
	tokens, err := Lex(expr)
	if err != nil {
		return nil, err
	}
	parser := &Parser{tokens, 0}
	ast, err := parser.ParseExpression()
	if err != nil {
		return nil, err
	}
	ctx := &EvaluationContext{values}
	table := ast.Accept(&BooleanSolver{}, ctx).([]TruthTableEntry)

	for _, entry := range table {
		slices.SortFunc(entry.Variables, sortVariables)
	}
	return table, nil
}

func sortVariables(a, b TruthTableVariable) int {
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
