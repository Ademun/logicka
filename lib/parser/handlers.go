package parser

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
)

type StatementHandler func(p *Parser) (ast.Stmt, error)
type PrefixHandler func(p *Parser) (ast.Expr, error)
type InfixHandler func(p *Parser, left ast.Expr, precedence Precedence) (ast.Expr, error)

type HandlerRegistry struct {
	statements map[lexer.TokenType]StatementHandler
	prefixes   map[lexer.TokenType]PrefixHandler
	infixes    map[lexer.TokenType]InfixHandler
}

func NewHandlerRegistry() *HandlerRegistry {
	registry := &HandlerRegistry{
		statements: make(map[lexer.TokenType]StatementHandler),
		prefixes:   make(map[lexer.TokenType]PrefixHandler),
		infixes:    make(map[lexer.TokenType]InfixHandler),
	}
	registry.registerDefaults()
	return registry
}

func (hr *HandlerRegistry) RegisterStatement(tokenType lexer.TokenType, handler StatementHandler) {
	hr.statements[tokenType] = handler
}

func (hr *HandlerRegistry) RegisterPrefix(tokenType lexer.TokenType, handler PrefixHandler) {
	hr.prefixes[tokenType] = handler
}

func (hr *HandlerRegistry) RegisterInfix(tokenType lexer.TokenType, handler InfixHandler) {
	hr.infixes[tokenType] = handler
}

func (hr *HandlerRegistry) GetStatement(tokenType lexer.TokenType) (StatementHandler, bool) {
	handler, exists := hr.statements[tokenType]
	return handler, exists
}

func (hr *HandlerRegistry) GetPrefix(tokenType lexer.TokenType) (PrefixHandler, bool) {
	handler, exists := hr.prefixes[tokenType]
	return handler, exists
}

func (hr *HandlerRegistry) GetInfix(tokenType lexer.TokenType) (InfixHandler, bool) {
	handler, exists := hr.infixes[tokenType]
	return handler, exists
}

func (hr *HandlerRegistry) registerDefaults() {
	hr.RegisterInfix(lexer.GlAssignment, (*Parser).handleAssignment)

	infixTokens := []lexer.TokenType{
		lexer.BlImplication, lexer.BlEquivalence, lexer.BlDisjunction, lexer.BlConjunction,
		lexer.CdEquals, lexer.CdNotEquals, lexer.CdGreater, lexer.CdLess,
		lexer.CdGreaterOrEqual, lexer.CdLessOrEqual,
		lexer.StElementOf, lexer.StNotElementOf, lexer.StSubset, lexer.StSuperset,
		lexer.ArAddition, lexer.ArSubtraction, lexer.ArMultiplication,
		lexer.ArDivision, lexer.ArModulus, lexer.ArPower,
		lexer.StUnion, lexer.StIntersection,
	}

	for _, tokenType := range infixTokens {
		hr.RegisterInfix(tokenType, (*Parser).handleInfix)
	}

	prefixTokens := []lexer.TokenType{
		lexer.BlNegation, lexer.BlForall, lexer.BlExists,
	}

	for _, tokenType := range prefixTokens {
		hr.RegisterPrefix(tokenType, (*Parser).handlePrefix)
	}

	hr.RegisterPrefix(lexer.GlQuote, (*Parser).handleStringExpression)
	hr.RegisterPrefix(lexer.BlForall, (*Parser).handleQuantifierExpression)
	hr.RegisterPrefix(lexer.BlExists, (*Parser).handleQuantifierExpression)
	hr.RegisterPrefix(lexer.GlLeftParenthesis, (*Parser).handleGroupedExpression)
	hr.RegisterPrefix(lexer.GlLeftBrace, (*Parser).handleBracedExpression)
	hr.RegisterPrefix(lexer.GlIdentifier, (*Parser).handleIdentifier)
	hr.RegisterPrefix(lexer.LtNumber, (*Parser).handleLiteral)
	hr.RegisterPrefix(lexer.LtTrue, (*Parser).handleLiteral)
	hr.RegisterPrefix(lexer.LtFalse, (*Parser).handleLiteral)
}
