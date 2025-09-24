package parser

import "logicka/lib/lexer"

type Precedence int

const (
	PrecedenceDefault Precedence = iota
	PrecedenceAssignment
	PrecedenceImplication
	PrecedenceEquivalence
	PrecedenceDisjunction
	PrecedenceConjunction
	PrecedenceComparison
	PrecedenceAddition
	PrecedenceMultiplication
	PrecedenceExponential
	PrecedencePrefix
	PrecedencePrimary
)

type precedenceMap map[lexer.TokenType]Precedence

var tokenPrecedences = precedenceMap{
	lexer.GlAssignment:      PrecedenceAssignment,
	lexer.BlImplication:     PrecedenceImplication,
	lexer.BlEquivalence:     PrecedenceEquivalence,
	lexer.BlDisjunction:     PrecedenceDisjunction,
	lexer.BlConjunction:     PrecedenceConjunction,
	lexer.CdEquals:          PrecedenceComparison,
	lexer.CdNotEquals:       PrecedenceComparison,
	lexer.CdGreater:         PrecedenceComparison,
	lexer.CdLess:            PrecedenceComparison,
	lexer.CdGreaterOrEqual:  PrecedenceComparison,
	lexer.CdLessOrEqual:     PrecedenceComparison,
	lexer.StElementOf:       PrecedenceComparison,
	lexer.StNotElementOf:    PrecedenceComparison,
	lexer.StSubset:          PrecedenceComparison,
	lexer.StSuperset:        PrecedenceComparison,
	lexer.ArAddition:        PrecedenceAddition,
	lexer.ArSubtraction:     PrecedenceAddition,
	lexer.StUnion:           PrecedenceAddition,
	lexer.ArMultiplication:  PrecedenceMultiplication,
	lexer.ArDivision:        PrecedenceMultiplication,
	lexer.ArModulus:         PrecedenceMultiplication,
	lexer.StIntersection:    PrecedenceMultiplication,
	lexer.ArPower:           PrecedenceExponential,
	lexer.BlNegation:        PrecedencePrefix,
	lexer.BlForall:          PrecedencePrefix,
	lexer.BlExists:          PrecedencePrefix,
	lexer.GlLeftParenthesis: PrecedencePrimary,
	lexer.GlLeftBrace:       PrecedencePrimary,
	lexer.GlIdentifier:      PrecedencePrimary,
	lexer.LtNumber:          PrecedencePrimary,
	lexer.LtTrue:            PrecedencePrimary,
	lexer.LtFalse:           PrecedencePrimary,
	lexer.GlQuote:           PrecedencePrimary,
}

func (pm precedenceMap) Get(tokenType lexer.TokenType) Precedence {
	if prec, exists := pm[tokenType]; exists {
		return prec
	}
	return 0
}
