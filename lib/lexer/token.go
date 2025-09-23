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
		return "GlLeftParenthesis"
	case GlRightParenthesis:
		return "GlRightParenthesis"
	case GlLeftBrace:
		return "GlLeftBrace"
	case GlRightBrace:
		return "GlRightBrace"
	case GlComma:
		return "GlComma"
	case GlSemicolon:
		return "GlSemicolon"
	case GlDot:
		return "GlDot"
	case GlPipe:
		return "GlPipe"
	case GlAssignment:
		return "GlAssignment"
	case GlIdentifier:
		return "GlIdentifier"
	case BlConjunction:
		return "BlConjunction"
	case BlDisjunction:
		return "BlDisjunction"
	case BlNegation:
		return "BlNegation"
	case BlEquivalence:
		return "BlEquivalence"
	case BlImplication:
		return "BlImplication"
	case BlForall:
		return "BlForall"
	case BlExists:
		return "BlExists"
	case ArAddition:
		return "ArAddition"
	case ArSubtraction:
		return "ArSubtraction"
	case ArMultiplication:
		return "ArMultiplication"
	case ArDivision:
		return "ArDivision"
	case ArModulus:
		return "ArModulus"
	case ArPower:
		return "ArPower"
	case CdEquals:
		return "CdEquals"
	case CdNotEquals:
		return "CdNotEquals"
	case CdGreater:
		return "CdGreater"
	case CdLess:
		return "CdLess"
	case CdGreaterOrEqual:
		return "CdGreaterOrEqual"
	case CdLessOrEqual:
		return "CdLessOrEqual"
	case CdIf:
		return "CdIf"
	case CdThen:
		return "CdThen"
	case CdElse:
		return "CdElse"
	case StElementOf:
		return "StElementOf"
	case StNotElementOf:
		return "StNotElementOf"
	case StUnion:
		return "StUnion"
	case StIntersection:
		return "StIntersection"
	case StSubset:
		return "StSubset"
	case StSuperset:
		return "StSuperset"
	case LtTrue:
		return "LtTrue"
	case LtFalse:
		return "LtFalse"
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
