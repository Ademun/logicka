package lexer

import "fmt"

const EOF = -1

type Token[T any] struct {
	Type  T
	Value string
	Pos   int
}

type UnexpectedRuneError struct {
	r   rune
	pos int
}

func (e UnexpectedRuneError) Error() string {
	return fmt.Sprintf("unexpected rune %c on pos %d", e.r, e.pos)
}

type Lexer[T any] interface {
	Lex(string) (Token[T], error)
}
