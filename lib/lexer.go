package lib

import (
	"fmt"
	"regexp"
	"strings"
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
}

type UnexpectedRuneError struct {
	r   rune
	pos int
}

func (e *UnexpectedRuneError) Error() string {
	return fmt.Sprintf("unexpected rune %c on pos %d", e.r, e.pos)
}

func Lex(str string) ([]Token, error) {
	str = sanitizeInput(str)
	var tokens []Token
	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch r {
		case '(':
			tokens = append(tokens, Token{Type: LPAREN, Value: string(r)})
		case ')':
			tokens = append(tokens, Token{Type: RPAREN, Value: string(r)})
		case 'A':
			tokens = append(tokens, Token{Type: FORALL, Value: string(r)})
		case 'E':
			tokens = append(tokens, Token{Type: EXISTS, Value: string(r)})
		case '-':
			if (i+1 < len(runes)) && runes[i+1] == '>' {
				tokens = append(tokens, Token{IMPL, "->"})
				i++
			} else {
				tokens = append(tokens, Token{NEG, string(r)})
			}
		case '!':
			tokens = append(tokens, Token{NEG, string(r)})
		case '~':
			tokens = append(tokens, Token{EQUIV, string(r)})
		case '&':
			tokens = append(tokens, Token{CONJ, string(r)})
		case '\\':
			if (i+1 < len(runes)) && runes[i+1] == '/' {
				tokens = append(tokens, Token{DISJ, "\\/"})
				i++
			} else {
				return nil, &UnexpectedRuneError{r, i}
			}
		case '1', '0':
			tokens = append(tokens, Token{Type: LIT, Value: string(r)})

		default:
			if unicode.IsLetter(r) {
				if unicode.IsLower(r) {
					var sb strings.Builder
					sb.WriteRune(r)
					for (i+1 < len(runes)) && unicode.IsLetter(runes[i+1]) {
						sb.WriteRune(runes[i+1])
						i++
					}
					tokens = append(tokens, Token{Type: VAR, Value: sb.String()})
				} else {
					tokens = append(tokens, Token{PRED, string(r)})
				}
			} else {
				return nil, &UnexpectedRuneError{r, i}
			}
		}
	}
	return tokens, nil
}

func sanitizeInput(input string) string {
	regex := regexp.MustCompile(`\s+`)
	return regex.ReplaceAllString(input, "")
}
