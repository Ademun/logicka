package lib

import (
	"fmt"
	"logicka/lib/lexer"
	"logicka/lib/parser"

	"github.com/sanity-io/litter"
)

type Logicka struct {
}

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

	return nil, nil
}
