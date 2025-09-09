package parser

import (
	"fmt"
	"logicka/lib/lexer"
	"testing"
)

func TestParser_BasicExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string // String representation for verification
	}{
		{
			name:     "simple_variable",
			input:    "a",
			expected: "a",
		},
		{
			name:     "literal_true",
			input:    "1",
			expected: "true",
		},
		{
			name:     "literal_false",
			input:    "0",
			expected: "false",
		},
		{
			name:     "simple_conjunction",
			input:    "a & b",
			expected: "a & b",
		},
		{
			name:     "simple_disjunction",
			input:    "a \\/ b",
			expected: "a \\/ b",
		},
		{
			name:     "simple_implication",
			input:    "a -> b",
			expected: "a -> b",
		},
		{
			name:     "simple_equivalence",
			input:    "a ~ b",
			expected: "a ~ b",
		},
		{
			name:     "simple_negation",
			input:    "!a",
			expected: "!a",
		},
		{
			name:     "negation_with_dash",
			input:    "-a",
			expected: "!a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := lexer.NewLexer(tt.input)
			tokens, err := lex.Lex()
			if err != nil {
				t.Fatalf("Lexer error: %v", err)
			}

			p := &Parser{Tokens: tokens}
			ast, err := p.ParseExpression()
			if err != nil {
				t.Fatalf("Parser error: %v", err)
			}

			if ast.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, ast.String())
			}
		})
	}
}

func TestParser_Precedence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "conjunction_over_disjunction",
			input:    "a \\/ b & c",
			expected: "a \\/ b & c",
		},
		{
			name:     "negation_highest",
			input:    "!a & b",
			expected: "!a & b",
		},
		{
			name:     "implication_lowest",
			input:    "a & b -> c",
			expected: "a & b -> c",
		},
		{
			name:     "equivalence_lowest",
			input:    "a -> b ~ c -> d",
			expected: "a -> b ~ c -> d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := lexer.NewLexer(tt.input)
			tokens, err := lex.Lex()
			if err != nil {
				t.Fatalf("Lexer error: %v", err)
			}

			p := &Parser{Tokens: tokens}
			ast, err := p.ParseExpression()
			if err != nil {
				t.Fatalf("Parser error: %v", err)
			}

			result := ast.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestParser_Grouping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple_grouping",
			input:    "(a)",
			expected: "(a)",
		},
		{
			name:     "grouping_changes_precedence",
			input:    "(a \\/ b) & c",
			expected: "(a \\/ b) & c",
		},
		{
			name:     "nested_grouping",
			input:    "((a & b) \\/ c)",
			expected: "((a & b) \\/ c)",
		},
		{
			name:     "multiple_groups",
			input:    "(a & b) \\/ (c & d)",
			expected: "(a & b) \\/ (c & d)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := lexer.NewLexer(tt.input)
			tokens, err := lex.Lex()
			if err != nil {
				t.Fatalf("Lexer error: %v", err)
			}

			p := &Parser{Tokens: tokens}
			ast, err := p.ParseExpression()
			if err != nil {
				t.Fatalf("Parser error: %v", err)
			}

			result := ast.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestParser_ComplexExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "complex_nested",
			input: "((a & b) \\/ (!c -> d)) ~ (e & f)",
		},
		{
			name:  "multiple_negations",
			input: "!!a & !b",
		},
		{
			name:  "chained_implications",
			input: "a -> b -> c",
		},
		{
			name:  "mixed_operators",
			input: "!a \\/ (b & c) -> (d ~ e)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := lexer.NewLexer(tt.input)
			tokens, err := lex.Lex()
			if err != nil {
				t.Fatalf("Lexer error: %v", err)
			}

			p := &Parser{Tokens: tokens}
			ast, err := p.ParseExpression()
			if err != nil {
				t.Fatalf("Parser error: %v", err)
			}

			// Verify that parsing completed successfully
			if ast == nil {
				t.Error("Expected AST but got nil")
			}
		})
	}
}

func TestParser_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "unmatched_left_paren",
			input:       "(a & b",
			expectError: true,
		},
		{
			name:        "unmatched_right_paren",
			input:       "a & b)",
			expectError: true,
		},
		{
			name:        "missing_operand",
			input:       "a &",
			expectError: true,
		},
		{
			name:        "empty_expression",
			input:       "",
			expectError: true,
		},
		{
			name:        "operator_without_operands",
			input:       "&",
			expectError: true,
		},
		{
			name:        "valid_expression",
			input:       "a & b",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := lexer.NewLexer(tt.input)
			tokens, err := lex.Lex()
			if err != nil && !tt.expectError {
				t.Fatalf("Unexpected lexer error: %v", err)
			}
			if err != nil && tt.expectError {
				return // Expected error at lexer level
			}

			p := &Parser{Tokens: tokens}
			te, err := p.ParseExpression()
			fmt.Println(te)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestParser_PredicateAndQuantifier(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple_predicate",
			input:    "P(x)",
			expected: "P", // Based on current PredicateNode.String() implementation
		},
		{
			name:     "universal_quantifier",
			input:    "A(x)",
			expected: "x", // Based on current QuantifierNode.String() implementation
		},
		{
			name:     "existential_quantifier",
			input:    "E(y)",
			expected: "y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := lexer.NewLexer(tt.input)
			tokens, err := lex.Lex()
			if err != nil {
				t.Fatalf("Lexer error: %v", err)
			}

			p := &Parser{Tokens: tokens}
			ast, err := p.ParseExpression()
			if err != nil {
				t.Fatalf("Parser error: %v", err)
			}

			result := ast.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
