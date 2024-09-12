package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // X == X
	LESSGREATER // X > X or X < X
	SUM         // X + X
	PRODUCT     // X * X
	PREFIX      // -X or !X
	CALL        // func(X)
	INDEX       // arr[X]
)

type (
	prefixFunc func() ast.Expression
	infixFunc  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	tokens []token.Token
	tok    token.Token

	current int

	prefixFuncs map[token.TokenType]prefixFunc
	infixFuncs  map[token.TokenType]infixFunc

	errors []string
}

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}

	p.tokens = []token.Token{}
	t := l.NextToken()
	for t.Type != token.EOF {
		p.tokens = append(p.tokens, t)
		t = l.NextToken()
	}
	p.tokens = append(p.tokens, t)
	p.tok = token.Token{}

	p.prefixFuncs = map[token.TokenType]prefixFunc{
		token.IDENT:    p.parseIdentifier,
		token.INT:      p.parseIntegerLiteral,
		token.BANG:     p.parsePrefixExpression,
		token.MINUS:    p.parsePrefixExpression,
		token.TRUE:     p.parseBooleanLiteral,
		token.FALSE:    p.parseBooleanLiteral,
		token.LPAREN:   p.parseGroupedExpression,
		token.IF:       p.parseIfExpression,
		token.FUNCTION: p.parseFunctionLiteral,
		token.STRING:   p.parseStringLiteral,
		token.LBRACKET: p.parseArrayLiteral,
		token.LBRACE:   p.parseHashLiteral,
		token.MACRO:    p.parseMacroExpression,
	}

	p.infixFuncs = map[token.TokenType]infixFunc{
		token.PLUS:     p.parseInfixExpression,
		token.MINUS:    p.parseInfixExpression,
		token.SLASH:    p.parseInfixExpression,
		token.ASTERISK: p.parseInfixExpression,
		token.EQ:       p.parseInfixExpression,
		token.NOT_EQ:   p.parseInfixExpression,
		token.LT:       p.parseInfixExpression,
		token.GT:       p.parseInfixExpression,
		token.LPAREN:   p.parseCallExpression,
		token.LBRACKET: p.parseIndexExpression,
	}

	p.errors = []string{}

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	// defer untrace(trace("ParseProgram"))
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	p.advance()

	for p.tok.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.advance()
	}

	return program
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) parseStatement() ast.Statement {
	// defer untrace(trace("parseStatement"))
	switch p.tok.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	// defer untrace(trace("parseLetStatement"))
	stmt := &ast.LetStatement{
		Token: p.tok,
	}

	if !p.expect(token.IDENT, "expected <IDENT> token following <let>") {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.tok,
		Value: p.tok.Lexeme,
	}

	if !p.expect(token.ASSIGN, "expected <=> token following <IDENT>") {
		return nil
	}

	p.advance()

	stmt.Value = p.parseExpression(LOWEST)

	if p.check(token.SEMICOLON) {
		p.advance()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	// defer untrace(trace("parseReturnStatement"))
	stmt := &ast.ReturnStatement{
		Token: p.tok,
	}

	p.advance()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.check(token.SEMICOLON) {
		p.advance()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{
		Token: p.tok,
	}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.check(token.SEMICOLON) {
		p.advance()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// defer untrace(trace("parseBlockStatement"))
	block := &ast.BlockStatement{
		Token:      p.tok,
		Statements: []ast.Statement{},
	}

	p.advance()

	for p.tok.Type != token.RBRACE && p.tok.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.advance()
	}

	return block
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// defer untrace(trace("parseExpression"))
	prefix := p.prefixFuncs[p.tok.Type]
	if prefix == nil {
		p.error(p.tok, fmt.Sprintf("no prefix parse function for <%s>", p.tok.Type))
		return nil
	}
	leftExpr := prefix()

	for precedence < p.precedence(p.peek().Type) {
		infix := p.infixFuncs[p.peek().Type]
		if infix == nil {
			return leftExpr
		}

		p.advance()

		leftExpr = infix(leftExpr)
	}

	return leftExpr
}

func (p *Parser) parseIdentifier() ast.Expression {
	// defer untrace(trace("parseIdentifier"))
	return &ast.Identifier{
		Token: p.tok,
		Value: p.tok.Lexeme,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	// defer untrace(trace("parseIntegerLiteral"))
	lit := &ast.IntegerLiteral{
		Token: p.tok,
	}

	value, err := strconv.ParseInt(p.tok.Lexeme, 0, 64)
	if err != nil {
		p.error(p.tok, fmt.Sprintf("ParseInt returned error: %s", err.Error()))
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	// defer untrace(trace("parseBooleanLiteral"))
	return &ast.BooleanLiteral{
		Token: p.tok,
		Value: p.tok.Type == token.TRUE,
	}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	// defer untrace(trace("parseFunctionLiteral"))
	lit := &ast.FunctionLiteral{
		Token: p.tok,
	}

	if !p.expect(token.LPAREN, "expected <(> token following <fn>") {
		return nil
	}

	idents := []*ast.Identifier{}

	p.advance()

	if p.tok.Type != token.RPAREN {
		idents = append(
			idents,
			&ast.Identifier{
				Token: p.tok,
				Value: p.tok.Lexeme,
			},
		)

		for p.check(token.COMMA) {
			p.advance()
			p.advance()

			idents = append(
				idents,
				&ast.Identifier{
					Token: p.tok,
					Value: p.tok.Lexeme,
				},
			)
		}

		if !p.expect(token.RPAREN, "expected <)> token following function parameters") {
			return nil
		}
	}

	lit.Parameters = idents

	if !p.expect(token.LBRACE, "expected <{> token following <)>") {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	// defer untrace(trace("parseStringLiteral"))
	return &ast.StringLiteral{
		Token: p.tok,
		Value: p.tok.Lexeme,
	}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	// defer untrace(trace("parseArrayLiteral"))
	lit := &ast.ArrayLiteral{
		Token: p.tok,
	}

	elems := []ast.Expression{}

	p.advance()

	if p.tok.Type != token.RBRACKET {
		elems = append(elems, p.parseExpression(LOWEST))

		for p.check(token.COMMA) {
			p.advance()
			p.advance()

			elems = append(elems, p.parseExpression(LOWEST))
		}

		if !p.expect(token.RBRACKET, "expected <]> token following array elements") {
			return nil
		}
	}

	lit.Elements = elems

	return lit

}

func (p *Parser) parseHashLiteral() ast.Expression {
	// defer untrace(trace("parseArrayLiteral"))
	lit := &ast.HashLiteral{
		Token: p.tok,
	}

	pairs := make(map[ast.Expression]ast.Expression)

	p.advance()

	if p.tok.Type != token.RBRACE {
		key := p.parseExpression(LOWEST)

		if !p.expect(token.COLON, "expected <:> token following hash key") {
			return nil
		}

		p.advance()

		value := p.parseExpression(LOWEST)

		pairs[key] = value

		for p.check(token.COMMA) {
			p.advance()
			p.advance()

			key := p.parseExpression(LOWEST)

			if !p.expect(token.COLON, "expected <:> token following hash key") {
				return nil
			}

			p.advance()

			value := p.parseExpression(LOWEST)

			pairs[key] = value
		}

		if !p.expect(token.RBRACE, "expected <}> token following hash pairs") {
			return nil
		}
	}

	lit.Pairs = pairs

	return lit

}

func (p *Parser) parsePrefixExpression() ast.Expression {
	// defer untrace(trace("parsePrefixExpression"))
	expr := &ast.PrefixExpression{
		Token:    p.tok,
		Operator: p.tok.Lexeme,
	}

	p.advance()

	expr.Right = p.parseExpression(PREFIX)

	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// defer untrace(trace("parseInfixExpression"))
	expr := &ast.InfixExpression{
		Token:    p.tok,
		Operator: p.tok.Lexeme,
		Left:     left,
	}

	precedence := p.precedence(p.tok.Type)
	p.advance()
	expr.Right = p.parseExpression(precedence)
	//                             ^^^ decrement here for right-associativity

	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	// defer untrace(trace("parseIfExpression"))
	expr := &ast.IfExpression{
		Token: p.tok,
	}

	if !p.expect(token.LPAREN, "expected <(> token following <if>") {
		return nil
	}

	p.advance()

	expr.Condition = p.parseExpression(LOWEST)

	if !p.expect(token.RPAREN, "expected <)> token following if condition") {
		return nil
	}

	if !p.expect(token.LBRACE, "expected '{' token following ')'") {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()
	if p.check(token.ELSE) {
		p.advance()

		if !p.expect(token.LBRACE, "expected <{> token following <else>") {
			return nil
		}

		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	// defer untrace(trace("parseCallExpression"))
	expr := &ast.CallExpression{
		Token:    p.tok,
		Function: left,
	}

	args := []ast.Expression{}

	p.advance()
	if p.tok.Type != token.RPAREN {
		args = append(args, p.parseExpression(LOWEST))

		for p.check(token.COMMA) {
			p.advance()
			p.advance()

			args = append(args, p.parseExpression(LOWEST))
		}

		if !p.expect(token.RPAREN, "expected <)> token following call arguments") {
			return nil
		}
	}

	expr.Arguments = args

	return expr

}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	// defer untrace(trace("parseIndexExpression"))
	expr := &ast.IndexExpression{
		Token:  p.tok,
		Struct: left,
	}

	p.advance()

	expr.Index = p.parseExpression(LOWEST)

	if !p.expect(token.RBRACKET, "expected <]> token following <INT>") {
		return nil
	}

	return expr
}

func (p *Parser) parseMacroExpression() ast.Expression {
	// defer untrace(trace("parseMacroExpression"))
	lit := &ast.MacroExpression{
		Token: p.tok,
	}

	if !p.expect(token.LPAREN, "expected <(> token following <fn>") {
		return nil
	}

	idents := []*ast.Identifier{}

	p.advance()

	if p.tok.Type != token.RPAREN {
		idents = append(
			idents,
			&ast.Identifier{
				Token: p.tok,
				Value: p.tok.Lexeme,
			},
		)

		for p.check(token.COMMA) {
			p.advance()
			p.advance()

			idents = append(
				idents,
				&ast.Identifier{
					Token: p.tok,
					Value: p.tok.Lexeme,
				},
			)
		}

		if !p.expect(token.RPAREN, "expected <)> token following function parameters") {
			return nil
		}
	}

	lit.Parameters = idents

	if !p.expect(token.LBRACE, "expected <{> token following <)>") {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	// defer untrace(trace("parseGroupedExpression"))
	p.advance()

	expr := p.parseExpression(LOWEST)

	if !p.expect(token.RPAREN, "expected <)> token following grouped expression") {
		return nil
	}

	return expr
}

func (p *Parser) advance() {
	if p.current >= len(p.tokens) {
		p.tok = token.Token{}
	} else {
		p.tok = p.tokens[p.current]
		p.current = p.current + 1

		for p.tok.Type == token.ILLEGAL {
			p.error(p.tok, "illegal token")
			p.tok = p.tokens[p.current]
			p.current = p.current + 1
		}
	}
}

func (p *Parser) peek() token.Token {
	if p.current >= len(p.tokens) {
		return token.Token{}
	}
	return p.tokens[p.current]
}

func (p *Parser) check(ttype token.TokenType) bool {
	return p.peek().Type == ttype
}

func (p *Parser) expect(ttype token.TokenType, msg string) bool {
	if p.check(ttype) {
		p.advance()
		return true
	} else {
		p.error(p.peek(), msg)
		return false
	}
}

func (p *Parser) precedence(ttype token.TokenType) int {
	if precedence, ok := precedences[ttype]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) error(tok token.Token, msg string) {
	input := p.l.Input()

	ln := 1
	col := 1

	for i := 0; i < tok.Offset; i = i + 1 {
		if input[i] == '\n' {
			col = 1
			ln = ln + 1
		}
		col = col + 1
	}

	msg = fmt.Sprintf("%d:%d: %s", ln, col, msg)
	p.errors = append(p.errors, msg)
}
