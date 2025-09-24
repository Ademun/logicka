package lexer

import (
	"fmt"
	"regexp"
)

type TokenType int

const (
	ErInvalid TokenType = iota
	// Global syntax tokens
	GlLeftParenthesis
	GlRightParenthesis
	GlLeftBrace
	GlRightBrace
	GlComma
	GlSemicolon
	GlDot
	GlPipe
	GlAssignment
	GlIdentifier
	// Boolean syntax tokens
	BlConjunction
	BlDisjunction
	BlNegation
	BlEquivalence
	BlImplication
	BlForall
	BlExists
	// Arithmetic syntax tokens
	ArAddition
	ArSubtraction
	ArMultiplication
	ArDivision
	ArModulus
	ArPower
	// Conditional syntax tokens
	CdEquals
	CdNotEquals
	CdGreater
	CdLess
	CdGreaterOrEqual
	CdLessOrEqual
	CdIf
	CdThen
	CdElse
	// Set syntax tockens
	StElementOf
	StNotElementOf
	StUnion
	StIntersection
	StSubset
	StSuperset
	// Literal syntax tokens
	LtTrue
	LtFalse
	LtNumber
	EOF
)

func (tt TokenType) String() string {
	switch tt {
	case ErInvalid:
		return "ErInvalid"
	case GlLeftParenthesis:
		return "("
	case GlRightParenthesis:
		return ")"
	case GlLeftBrace:
		return "{"
	case GlRightBrace:
		return "}"
	case GlComma:
		return ","
	case GlSemicolon:
		return ";"
	case GlDot:
		return "."
	case GlPipe:
		return "|"
	case GlAssignment:
		return ":="
	case GlIdentifier:
		return "GlIdentifier"
	case BlConjunction:
		return "⋀"
	case BlDisjunction:
		return "⋁"
	case BlNegation:
		return "¬"
	case BlEquivalence:
		return "↔"
	case BlImplication:
		return "→"
	case BlForall:
		return "∀"
	case BlExists:
		return "∃"
	case ArAddition:
		return "+"
	case ArSubtraction:
		return "-"
	case ArMultiplication:
		return "*"
	case ArDivision:
		return "/"
	case ArModulus:
		return "%"
	case ArPower:
		return "^"
	case CdEquals:
		return "=="
	case CdNotEquals:
		return "!="
	case CdGreater:
		return ">"
	case CdLess:
		return "<"
	case CdGreaterOrEqual:
		return ">="
	case CdLessOrEqual:
		return "<="
	case CdIf:
		return "if"
	case CdThen:
		return "then"
	case CdElse:
		return "else"
	case StElementOf:
		return "∈"
	case StNotElementOf:
		return "∉"
	case StUnion:
		return "⋃"
	case StIntersection:
		return "⋂"
	case StSubset:
		return "⊆"
	case StSuperset:
		return "⊇"
	case LtTrue:
		return "T"
	case LtFalse:
		return "F"
	case LtNumber:
		return "LtNumber"
	case EOF:
		return "EOF"
	default:
		return fmt.Sprintf("TokenType(%d)", int(tt))
	}
}

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
	"(":  GlLeftParenthesis,
	")":  GlRightParenthesis,
	"{":  GlLeftBrace,
	"}":  GlRightBrace,
	",":  GlComma,
	";":  GlSemicolon,
	".":  GlDot,
	"|":  GlPipe,
	":=": GlAssignment,
	"==": CdEquals,
	"!=": CdNotEquals,
	"+":  ArAddition,
	"-":  ArSubtraction,
	"*":  ArMultiplication,
	"/":  ArDivision,
	"%":  ArModulus,
	"^":  ArPower,
	">":  CdGreater,
	"<":  CdLess,
	">=": CdGreaterOrEqual,
	"<=": CdLessOrEqual,
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
	"conjunction":    BlConjunction,
	"disjunction":    BlDisjunction,
	"negation":       BlNegation,
	"equivalence":    BlEquivalence,
	"implication":    BlImplication,
	"forall":         BlForall,
	"exists":         BlExists,
	"element_of":     StElementOf,
	"not_element_of": StNotElementOf,
	"union":          StUnion,
	"intersection":   StIntersection,
	"subset":         StSubset,
	"superset":       StSuperset,
	"if":             CdIf,
	"then":           CdThen,
	"else":           CdElse,
	"T":              LtTrue,
	"F":              LtFalse,
}

var (
	numberRegex     = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?`)
	identifierRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`)
)
