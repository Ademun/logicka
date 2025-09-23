package parser

import (
	"logicka/lib/ast"
	"logicka/lib/lexer"
)

func (p *Parser) parseStatement() (ast.Stmt, error) {
	if handler, exists := p.registry.GetStatement(p.current().Type); exists {
		return handler(p)
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseExpressionStatement() (ast.Stmt, error) {
	expr, err := p.parseExpression(PrecedenceDefault)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(lexer.GlSemicolon); err != nil {
		return nil, err
	}

	return ast.NewExprStmt(expr), nil
}
