package lexer

import (
	"fmt"
	"unicode"
)

type UnknownSymbolError struct {
	symbol string
	pos    int
}

func (e UnknownSymbolError) Error() string {
	return fmt.Sprintf("unknown symbol '%s' at position %d", e.symbol, e.pos)
}

type Lexer struct {
	input string
	pos   int
}

func NewLexer() *Lexer {
	return &Lexer{
		pos: 0,
	}
}

func (l *Lexer) Tokenize(input string) ([]*Token, error) {
	l.input = input
	l.pos = 0

	tokens := make([]*Token, 0)
	for l.pos < len(l.input) {
		if l.IsWhitespace() {
			l.pos++
			continue
		}

		if token := l.TryParseSymbolToken(); token != nil {
			tokens = append(tokens, token)
			continue
		}

		if token := l.TryParseNumberToken(); token != nil {
			tokens = append(tokens, token)
			continue
		}

		if token := l.TryParseIdentifierOrKeywordToken(); token != nil {
			tokens = append(tokens, token)
			continue
		}

		return nil, &UnknownSymbolError{symbol: string(l.input[l.pos]), pos: l.pos}
	}
	return tokens, nil
}

func (l *Lexer) IsWhitespace() bool {
	return unicode.IsSpace(rune(l.input[l.pos]))
}

func (l *Lexer) TryParseSymbolToken() *Token {
	start := l.pos
	for scope := maxSymbolLength; scope > 0; scope-- {
		end := l.pos + scope
		if end > len(l.input) {
			continue
		}
		candidate := l.input[l.pos:end]
		if found, ok := symbolTokens[candidate]; ok {
			l.pos += scope
			return NewToken(found, candidate, start)
		}
	}
	return nil
}

func (l *Lexer) TryParseNumberToken() *Token {
	if match := numberRegex.FindString(l.input[l.pos:]); match != "" {
		start := l.pos
		l.pos += len(match)
		return NewToken(LT_NUMBER, match, start)
	}
	return nil
}

func (l *Lexer) TryParseIdentifierOrKeywordToken() *Token {
	if match := identifierRegex.FindString(l.input[l.pos:]); match != "" {
		start := l.pos
		l.pos += len(match)
		if found, ok := keywordTokens[match]; ok {
			return NewToken(found, match, start)
		}
		return NewToken(GL_IDENTIFIER, match, start)
	}
	return nil
}
