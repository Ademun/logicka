package lib

import (
	"fmt"
)

type Parser struct {
	Tokens []Token
	pos    int
}

func (p *Parser) current() Token {
	if p.pos >= len(p.Tokens) {
		return Token{Type: EOF, Value: ""}
	}
	return p.Tokens[p.pos]
}

func (p *Parser) peek() Token {
	if p.pos+1 >= len(p.Tokens) {
		return Token{Type: EOF, Value: ""}
	}
	return p.Tokens[p.pos+1]
}

func (p *Parser) advance() {
	if p.pos < len(p.Tokens) {
		p.pos++
	}
}

func (p *Parser) expect(tokenType TokenType) error {
	if p.current().Type != tokenType {
		return fmt.Errorf("expected %s, got %s", tokenType.String(), p.current().Type.String())
	}
	p.advance()
	return nil
}

// <expr> ::= <equal>
func (p *Parser) ParseExpression() (ASTNode, error) {
	return p.parseEqual()
}

// <equal> ::= <impl> ("~" <impl>)*
func (p *Parser) parseEqual() (ASTNode, error) {
	left, err := p.parseImpl()
	if err != nil {
		return nil, err
	}

	for p.current().Type == EQUIV {
		p.advance() // consume "~"
		right, err := p.parseImpl()
		if err != nil {
			return nil, err
		}
		left = &BinaryNode{Operator: EQUIV, Left: left, Right: right}
	}

	return left, nil
}

// <impl> ::= <or> ("->" <or>)*
func (p *Parser) parseImpl() (ASTNode, error) {
	left, err := p.parseOr()
	if err != nil {
		return nil, err
	}

	for p.current().Type == IMPL {
		p.advance() // consume "->"
		right, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		left = &BinaryNode{Operator: IMPL, Left: left, Right: right}
	}

	return left, nil
}

// <or> ::= <and> ("\\/" <and>)*
func (p *Parser) parseOr() (ASTNode, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	for p.current().Type == DISJ {
		p.advance() // consume "\\/"
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &BinaryNode{Operator: DISJ, Left: left, Right: right}
	}

	return left, nil
}

// <and> ::= <not> ("&" <not>)*
func (p *Parser) parseAnd() (ASTNode, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}

	for p.current().Type == CONJ {
		p.advance() // consume "&"
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		left = &BinaryNode{Operator: CONJ, Left: left, Right: right}
	}

	return left, nil
}

// <not> ::= "-" <pred> | <pred>
func (p *Parser) parseNot() (ASTNode, error) {
	for p.current().Type == NEG {
		p.advance() // consume "-"
		expr, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return &UnaryNode{Operator: NEG, Operand: expr}, nil
	}
	return p.parsePred()
}

// <pred> ::= [A-Z] <quant> | [A-Z] "(" <quant> ")" | <quant>
func (p *Parser) parsePred() (ASTNode, error) {
	if p.current().Type == PRED {
		predName := p.current().Value
		p.advance()

		// Check if this is a predicate with arguments: [A-Z] "(" <quant> ")"
		if p.current().Type == LPAREN {
			p.advance() // consume "("

			var args []string
			// Parse arguments (variables)
			for p.current().Type == VAR {
				args = append(args, p.current().Value)
				p.advance()
			}

			if err := p.expect(RPAREN); err != nil {
				return nil, err
			}
			return &PredicateNode{Name: predName, Body: args}, nil
		} else {
			// This is a predicate followed by a quantifier: [A-Z] <quant>
			quant, err := p.parseQuant()
			if err != nil {
				return nil, err
			}
			return &PredicateNode{Name: predName, Body: quant}, nil
		}
	}
	return p.parseQuant()
}

// <quant> ::= ("A" | "E") <primary> | ("A" | "E") "(" <primary> ")" | <primary>
func (p *Parser) parseQuant() (ASTNode, error) {
	if p.current().Type == FORALL || p.current().Type == EXISTS {
		quantType := p.current().Type
		p.advance()

		// Check for quantifier with parentheses: ("A" | "E") "(" <primary> ")"
		if p.current().Type == LPAREN {
			p.advance() // consume "("
			// Expect a variable
			if p.current().Type != VAR {
				return nil, fmt.Errorf("expected variable after quantifier, got %s", p.current().Type.String())
			}
			variable := p.current().Value
			p.advance()

			if err := p.expect(RPAREN); err != nil {
				return nil, err
			}

			return &QuantifierNode{Type: quantType, Variable: variable, Domain: nil}, nil
		} else {
			// Simple quantifier: ("A" | "E") <primary>
			body, err := p.parsePrimary()
			if err != nil {
				return nil, err
			}
			return &QuantifierNode{Type: quantType, Variable: "", Domain: body}, nil
		}
	}
	return p.parsePrimary()
}

// <primary> ::= [a-z] | "(" <expr> ")"
func (p *Parser) parsePrimary() (ASTNode, error) {
	if p.current().Type == VAR {
		variable := p.current().Value
		p.advance()
		return &VariableNode{Name: variable}, nil
	}

	if p.current().Type == LPAREN {
		p.advance() // consume "("
		expr, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.expect(RPAREN); err != nil {
			return nil, err
		}
		return &GroupingNode{Expr: expr}, nil
	}
	return nil, fmt.Errorf("expected variable or '(', got %s", p.current().Type.String())
}
