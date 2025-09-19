package lexer

import "regexp"

type TokenType int

const (
	ER_INVALID TokenType = iota
	// Global syntax tokens
	GL_LEFT_PARENTHESIS
	GL_RIGHT_PARENTHESIS
	GL_LEFT_BRACE
	GL_RIGHT_BRACE
	GL_COMMA
	GL_DOT
	GL_PIPE
	GL_ASSIGNMENT
	GL_IDENTIFIER
	// Boolean syntax tokens
	BL_CONJUNCTION
	BL_DISJUNCTION
	BL_NEGATION
	BL_EQUIVALENCE
	BL_IMPLICATION
	BL_FORALL
	BL_EXISTS
	// Arithmetic syntax tokens
	AR_ADDITION
	AR_SUBSTRACTION
	AR_MULTIPLICATION
	AR_DIVISION
	AR_MODULUS
	AR_POWER
	// Conditional syntax tokens
	CD_EQUALS
	CD_NOT_EQUALS
	CD_GREATER
	CD_LESS
	CD_GREATER_OR_EQUAL
	CD_LESS_OR_EQUAL
	CD_IF
	CD_THEN
	CD_ELSE
	// Set syntax tockens
	ST_ELEMENT_OF
	ST_NOT_ELEMENT_OF
	ST_UNION
	ST_INTERSECTION
	ST_SUBSET
	ST_SUPERSET
	// Literal syntax tokens
	LT_TRUE
	LT_FALSE
	LT_NUMBER
)

type Token struct {
	Type  TokenType
	Value string
	Pos   int
}

func NewToken(ttype TokenType, value string, pos int) *Token {
	return &Token{
		Type:  ttype,
		Value: value,
		Pos:   pos,
	}
}

var symbolTokens = map[string]TokenType{
	"(":  GL_LEFT_PARENTHESIS,
	")":  GL_RIGHT_PARENTHESIS,
	"{":  GL_LEFT_BRACE,
	"}":  GL_RIGHT_BRACE,
	",":  GL_COMMA,
	".":  GL_DOT,
	"|":  GL_PIPE,
	":=": GL_ASSIGNMENT,
	"==": CD_EQUALS,
	"!=": CD_NOT_EQUALS,
	"+":  AR_ADDITION,
	"-":  AR_SUBSTRACTION,
	"*":  AR_MULTIPLICATION,
	"/":  AR_DIVISION,
	"%":  AR_MODULUS,
	"^":  AR_POWER,
	">":  CD_GREATER,
	"<":  CD_LESS,
	">=": CD_GREATER_OR_EQUAL,
	"<=": CD_LESS_OR_EQUAL,
}

var maxSymbolLength int

func init() {
	maxSymbolLength = 0
	for symbol := range symbolTokens {
		if len(symbol) > maxSymbolLength {
			maxSymbolLength = len(symbol)
		}
	}
}

var keywordTokens = map[string]TokenType{
	"conjunction":    BL_CONJUNCTION,
	"disjunction":    BL_DISJUNCTION,
	"negation":       BL_NEGATION,
	"equivalence":    BL_EQUIVALENCE,
	"implication":    BL_IMPLICATION,
	"forall":         BL_FORALL,
	"exists":         BL_EXISTS,
	"element_of":     ST_ELEMENT_OF,
	"not_element_of": ST_NOT_ELEMENT_OF,
	"union":          ST_UNION,
	"intersection":   ST_INTERSECTION,
	"subset":         ST_SUBSET,
	"superset":       ST_SUPERSET,
	"if":             CD_IF,
	"then":           CD_THEN,
	"else":           CD_ELSE,
	"T":              LT_TRUE,
	"F":              LT_FALSE,
}

var (
	numberRegex     = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?`)
	identifierRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`)
)
