package parser

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
)

type Parser struct {
	tokens   []*lexer.Token
	position int
	registry *HandlerRegistry
}

func NewParser() *Parser {
	return &Parser{
		registry: NewHandlerRegistry(),
	}
}

func (p *Parser) Parse(tokens []*lexer.Token) (ast.Stmt, error) {
	if len(tokens) == 0 {
		return nil, newParseError(nil, "no tokens to parse")
	}

	p.tokens = tokens
	p.position = 0

	var statements []ast.Stmt

	for p.canConsume() {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	return ast.NewBlockStmt(statements), nil
}

func (p *Parser) current() *lexer.Token {
	if p.position >= len(p.tokens) {
		return &lexer.Token{Type: lexer.EOF}
	}
	return p.tokens[p.position]
}

func (p *Parser) consume() *lexer.Token {
	token := p.current()
	p.position++
	return token
}

func (p *Parser) peek() *lexer.Token {
	if p.position+1 >= len(p.tokens) {
		return &lexer.Token{Type: lexer.EOF}
	}
	return p.tokens[p.position+1]
}

func (p *Parser) canConsume() bool {
	return p.position < len(p.tokens) && p.current().Type != lexer.EOF
}

func (p *Parser) expect(expected lexer.TokenType) (*lexer.Token, error) {
	current := p.current()
	if current.Type != expected {
		return nil, newUnexpectedTokenError(&lexer.Token{Type: expected}, current)
	}
	return p.consume(), nil
}
