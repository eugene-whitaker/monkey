package lexer

import (
	"testing"

	"github.com/eugene-whitaker/monkey/token"
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
			"let string = \"hello world\";",
			[]TokenTest{
				{token.LET, "let"},
				{token.IDENT, "string"},
				{token.ASSIGN, "="},
				{token.STRING, "hello world"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			"let array = [1, 2];",
			[]TokenTest{
				{token.LET, "let"},
				{token.IDENT, "array"},
				{token.ASSIGN, "="},
				{token.LBRACKET, "["},
				{token.INT, "1"},
				{token.COMMA, ","},
				{token.INT, "2"},
				{token.RBRACKET, "]"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			"let hash = {\"key\": \"value\"};",
			[]TokenTest{
				{token.LET, "let"},
				{token.IDENT, "hash"},
				{token.ASSIGN, "="},
				{token.LBRACE, "{"},
				{token.STRING, "key"},
				{token.COLON, ":"},
				{token.STRING, "value"},
				{token.RBRACE, "}"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			"macro(x, y) { x + y; };",
			[]TokenTest{
				{token.MACRO, "macro"},
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
	}

	for i, test := range tests {
		if !testTokens(t, i, test.input, test.tests) {
			continue
		}
	}
}

func testTokens(t *testing.T, idx int, input string, tests []TokenTest) bool {
	l := NewLexer(input)
	for j, expected := range tests {
		actual := l.NextToken()

		if expected.ttype != actual.Type {
			t.Errorf("tests[%d][%d] - %q ==> expected: %q actual: %q", idx, j, input, expected.ttype, actual.Type)
			return false
		}

		if expected.lexeme != actual.Lexeme {
			t.Errorf("tests[%d][%d] - %q ==> expected: %q actual: %q", idx, j, input, expected.lexeme, actual.Lexeme)
			return false
		}
	}

	return true
}
