package lexer

import "unicode"

type SetTokenType int

const (
	LBRACE    SetTokenType = iota // {
	RBRACE                        // }
	NAME                          // Uppercase single letter
	ASSIGN                        // =
	ELEMENT                       // Lowercase single letter or any comparable type
	COMMA                         // ,
	IN                            // ∈
	UNION                         // ∪
	INTERSECT                     // ∩
	SYMDIFF                       // ⊕
	ADD                           // +
	SUBSTRACT                     // -
)

func (t SetTokenType) String() string {
	switch t {
	case LBRACE:
		return "{"
	case RBRACE:
		return "}"
	case NAME:
		return "NAME"
	case ASSIGN:
		return "="
	case ELEMENT:
		return "ELEMENT"
	case COMMA:
		return ","
	case IN:
		return "∈"
	case UNION:
		return "∪"
	case INTERSECT:
		return "∩"
	case SYMDIFF:
		return "⊕"
	case ADD:
		return "+"
	case SUBSTRACT:
		return "-"
	default:
		return "UNKNOWN"
	}
}

type SetLexer struct {
	input []rune
	pos   int
}

func NewSetLexer(input string) *SetLexer {
	return &SetLexer{input: []rune(input), pos: 0}
}

func (l *SetLexer) Lex() ([]Token[SetTokenType], error) {
	var tokens []Token[SetTokenType]

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

func (l *SetLexer) NextToken() (Token[SetTokenType], error) {
	r := l.input[l.pos]
	startPos := l.pos

	switch r {
	case '{':
		l.pos++
		return Token[SetTokenType]{Type: LBRACE, Value: "{", Pos: startPos}, nil
	case '}':
		l.pos++
		return Token[SetTokenType]{Type: RBRACE, Value: "}", Pos: startPos}, nil
	case '=':
		l.pos++
		return Token[SetTokenType]{Type: ASSIGN, Value: "=", Pos: startPos}, nil
	case ',':
		l.pos++
		return Token[SetTokenType]{Type: COMMA, Value: ",", Pos: startPos}, nil
	case '∈':
		l.pos++
		return Token[SetTokenType]{Type: IN, Value: "∈", Pos: startPos}, nil
	case '∪':
		l.pos++
		return Token[SetTokenType]{Type: UNION, Value: "∪", Pos: startPos}, nil
	case '∩':
		l.pos++
		return Token[SetTokenType]{Type: INTERSECT, Value: "∩", Pos: startPos}, nil
	case '⊕':
		l.pos++
		return Token[SetTokenType]{Type: SYMDIFF, Value: "⊕", Pos: startPos}, nil
	case '+':
		l.pos++
		return Token[SetTokenType]{Type: ADD, Value: "+", Pos: startPos}, nil
	case '-':
		l.pos++
		return Token[SetTokenType]{Type: SUBSTRACT, Value: "-", Pos: startPos}, nil
	default:
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return l.lexIdentifier()
		}
		return Token[SetTokenType]{}, UnexpectedRuneError{pos: l.pos, r: r}
	}
}

func (l *SetLexer) lexIdentifier() (Token[SetTokenType], error) {
	r := l.input[l.pos]
	startPos := l.pos
	l.pos++

	if unicode.IsLower(r) || unicode.IsNumber(r) {
		return Token[SetTokenType]{Type: ELEMENT, Value: string(r), Pos: startPos}, nil
	}
	return Token[SetTokenType]{Type: NAME, Value: string(r), Pos: startPos}, nil
}
