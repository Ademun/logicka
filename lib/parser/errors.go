package parser

import (
	"fmt"
	"logicka/lib/lexer"
)

type ParseError struct {
	Token   *lexer.Token
	Message string
}

func (e *ParseError) Error() string {
	if e.Token != nil {
		return fmt.Sprintf("parse error at pos %d: %s (token: \"%s\")",
			e.Token.Pos, e.Message, e.Token.Value)
	}
	return fmt.Sprintf("parse error: %s", e.Message)
}

func newParseError(token *lexer.Token, message string) *ParseError {
	return &ParseError{Token: token, Message: message}
}

func newUnexpectedTokenError(expected, actual *lexer.Token) *ParseError {
	message := fmt.Sprintf("expected %s, got %s", expected.Type, actual.Type)
	return newParseError(actual, message)
}
