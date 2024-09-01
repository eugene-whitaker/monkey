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
						Token: token.Token{Type: token.LET, Literal: "let"},
						Name: &Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Value: "x",
						},
						Value: &Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "y"},
							Value: "y",
						},
					},
				},
			},
			"let x = y;",
		},
	}

	for _, tt := range tests {
		actual := tt.program.String()
		if tt.expected != actual {
			t.Errorf("program.String() ==> expected: %q actual: %q", tt.expected, actual)
		}
	}
}
