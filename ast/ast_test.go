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

	for i, test := range tests {
		actual := test.program.String()

		if test.expected != actual {
			t.Errorf("test[%d] - program.String() ==> expected: %q actual: %q", i, test.expected, actual)
		}
	}
}
