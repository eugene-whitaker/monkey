package ast

import (
	"monkey/token"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		program  *Program
		expected string
	}{
		{
			&Program{
				Statements: []Statement{
					&LetStatement{
						Token: token.Token{Type: token.LET, Lexeme: "let"},
						Name: &Identifier{
							Token: token.Token{Type: token.IDENT, Lexeme: "ident"},
							Value: "ident",
						},
						Value: &Identifier{
							Token: token.Token{Type: token.IDENT, Lexeme: "value"},
							Value: "value",
						},
					},
				},
			},
			"let ident = value;",
		},
	}

	for _, tt := range tests {
		actual := tt.program.String()

		if tt.expected != actual {
			t.Errorf("program.String() ==> expected: %q actual: %q", tt.expected, actual)
		}
	}
}
