package lexer

import (
	"reflect"
	"testing"
)

func TestLexer_BasicTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token[BooleanTokenType]
	}{
		{
			name:  "parentheses",
			input: "()",
			expected: []Token[BooleanTokenType]{
				{Type: LPAREN, Value: "(", Pos: 0},
				{Type: RPAREN, Value: ")", Pos: 1},
			},
		},
		{
			name:  "quantifiers",
			input: "AE",
			expected: []Token[BooleanTokenType]{
				{Type: FORALL, Value: "A", Pos: 0},
				{Type: EXISTS, Value: "E", Pos: 1},
			},
		},
		{
			name:  "implication",
			input: "->",
			expected: []Token[BooleanTokenType]{
				{Type: IMPL, Value: "->", Pos: 0},
			},
		},
		{
			name:  "equivalence",
			input: "~",
			expected: []Token[BooleanTokenType]{
				{Type: EQUIV, Value: "~", Pos: 0},
			},
		},
		{
			name:  "conjunction_disjunction",
			input: "&\\/",
			expected: []Token[BooleanTokenType]{
				{Type: CONJ, Value: "&", Pos: 0},
				{Type: DISJ, Value: "\\/", Pos: 1},
			},
		},
		{
			name:  "negations",
			input: "!-",
			expected: []Token[BooleanTokenType]{
				{Type: NEG, Value: "!", Pos: 0},
				{Type: NEG, Value: "-", Pos: 1},
			},
		},
		{
			name:  "literals",
			input: "10",
			expected: []Token[BooleanTokenType]{
				{Type: LIT, Value: "1", Pos: 0},
				{Type: LIT, Value: "0", Pos: 1},
			},
		},
		{
			name:  "variables",
			input: "abc",
			expected: []Token[BooleanTokenType]{
				{Type: VAR, Value: "a", Pos: 0},
				{Type: VAR, Value: "b", Pos: 1},
				{Type: VAR, Value: "c", Pos: 2},
			},
		},
		{
			name:  "predicates",
			input: "BC",
			expected: []Token[BooleanTokenType]{
				{Type: PRED, Value: "B", Pos: 0},
				{Type: PRED, Value: "C", Pos: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := NewBooleanLexer(tt.input)
			tokens, err := lex.Lex()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !reflect.DeepEqual(tokens, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, tokens)
			}
		})
	}
}

func TestLexer_ComplexExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []BooleanTokenType
	}{
		{
			name:  "simple_implication",
			input: "a -> b",
			want:  []BooleanTokenType{VAR, IMPL, VAR},
		},
		{
			name:  "conjunction_with_negation",
			input: "!a & b",
			want:  []BooleanTokenType{NEG, VAR, CONJ, VAR},
		},
		{
			name:  "complex_expression",
			input: "(a \\/ b) -> (!c & d)",
			want: []BooleanTokenType{
				LPAREN, VAR, DISJ, VAR, RPAREN,
				IMPL,
				LPAREN, NEG, VAR, CONJ, VAR, RPAREN,
			},
		},
		{
			name:  "equivalence_expression",
			input: "a ~ (b & c)",
			want: []BooleanTokenType{
				VAR, EQUIV,
				LPAREN, VAR, CONJ, VAR, RPAREN,
			},
		},
		{
			name:  "quantifier_expression",
			input: "A(x) P(x)",
			want: []BooleanTokenType{
				FORALL, LPAREN, VAR, RPAREN,
				PRED, LPAREN, VAR, RPAREN,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := NewBooleanLexer(tt.input)
			tokens, err := lex.Lex()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(tokens) != len(tt.want) {
				t.Fatalf("Expected %d tokens, got %d", len(tt.want), len(tokens))
			}

			for i, expectedType := range tt.want {
				if tokens[i].Type != expectedType {
					t.Errorf("Token %d: expected %s, got %s", i, expectedType, tokens[i].Type)
				}
			}
		})
	}
}

func TestLexer_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "invalid_character",
			input:       "a @ b",
			expectError: true,
		},
		{
			name:        "incomplete_disjunction",
			input:       "a \\ b",
			expectError: true,
		},
		{
			name:        "numbers_other_than_01",
			input:       "a & 2",
			expectError: true,
		},
		{
			name:        "valid_expression",
			input:       "a -> b",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := NewBooleanLexer(tt.input)
			_, err := lex.Lex()
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestLexer_WhitespaceHandling(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []BooleanTokenType
	}{
		{
			name:  "spaces_between_tokens",
			input: "  a   &   b  ",
			want:  []BooleanTokenType{VAR, CONJ, VAR},
		},
		{
			name:  "tabs_and_newlines",
			input: "a\t->\nb",
			want:  []BooleanTokenType{VAR, IMPL, VAR},
		},
		{
			name:  "no_whitespace",
			input: "a&b",
			want:  []BooleanTokenType{VAR, CONJ, VAR},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := NewBooleanLexer(tt.input)
			tokens, err := lex.Lex()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(tokens) != len(tt.want) {
				t.Fatalf("Expected %d tokens, got %d", len(tt.want), len(tokens))
			}

			for i, expectedType := range tt.want {
				if tokens[i].Type != expectedType {
					t.Errorf("Token %d: expected %s, got %s", i, expectedType, tokens[i].Type)
				}
			}
		})
	}
}
