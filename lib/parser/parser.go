package parser

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
)

type Parser struct {
	tokens []*lexer.Token
	pos    int
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(tokens []*lexer.Token) (ast.Stmt, error) {
	p.tokens = tokens
	p.pos = 0
}

func (p *Parser) current() *lexer.Token {
	return p.tokens[p.pos]
}

func (p *Parser) consume() *lexer.Token {
	tk := p.current()
	p.pos++
	return tk
}

func (p *Parser) peek() *lexer.Token {
	if p.canConsume() {
		return p.tokens[p.pos+1]
	}
	return nil
}

func (p *Parser) canConsume() bool {
	return p.pos < len(p.tokens) && p.current().Type != lexer.EOF
}
