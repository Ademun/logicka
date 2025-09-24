package parser

import (
	"fmt"
	"logicka/lib/ast"
	"logicka/lib/lexer"
	"strconv"
)

func (p *Parser) parseExpression(minPrecedence Precedence) (ast.Expr, error) {
	token := p.current()
	handler, exists := p.registry.GetPrefix(token.Type)
	if !exists {
		return nil, newParseError(token, "unexpected token, expected expression")
	}

	left, err := handler(p)
	if err != nil {
		return nil, err
	}

	for tokenPrecedences.Get(p.current().Type) > minPrecedence {
		token = p.current()
		infixHandler, exists := p.registry.GetInfix(token.Type)
		if !exists {
			break
		}

		left, err = infixHandler(p, left, tokenPrecedences.Get(p.current().Type))
		if err != nil {
			return nil, err
		}
	}

	return left, nil
}

func (p *Parser) handleAssignment(left ast.Expr, precedence Precedence) (ast.Expr, error) {
	p.consume()

	value, err := p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}

	return ast.NewAssignmentExpr(left, value)
}

func (p *Parser) handlePrefix() (ast.Expr, error) {
	token := p.consume()
	operand, err := p.parseExpression(PrecedencePrefix)
	if err != nil {
		return nil, err
	}

	return ast.NewUnaryExpr(*token, operand), nil
}

func (p *Parser) handleInfix(left ast.Expr, precedence Precedence) (ast.Expr, error) {
	token := p.consume()
	right, err := p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}

	return ast.NewBinaryExpr(*token, left, right), nil
}

func (p *Parser) handleStringExpression() (ast.Expr, error) {
	p.consume()

	body, err := p.parseExpression(PrecedenceDefault)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lexer.GlQuote); err != nil {
		return nil, err
	}

	return ast.NewStringExpr(body.String()), nil
}

func (p *Parser) handleQuantifierExpression() (ast.Expr, error) {
	token := p.consume()
	variable, err := p.expect(lexer.GlIdentifier)
	if err != nil {
		return nil, err
	}
	fmt.Println(variable)
	var domain ast.Expr
	if p.current().Type == lexer.StElementOf {
		p.consume()
		domain, err = p.handleIdentifier()
		if err != nil {
			return nil, err
		}
	}
	body, err := p.parseExpression(PrecedencePrefix)
	if err != nil {
		return nil, err
	}

	return ast.NewQuantifierExpr(token.Type, ast.NewIdentifierExpr(variable.Value), domain, body)
}

func (p *Parser) handleGroupedExpression() (ast.Expr, error) {
	p.consume()

	body, err := p.parseExpression(PrecedenceDefault)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(lexer.GlRightParenthesis); err != nil {
		return nil, err
	}

	return ast.NewGroupingExpr(body), nil
}

func (p *Parser) handleBracedExpression() (ast.Expr, error) {
	p.consume()

	var elements []ast.Expr

	if p.current().Type == lexer.GlRightBrace {
		p.consume()
		return ast.NewBracedExpression(elements), nil
	}

	for {
		element, err := p.parseExpression(PrecedenceAssignment)
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)

		if p.current().Type == lexer.GlRightBrace {
			p.consume()
			break
		}

		if _, err := p.expect(lexer.GlComma); err != nil {
			return nil, err
		}
	}

	return ast.NewBracedExpression(elements), nil
}

func (p *Parser) handleIdentifier() (ast.Expr, error) {
	name := p.consume()

	if p.current().Type != lexer.GlLeftParenthesis {
		return ast.NewIdentifierExpr(name.Value), nil
	}

	args, err := p.parseArgumentList()
	if err != nil {
		return nil, err
	}

	if p.current().Type == lexer.GlAssignment {
		p.consume()
		body, err := p.parseExpression(PrecedenceDefault)
		if err != nil {
			return nil, err
		}
		return ast.NewFunctionDeclExpr(name.Value, args, body)
	}

	return ast.NewFunctionCallExpr(name.Value, args), nil
}

func (p *Parser) parseArgumentList() ([]ast.Expr, error) {
	if _, err := p.expect(lexer.GlLeftParenthesis); err != nil {
		return nil, err
	}

	var args []ast.Expr

	if p.current().Type == lexer.GlRightParenthesis {
		p.consume()
		return args, nil
	}

	for {
		arg, err := p.parseExpression(PrecedenceDefault)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		if p.current().Type == lexer.GlRightParenthesis {
			p.consume()
			break
		}

		if _, err := p.expect(lexer.GlComma); err != nil {
			return nil, err
		}
	}

	return args, nil
}

func (p *Parser) handleLiteral() (ast.Expr, error) {
	literal := p.consume()

	switch literal.Type {
	case lexer.LtNumber:
		value, err := strconv.ParseFloat(literal.Value, 64)
		if err != nil {
			return nil, newParseError(literal, "invalid number format")
		}
		return ast.NewNumberExpr(value), nil

	case lexer.LtTrue:
		return ast.NewBooleanExpr(true), nil

	case lexer.LtFalse:
		return ast.NewBooleanExpr(false), nil

	default:
		return nil, newParseError(literal, "unknown literal type")
	}
}
