package lexer

import (
	"testing"

	"monkey/token"
)

type TokenTest struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func TestNextToken(t *testing.T) {
	tests := []struct {
		input    string
		expected []TokenTest
	}{
		{
			"let five = 5;",
			[]TokenTest{
				{token.LET, "let"},
				{token.IDENT, "five"},
				{token.ASSIGN, "="},
				{token.INT, "5"},
				{token.SEMICOLON, ";"},
			},
		},
		{
			"let ten = 10;",
			[]TokenTest{
				{token.LET, "let"},
				{token.IDENT, "ten"},
				{token.ASSIGN, "="},
				{token.INT, "10"},
				{token.SEMICOLON, ";"},
			},
		},
		{
			"let add = fn(x, y) { x + y; };",
			[]TokenTest{
				{token.LET, "let"},
				{token.IDENT, "add"},
				{token.ASSIGN, "="},
				{token.FUNCTION, "fn"},
				{token.LPAREN, "("},
				{token.IDENT, "x"},
				{token.COMMA, ","},
				{token.IDENT, "y"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.IDENT, "x"},
				{token.PLUS, "+"},
				{token.IDENT, "y"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.SEMICOLON, ";"},
			},
		},
		{
			"let result = add(five, ten);",
			[]TokenTest{
				{token.LET, "let"},
				{token.IDENT, "result"},
				{token.ASSIGN, "="},
				{token.IDENT, "add"},
				{token.LPAREN, "("},
				{token.IDENT, "five"},
				{token.COMMA, ","},
				{token.IDENT, "ten"},
				{token.RPAREN, ")"},
				{token.SEMICOLON, ";"},
			},
		},
		{
			"!-/*5;",
			[]TokenTest{
				{token.BANG, "!"},
				{token.MINUS, "-"},
				{token.SLASH, "/"},
				{token.ASTERISK, "*"},
				{token.INT, "5"},
				{token.SEMICOLON, ";"},
			},
		},
		{
			"5 < 10 > 5;",
			[]TokenTest{
				{token.INT, "5"},
				{token.LT, "<"},
				{token.INT, "10"},
				{token.GT, ">"},
				{token.INT, "5"},
				{token.SEMICOLON, ";"},
			},
		},
		{
			"if (5 < 10) { return true; } else { return false; };",
			[]TokenTest{
				{token.IF, "if"},
				{token.LPAREN, "("},
				{token.INT, "5"},
				{token.LT, "<"},
				{token.INT, "10"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.RETURN, "return"},
				{token.TRUE, "true"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.ELSE, "else"},
				{token.LBRACE, "{"},
				{token.RETURN, "return"},
				{token.FALSE, "false"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.SEMICOLON, ";"},
			},
		},
		{
			"10 == 10;",
			[]TokenTest{
				{token.INT, "10"},
				{token.EQ, "=="},
				{token.INT, "10"},
				{token.SEMICOLON, ";"},
			},
		},
		{
			"10 != 9;",
			[]TokenTest{
				{token.INT, "10"},
				{token.NOT_EQ, "!="},
				{token.INT, "9"},
				{token.SEMICOLON, ";"},
			},
		},
	}

	for i, tt := range tests {
		l := New(tt.input)

		for j, expected := range tt.expected {
			actual := l.NextToken()

			if expected.expectedType != actual.Type {
				t.Fatalf("tests[%d][%d] ==> expected: %q actual: %q", i, j, expected.expectedType, actual.Type)
			}

			if expected.expectedLiteral != actual.Literal {
				t.Fatalf("tests[%d][%d] ==> expected: %q actual: %q", i, j, expected.expectedType, actual.Type)
			}
		}
	}
}
