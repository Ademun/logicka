package parser

import (
	"fmt"
	"logicka/lib/ast"
	"logicka/lib/lexer"
)

type Parser struct {
	Tokens []lexer.Token
	pos    int
}

func (p *Parser) current() lexer.Token {
	if p.pos >= len(p.Tokens) {
		return lexer.Token{Type: lexer.EOF, Value: ""}
	}
	return p.Tokens[p.pos]
}

func (p *Parser) peek() lexer.Token {
	if p.pos+1 >= len(p.Tokens) {
		return lexer.Token{Type: lexer.EOF, Value: ""}
	}
	return p.Tokens[p.pos+1]
}

func (p *Parser) advance() {
	if p.pos < len(p.Tokens) {
		p.pos++
	}
}

func (p *Parser) expect(tokenType lexer.TokenType) error {
	if p.current().Type != tokenType {
		return fmt.Errorf("expected %s, got %s", tokenType.String(), p.current().Type.String())
	}
	p.advance()
	return nil
}

// <expr> ::= <equal>
func (p *Parser) ParseExpression() (ast.ASTNode, error) {
	expr, err := p.parseEqual()
	if err != nil {
		return nil, err
	}

	if p.current().Type != lexer.EOF {
		return nil, fmt.Errorf("unexpected token %s at pos %d, expected end of expression", p.current().Type.String(), p.current().Pos)
	}

	return expr, nil
}

// <equal> ::= <impl> ("~" <impl>)*
func (p *Parser) parseEqual() (ast.ASTNode, error) {
	left, err := p.parseImpl()
	if err != nil {
		return nil, err
	}

	for p.current().Type == lexer.EQUIV {
		p.advance() // consume "~"
		right, err := p.parseImpl()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryNode{Operator: lexer.EQUIV, Left: left, Right: right}
	}

	return left, nil
}

// <impl> ::= <or> ("->" <or>)*
func (p *Parser) parseImpl() (ast.ASTNode, error) {
	left, err := p.parseOr()
	if err != nil {
		return nil, err
	}

	for p.current().Type == lexer.IMPL {
		p.advance() // consume "->"
		right, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryNode{Operator: lexer.IMPL, Left: left, Right: right}
	}

	return left, nil
}

// <or> ::= <and> ("\\/" <and>)*
func (p *Parser) parseOr() (ast.ASTNode, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	for p.current().Type == lexer.DISJ {
		p.advance() // consume "\\/"
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryNode{Operator: lexer.DISJ, Left: left, Right: right}
	}

	return left, nil
}

// <and> ::= <not> ("&" <not>)*
func (p *Parser) parseAnd() (ast.ASTNode, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}

	for p.current().Type == lexer.CONJ {
		p.advance() // consume "&"
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryNode{Operator: lexer.CONJ, Left: left, Right: right}
	}

	return left, nil
}

// <not> ::= "-" <pred> | <pred>
func (p *Parser) parseNot() (ast.ASTNode, error) {
	for p.current().Type == lexer.NEG {
		p.advance() // consume "-"
		expr, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryNode{Operator: lexer.NEG, Operand: expr}, nil
	}
	return p.parsePred()
}

// <pred> ::= [A-Z] <quant> | [A-Z] "(" <quant> ")" | <quant>
func (p *Parser) parsePred() (ast.ASTNode, error) {
	if p.current().Type == lexer.PRED {
		predName := p.current().Value
		p.advance()

		// Check if this is a predicate with arguments: [A-Z] "(" <quant> ")"
		if p.current().Type == lexer.LPAREN {
			p.advance() // consume "("

			var args []string
			// Parse arguments (variables)
			for p.current().Type == lexer.VAR {
				args = append(args, p.current().Value)
				p.advance()
			}

			if err := p.expect(lexer.RPAREN); err != nil {
				return nil, err
			}
			return &ast.PredicateNode{Name: predName, Body: args}, nil
		} else {
			// This is a predicate followed by a quantifier: [A-Z] <quant>
			quant, err := p.parseQuant()
			if err != nil {
				return nil, err
			}
			return &ast.PredicateNode{Name: predName, Body: quant}, nil
		}
	}
	return p.parseQuant()
}

// <quant> ::= ("A" | "E") <primary> | ("A" | "E") "(" <primary> ")" | <primary>
func (p *Parser) parseQuant() (ast.ASTNode, error) {
	if p.current().Type == lexer.FORALL || p.current().Type == lexer.EXISTS {
		quantType := p.current().Type
		p.advance()

		// Check for quantifier with parentheses: ("A" | "E") "(" <primary> ")"
		if p.current().Type == lexer.LPAREN {
			p.advance() // consume "("
			// Expect a variable
			if p.current().Type != lexer.VAR {
				return nil, fmt.Errorf("expected variable after quantifier, got %s", p.current().Type.String())
			}
			variable := p.current().Value
			p.advance()

			if err := p.expect(lexer.RPAREN); err != nil {
				return nil, err
			}

			return &ast.QuantifierNode{Type: quantType, Variable: variable, Domain: nil}, nil
		} else {
			// Simple quantifier: ("A" | "E") <primary>
			body, err := p.parsePrimary()
			if err != nil {
				return nil, err
			}
			return &ast.QuantifierNode{Type: quantType, Variable: "", Domain: body}, nil
		}
	}
	return p.parsePrimary()
}

// <primary> ::= [a-z] | "(" <expr> ")"
func (p *Parser) parsePrimary() (ast.ASTNode, error) {
	if p.current().Type == lexer.LIT {
		literal := p.current().Value
		p.advance()
		if literal == "1" {
			return &ast.LiteralNode{Value: true}, nil
		}
		return &ast.LiteralNode{Value: false}, nil
	}
	if p.current().Type == lexer.VAR {
		variable := p.current().Value
		p.advance()
		return &ast.VariableNode{Name: variable}, nil
	}

	if p.current().Type == lexer.LPAREN {
		p.advance() // consume "("
		expr, err := p.parseEqual()
		if err != nil {
			return nil, err
		}
		if err := p.expect(lexer.RPAREN); err != nil {
			return nil, err
		}
		return &ast.GroupingNode{Expr: expr}, nil
	}

	return nil, fmt.Errorf("expected variable or '(', got %s", p.current().Type.String())
}
