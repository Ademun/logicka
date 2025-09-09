package lexer

import (
	"fmt"
	"unicode"
)

type TokenType int

const (
	LPAREN TokenType = iota // (
	RPAREN                  // )
	FORALL                  // A
	EXISTS                  // E
	IMPL                    // ->
	EQUIV                   // ~
	CONJ                    // &
	DISJ                    // \/
	NEG                     // ! -
	PRED
	VAR
	LIT // 1 0
	EOF
)

func (t TokenType) String() string {
	switch t {
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case FORALL:
		return "A"
	case EXISTS:
		return "E"
	case IMPL:
		return "->"
	case EQUIV:
		return "~"
	case CONJ:
		return "&"
	case DISJ:
		return "\\/"
	case NEG:
		return "!"
	case PRED:
		return "PREDICATE"
	case VAR:
		return "VARIABLE"
	case LIT:
		return "LITERAL"
	case EOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}

type Token struct {
	Type  TokenType
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

type Lexer struct {
	input []rune
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{[]rune(input), 0}
}

func (l *Lexer) Lex() ([]Token, error) {
	var tokens []Token

	for l.pos < len(l.input) {
		if unicode.IsSpace(l.input[l.pos]) {
			l.pos++
			continue
		}
		token, err := l.NextToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (l *Lexer) NextToken() (Token, error) {
	r := l.input[l.pos]
	startPos := l.pos

	switch r {
	case '(':
		l.pos++
		return Token{Type: LPAREN, Value: "(", Pos: startPos}, nil
	case ')':
		l.pos++
		return Token{Type: RPAREN, Value: ")", Pos: startPos}, nil
	case 'A':
		l.pos++
		return Token{Type: FORALL, Value: "A", Pos: startPos}, nil
	case 'E':
		l.pos++
		return Token{Type: EXISTS, Value: "E", Pos: startPos}, nil
	case '-':
		return l.lexImpl()
	case '!':
		l.pos++
		return Token{Type: NEG, Value: "!", Pos: startPos}, nil
	case '~':
		l.pos++
		return Token{Type: EQUIV, Value: "~", Pos: startPos}, nil
	case '&':
		l.pos++
		return Token{Type: CONJ, Value: "&", Pos: startPos}, nil
	case '\\':
		return l.lexDisj()
	case '1', '0':
		l.pos++
		return Token{Type: LIT, Value: string(r), Pos: startPos}, nil
	default:
		if unicode.IsLetter(r) {
			return l.lexIdentifier()
		}
		return Token{}, UnexpectedRuneError{pos: l.pos, r: r}
	}
}

func (l *Lexer) lexImpl() (Token, error) {
	startPos := l.pos
	if l.pos+1 < len(l.input) && l.input[l.pos+1] == '>' {
		l.pos += 2
		return Token{Type: IMPL, Value: "->", Pos: startPos}, nil
	}
	l.pos++
	return Token{Type: NEG, Value: "-", Pos: startPos}, nil
}

// lexDisj обрабатывает дизъюнкцию (\/)
func (l *Lexer) lexDisj() (Token, error) {
	startPos := l.pos
	if l.pos+1 < len(l.input) && l.input[l.pos+1] == '/' {
		l.pos += 2
		return Token{Type: DISJ, Value: "\\/", Pos: startPos}, nil
	}
	return Token{}, fmt.Errorf("expected '\\/' at position %d", startPos)
}

func (l *Lexer) lexIdentifier() (Token, error) {
	r := l.input[l.pos]
	startPos := l.pos
	l.pos++

	if unicode.IsLower(r) {
		return Token{Type: VAR, Value: string(r), Pos: startPos}, nil
	}
	return Token{Type: PRED, Value: string(r), Pos: startPos}, nil
}
