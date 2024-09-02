package lexer

import "monkey/token"

type Lexer struct {
	input   string
	start   int
	current int

	ch byte
}

func New(input string) *Lexer {
	l := &Lexer{
		input:   input,
		start:   0,
		current: 0,
	}
	l.ch = l.peek()
	return l
}

func (l *Lexer) NextToken() token.Token {
	l.advance()

	for l.ch != 0 {
		switch l.ch {
		case '=':
			return l.emitTwoCharToken('=', token.EQ, token.ASSIGN)
		case '+':
			return l.emitOneCharToken(token.PLUS)
		case '-':
			return l.emitOneCharToken(token.MINUS)
		case '!':
			return l.emitTwoCharToken('=', token.NOT_EQ, token.BANG)
		case '/':
			return l.emitOneCharToken(token.SLASH)
		case '*':
			return l.emitOneCharToken(token.ASTERISK)
		case '<':
			return l.emitOneCharToken(token.LT)
		case '>':
			return l.emitOneCharToken(token.GT)
		case ';':
			return l.emitOneCharToken(token.SEMICOLON)
		case ',':
			return l.emitOneCharToken(token.COMMA)
		case '(':
			return l.emitOneCharToken(token.LPAREN)
		case ')':
			return l.emitOneCharToken(token.RPAREN)
		case '{':
			return l.emitOneCharToken(token.LBRACE)
		case '}':
			return l.emitOneCharToken(token.RBRACE)
		case ' ', '\t', '\r':
		case '\n':
		default:
			if isLetter(l.ch) {
				literal := l.lexeme(isLetter)
				return token.Token{
					Type:    token.LookupKeyword(literal),
					Literal: literal,
				}
			} else if isDigit(l.ch) {
				return token.Token{Type: token.INT, Literal: l.lexeme(isDigit)}
			} else {
				return l.emitOneCharToken(token.ILLEGAL)
			}
		}

		l.advance()
	}

	return token.Token{Type: token.EOF, Literal: ""}
}

func (l *Lexer) advance() {
	l.start = l.current
	if l.current >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.current]
		l.current = l.current + 1
	}
}

func (l *Lexer) consume() {
	if l.current >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.current]
		l.current = l.current + 1
	}
}

func (l *Lexer) peek() byte {
	if l.current >= len(l.input) {
		return 0
	}
	return l.input[l.current]
}

func (l *Lexer) match(expected byte) bool {
	if l.current >= len(l.input) {
		return false
	}

	if l.input[l.current] != expected {
		return false
	}

	l.current = l.current + 1
	return true
}

func (l *Lexer) lexeme(condition func(byte) bool) string {
	for condition(l.peek()) {
		l.consume()
	}
	return l.input[l.start:l.current]
}

func (l *Lexer) emitOneCharToken(tokenType token.TokenType) token.Token {
	return token.Token{Type: tokenType, Literal: string(l.ch)}
}

func (l *Lexer) emitTwoCharToken(match byte, matched, mismatched token.TokenType) token.Token {
	if l.match(match) {
		return token.Token{Type: matched, Literal: l.input[l.start:l.current]}
	} else {
		return l.emitOneCharToken(mismatched)
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
