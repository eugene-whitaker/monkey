package lexer

import "monkey/token"

type Lexer struct {
	input string
	ch    byte

	start   int
	current int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:   input,
		ch:      0,
		start:   0,
		current: 0,
	}
}

func (l *Lexer) NextToken() token.Token {
	l.advance()

	for l.ch != 0 {
		switch l.ch {
		case '=':
			if l.match('=') {
				return token.Token{
					Type:   token.EQ,
					Lexeme: l.input[l.start:l.current],
					Offset: l.start,
					Length: l.current - l.start,
				}
			} else {
				return l.emit(token.ASSIGN)
			}
		case '+':
			return l.emit(token.PLUS)
		case '-':
			return l.emit(token.MINUS)
		case '!':
			if l.match('=') {
				return token.Token{
					Type:   token.NOT_EQ,
					Lexeme: l.input[l.start:l.current],
					Offset: l.start,
					Length: l.current - l.start,
				}
			} else {
				return l.emit(token.BANG)
			}
		case '/':
			return l.emit(token.SLASH)
		case '*':
			return l.emit(token.ASTERISK)
		case '<':
			return l.emit(token.LT)
		case '>':
			return l.emit(token.GT)
		case ';':
			return l.emit(token.SEMICOLON)
		case ',':
			return l.emit(token.COMMA)
		case '(':
			return l.emit(token.LPAREN)
		case ')':
			return l.emit(token.RPAREN)
		case '{':
			return l.emit(token.LBRACE)
		case '}':
			return l.emit(token.RBRACE)
		case ' ', '\t', '\r', '\n':
		case '"':
			return l.string()
		default:
			if isLetter(l.ch) {
				literal := l.identifier()
				return token.Token{
					Type:   token.LookupKeyword(literal),
					Lexeme: literal,
					Offset: l.start,
					Length: l.current - l.start,
				}
			} else if isDigit(l.ch) {
				literal := l.number()
				return token.Token{
					Type:   token.INT,
					Lexeme: literal,
					Offset: l.start,
					Length: l.current - l.start,
				}
			} else {
				return l.emit(token.ILLEGAL)
			}
		}

		l.advance()
	}

	return token.Token{
		Type:   token.EOF,
		Lexeme: "",
		Offset: l.start,
		Length: 0,
	}
}

func (l *Lexer) Input() string {
	return l.input
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

func (l *Lexer) identifier() string {
	for isLetter(l.peek()) {
		l.consume()
	}
	return l.input[l.start:l.current]
}

func (l *Lexer) number() string {
	for isDigit(l.peek()) {
		l.consume()
	}
	return l.input[l.start:l.current]
}

func (l *Lexer) string() token.Token {
	for l.peek() != '"' && l.peek() != 0 {
		l.consume()
	}

	if l.peek() != '"' {
		return token.Token{
			Type:   token.ILLEGAL,
			Lexeme: l.input[l.start:l.current],
			Offset: l.start,
			Length: l.start - l.current,
		}
	}

	l.consume()

	return token.Token{
		Type:   token.STRING,
		Lexeme: l.input[l.start+1 : l.current-1],
		Offset: l.start,
		Length: l.start - l.current,
	}
}

func (l *Lexer) emit(ttype token.TokenType) token.Token {
	return token.Token{
		Type:   ttype,
		Lexeme: string(l.ch),
		Offset: l.start,
		Length: 1,
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
