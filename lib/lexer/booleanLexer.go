package lexer

import (
	"fmt"
	"unicode"
)

type BooleanTokenType int

const (
	LPAREN BooleanTokenType = iota // (
	RPAREN                         // )
	FORALL                         // A
	EXISTS                         // E
	IMPL                           // ->
	EQUIV                          // ~
	CONJ                           // &
	DISJ                           // \/
	NEG                            // ! -
	PRED
	VAR
	LIT // 1 0
)

func (t BooleanTokenType) String() string {
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
		return "→"
	case EQUIV:
		return "~"
	case CONJ:
		return "∧"
	case DISJ:
		return "∨"
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

type BooleanLexer struct {
	input []rune
	pos   int
}

func NewBooleanLexer(input string) *BooleanLexer {
	return &BooleanLexer{[]rune(input), 0}
}

func (l *BooleanLexer) Lex() ([]Token[BooleanTokenType], error) {
	var tokens []Token[BooleanTokenType]

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

func (l *BooleanLexer) NextToken() (Token[BooleanTokenType], error) {
	r := l.input[l.pos]
	startPos := l.pos

	switch r {
	case '(':
		l.pos++
		return Token[BooleanTokenType]{Type: LPAREN, Value: "(", Pos: startPos}, nil
	case ')':
		l.pos++
		return Token[BooleanTokenType]{Type: RPAREN, Value: ")", Pos: startPos}, nil
	case 'A':
		l.pos++
		return Token[BooleanTokenType]{Type: FORALL, Value: "A", Pos: startPos}, nil
	case 'E':
		l.pos++
		return Token[BooleanTokenType]{Type: EXISTS, Value: "E", Pos: startPos}, nil
	case '-':
		return l.lexImpl()
	case '!':
		l.pos++
		return Token[BooleanTokenType]{Type: NEG, Value: "!", Pos: startPos}, nil
	case '~':
		l.pos++
		return Token[BooleanTokenType]{Type: EQUIV, Value: "~", Pos: startPos}, nil
	case '&':
		l.pos++
		return Token[BooleanTokenType]{Type: CONJ, Value: "&", Pos: startPos}, nil
	case '\\':
		return l.lexDisj()
	case '1', '0':
		l.pos++
		return Token[BooleanTokenType]{Type: LIT, Value: string(r), Pos: startPos}, nil
	default:
		if unicode.IsLetter(r) {
			return l.lexIdentifier()
		}
		return Token[BooleanTokenType]{}, UnexpectedRuneError{pos: l.pos, r: r}
	}
}

func (l *BooleanLexer) lexImpl() (Token[BooleanTokenType], error) {
	startPos := l.pos
	if l.pos+1 < len(l.input) && l.input[l.pos+1] == '>' {
		l.pos += 2
		return Token[BooleanTokenType]{Type: IMPL, Value: "->", Pos: startPos}, nil
	}
	l.pos++
	return Token[BooleanTokenType]{Type: NEG, Value: "-", Pos: startPos}, nil
}

// lexDisj обрабатывает дизъюнкцию (\/)
func (l *BooleanLexer) lexDisj() (Token[BooleanTokenType], error) {
	startPos := l.pos
	if l.pos+1 < len(l.input) && l.input[l.pos+1] == '/' {
		l.pos += 2
		return Token[BooleanTokenType]{Type: DISJ, Value: "\\/", Pos: startPos}, nil
	}
	return Token[BooleanTokenType]{}, fmt.Errorf("expected '\\/' at position %d", startPos)
}

func (l *BooleanLexer) lexIdentifier() (Token[BooleanTokenType], error) {
	r := l.input[l.pos]
	startPos := l.pos
	l.pos++

	if unicode.IsLower(r) {
		return Token[BooleanTokenType]{Type: VAR, Value: string(r), Pos: startPos}, nil
	}
	return Token[BooleanTokenType]{Type: PRED, Value: string(r), Pos: startPos}, nil
}
