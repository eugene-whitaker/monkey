package lexer

import (
	"monkey/token"
	"testing"
)

type TokenTest struct {
	ttype  token.TokenType
	lexeme string
}

func TestNextToken(t *testing.T) {
	tests := []struct {
		input string
		tests []TokenTest
	}{
		{
			"let five = 5;",
			[]TokenTest{
				{token.LET, "let"},
				{token.IDENT, "five"},
				{token.ASSIGN, "="},
				{token.INT, "5"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
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
				{token.EOF, ""},
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
				{token.EOF, ""},
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
				{token.EOF, ""},
			},
		},
		{
			"let multiply = fn(x, y) { x * y; };",
			[]TokenTest{
				{token.LET, "let"},
				{token.IDENT, "multiply"},
				{token.ASSIGN, "="},
				{token.FUNCTION, "fn"},
				{token.LPAREN, "("},
				{token.IDENT, "x"},
				{token.COMMA, ","},
				{token.IDENT, "y"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.IDENT, "x"},
				{token.ASTERISK, "*"},
				{token.IDENT, "y"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			"let subtract = fn(x, y) { x - y; };",
			[]TokenTest{
				{token.LET, "let"},
				{token.IDENT, "subtract"},
				{token.ASSIGN, "="},
				{token.FUNCTION, "fn"},
				{token.LPAREN, "("},
				{token.IDENT, "x"},
				{token.COMMA, ","},
				{token.IDENT, "y"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.IDENT, "x"},
				{token.MINUS, "-"},
				{token.IDENT, "y"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			"let divide = fn(x, y) { x / y; };",
			[]TokenTest{
				{token.LET, "let"},
				{token.IDENT, "divide"},
				{token.ASSIGN, "="},
				{token.FUNCTION, "fn"},
				{token.LPAREN, "("},
				{token.IDENT, "x"},
				{token.COMMA, ","},
				{token.IDENT, "y"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.IDENT, "x"},
				{token.SLASH, "/"},
				{token.IDENT, "y"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			"if (x > y) { return true; } else { return false; };",
			[]TokenTest{
				{token.IF, "if"},
				{token.LPAREN, "("},
				{token.IDENT, "x"},
				{token.GT, ">"},
				{token.IDENT, "y"},
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
				{token.EOF, ""},
			},
		},
		{
			"if (x < y) { return true; } else { return false; };",
			[]TokenTest{
				{token.IF, "if"},
				{token.LPAREN, "("},
				{token.IDENT, "x"},
				{token.LT, "<"},
				{token.IDENT, "y"},
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
				{token.EOF, ""},
			},
		},
		{
			"if (x == y) { return true; } else { return false; };",
			[]TokenTest{
				{token.IF, "if"},
				{token.LPAREN, "("},
				{token.IDENT, "x"},
				{token.EQ, "=="},
				{token.IDENT, "y"},
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
				{token.EOF, ""},
			},
		},
		{
			"if (x != y) { return true; } else { return false; };",
			[]TokenTest{
				{token.IF, "if"},
				{token.LPAREN, "("},
				{token.IDENT, "x"},
				{token.NOT_EQ, "!="},
				{token.IDENT, "y"},
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
				{token.EOF, ""},
			},
		},
		{
			"let string = \"chars\";",
			[]TokenTest{
				{token.LET, "let"},
				{token.IDENT, "string"},
				{token.ASSIGN, "="},
				{token.STRING, "chars"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			"let string = \"s p a c e\";",
			[]TokenTest{
				{token.LET, "let"},
				{token.IDENT, "string"},
				{token.ASSIGN, "="},
				{token.STRING, "s p a c e"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
	}

	for i, tt := range tests {
		l := NewLexer(tt.input)

		for j, expected := range tt.tests {
			actual := l.NextToken()

			if expected.ttype != actual.Type {
				t.Fatalf("tests[%d][%d] - %q ==> expected: %q actual: %q", i, j, tt.input, expected.ttype, actual.Type)
			}

			if expected.lexeme != actual.Lexeme {
				t.Fatalf("tests[%d][%d] - %q ==> expected: %q actual: %q", i, j, tt.input, expected.lexeme, actual.Lexeme)
			}
		}
	}
}
