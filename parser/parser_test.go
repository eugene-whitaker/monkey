package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

type StatementTest interface {
	statement()
}

type LetStatementTest struct {
	expectedName string
}

func (lst LetStatementTest) statement() {}

type ReturnStatementTest struct{}

func (rst ReturnStatementTest) statement() {}

type ExpressionStatementTest struct {
	test ExpressionTest
}

func (est ExpressionStatementTest) statement() {}

type BlockStatementTest struct {
	tests []StatementTest
}

func (bst BlockStatementTest) statement() {}

type ExpressionTest interface {
	expression()
}

type IdentifierTest string

func (it IdentifierTest) expression() {}

type IntegerLiteralTest int64

func (ilt IntegerLiteralTest) expression() {}

type BooleanTest bool

func (bt BooleanTest) expression() {}

type PrefixExpressionTest struct {
	operator   string
	rightValue ExpressionTest
}

func (pet PrefixExpressionTest) expression() {}

type InfixExpressionTest struct {
	leftValue  ExpressionTest
	operator   string
	rightValue ExpressionTest
}

func (iet InfixExpressionTest) expression() {}

type IfExpressionTest struct {
	condition   ExpressionTest
	consequence *BlockStatementTest
	alternative *BlockStatementTest
}

func (ifet IfExpressionTest) expression() {}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input string
		tests []StatementTest
	}{
		{
			"let x = 5;",
			[]StatementTest{
				LetStatementTest{
					"x",
				},
			},
		},
		{
			"let y = 10;",
			[]StatementTest{
				LetStatementTest{
					"y",
				},
			},
		},
		{
			"let foobar = 838383;",
			[]StatementTest{
				LetStatementTest{
					"foobar",
				},
			},
		},
	}

	for _, tt := range tests {
		if !testProgram(t, tt.input, tt.tests) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input string
		tests []StatementTest
	}{
		{
			"return 5;",
			[]StatementTest{
				ReturnStatementTest{},
			},
		},
		{
			"return 10;",
			[]StatementTest{
				ReturnStatementTest{},
			},
		},
		{
			"return 993322;",
			[]StatementTest{
				ReturnStatementTest{},
			},
		},
	}

	for _, tt := range tests {
		if !testProgram(t, tt.input, tt.tests) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	tests := []struct {
		input string
		tests []StatementTest
	}{
		{
			"foobar;",
			[]StatementTest{
				ExpressionStatementTest{
					IdentifierTest("foobar"),
				},
			},
		},
	}

	for _, tt := range tests {
		if !testProgram(t, tt.input, tt.tests) {
			return
		}
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	tests := []struct {
		input string
		tests []StatementTest
	}{
		{
			"5;",
			[]StatementTest{
				ExpressionStatementTest{
					IntegerLiteralTest(5),
				},
			},
		},
	}

	for _, tt := range tests {
		if !testProgram(t, tt.input, tt.tests) {
			return
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		tests []StatementTest
	}{
		{
			"true;",
			[]StatementTest{
				ExpressionStatementTest{
					BooleanTest(true),
				},
			},
		},
		{
			"false;",
			[]StatementTest{
				ExpressionStatementTest{
					BooleanTest(false),
				},
			},
		},
	}

	for _, tt := range tests {
		if !testProgram(t, tt.input, tt.tests) {
			return
		}
	}
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		input string
		tests []StatementTest
	}{
		{
			"!5;",
			[]StatementTest{
				ExpressionStatementTest{
					PrefixExpressionTest{
						"!",
						IntegerLiteralTest(5),
					},
				},
			},
		},
		{
			"-15;",
			[]StatementTest{
				ExpressionStatementTest{
					PrefixExpressionTest{
						"-",
						IntegerLiteralTest(15),
					},
				},
			},
		},
		{
			"!true;",
			[]StatementTest{
				ExpressionStatementTest{
					PrefixExpressionTest{
						"!",
						BooleanTest(true),
					},
				},
			},
		},
		{
			"!false;",
			[]StatementTest{
				ExpressionStatementTest{
					PrefixExpressionTest{
						"!",
						BooleanTest(false),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		if !testProgram(t, tt.input, tt.tests) {
			return
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	tests := []struct {
		input string
		tests []StatementTest
	}{
		{
			"5 + 5;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IntegerLiteralTest(5),
						"+",
						IntegerLiteralTest(5),
					},
				},
			},
		},
		{
			"5 - 5;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IntegerLiteralTest(5),
						"-",
						IntegerLiteralTest(5),
					},
				},
			},
		},
		{
			"5 * 5;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IntegerLiteralTest(5),
						"*",
						IntegerLiteralTest(5),
					},
				},
			},
		},
		{
			"5 / 5;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IntegerLiteralTest(5),
						"/",
						IntegerLiteralTest(5),
					},
				},
			},
		},
		{
			"5 > 5;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IntegerLiteralTest(5),
						">",
						IntegerLiteralTest(5),
					},
				},
			},
		},
		{
			"5 == 5;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IntegerLiteralTest(5),
						"==",
						IntegerLiteralTest(5),
					},
				},
			},
		},
		{
			"5 != 5;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IntegerLiteralTest(5),
						"!=",
						IntegerLiteralTest(5),
					},
				},
			},
		},
		{
			"true == true",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						BooleanTest(true),
						"==",
						BooleanTest(true),
					},
				},
			},
		},
		{
			"true != false",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						BooleanTest(true),
						"!=",
						BooleanTest(false),
					},
				},
			},
		},
		{
			"false == false",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						BooleanTest(false),
						"==",
						BooleanTest(false),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		if !testProgram(t, tt.input, tt.tests) {
			return
		}
	}
}

func TestIfExpression(t *testing.T) {
	tests := []struct {
		input string
		tests []StatementTest
	}{
		{
			"if (x < y) { x };",
			[]StatementTest{
				ExpressionStatementTest{
					IfExpressionTest{
						InfixExpressionTest{
							IdentifierTest("x"),
							"<",
							IdentifierTest("y"),
						},
						&BlockStatementTest{
							[]StatementTest{
								ExpressionStatementTest{IdentifierTest("x")},
							},
						},
						nil,
					},
				},
			},
		},
		{
			"if (x < y) { x } else { y };",
			[]StatementTest{
				ExpressionStatementTest{
					IfExpressionTest{
						InfixExpressionTest{
							IdentifierTest("x"),
							"<",
							IdentifierTest("y"),
						},
						&BlockStatementTest{
							[]StatementTest{
								ExpressionStatementTest{IdentifierTest("x")},
							},
						},
						&BlockStatementTest{
							[]StatementTest{
								ExpressionStatementTest{IdentifierTest("y")},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		if !testProgram(t, tt.input, tt.tests) {
			return
		}
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testProgram(t *testing.T, input string, tests []StatementTest) bool {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Errorf("ParseProgram() returned nil")
		return false
	}

	if len(program.Statements) != len(tests) {
		t.Errorf("program.Statements does not contain %d statements. got=%d", len(tests), len(program.Statements))
		return false
	}

	for i, s := range program.Statements {
		if !testStatement(t, s, tests[i]) {
			return false
		}
	}

	return true
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral() not 'let'. got=%q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt not *ast.LetStatement. got=%T(%s)", stmt, stmt)
		return false
	}

	if !testIdentifier(t, letStmt.Name, name) {
		return false
	}

	return true
}

func testReturnStatement(t *testing.T, stmt ast.Statement) bool {
	if stmt.TokenLiteral() != "return" {
		t.Errorf("stmt.TokenLiteral() not 'return'. got=%q", stmt.TokenLiteral())
		return false
	}

	_, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("stmt not *ast.ReturnStatement. got=%T(%s)", stmt, stmt)
		return false
	}

	return true
}

func testExpressionStatement(t *testing.T, stmt ast.Statement, test ExpressionTest) bool {
	expStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("stmt not *ast.ExpressionStatement. got=%T(%s)", stmt, stmt)
		return false
	}

	if !testExpression(t, expStmt.Expression, test) {
		return false
	}

	return true
}

func testBlockStatement(t *testing.T, stmt ast.Statement, tests []StatementTest) bool {
	if stmt.TokenLiteral() != "{" {
		t.Errorf("stmt.TokenLiteral() not '{'. got=%q", stmt.TokenLiteral())
		return false
	}

	blockStmt, ok := stmt.(*ast.BlockStatement)
	if !ok {
		t.Errorf("stmt not *ast.BlockStatement. got=%T(%s)", stmt, stmt)
		return false
	}

	if len(blockStmt.Statements) != len(tests) {
		t.Errorf("blockStmt.Statements does not contain %d statements. got=%d", len(tests), len(blockStmt.Statements))
		return false
	}

	for i, s := range blockStmt.Statements {
		if !testStatement(t, s, tests[i]) {
			return false
		}
	}

	return true
}

func testStatement(t *testing.T, stmt ast.Statement, expected StatementTest) bool {
	switch tt := expected.(type) {
	case LetStatementTest:
		return testLetStatement(t, stmt, tt.expectedName)
	case ReturnStatementTest:
		return testReturnStatement(t, stmt)
	case ExpressionStatementTest:
		return testExpressionStatement(t, stmt, tt.test)
	case *BlockStatementTest:
		return testBlockStatement(t, stmt, tt.tests)
	}
	t.Errorf("type of expected not handled. got=%T", expected)
	return false
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T(%s)", exp, exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	integ, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("exp not *ast.IntegerLiteral. got=%T(%s)", exp, exp)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral() not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func testBoolean(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T(%s)", exp, exp)
		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value not %t. got=%t", value, boolean.Value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral() not %t. got=%s", value, boolean.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected ExpressionTest) bool {
	switch v := expected.(type) {
	case IdentifierTest:
		return testIdentifier(t, exp, string(v))
	case IntegerLiteralTest:
		return testIntegerLiteral(t, exp, int64(v))
	case BooleanTest:
		return testBoolean(t, exp, bool(v))
	}
	t.Errorf("type of expected not handled. got=%T", expected)
	return false
}

func testPrefixExpression(t *testing.T, exp ast.Expression, operator string, right ExpressionTest) bool {
	opExp, ok := exp.(*ast.PrefixExpression)
	if !ok {
		t.Errorf("exp is not ast.PrefixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left ExpressionTest, operator string, right ExpressionTest) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testIfExpression(t *testing.T, exp ast.Expression, condition ExpressionTest, consequence *BlockStatementTest, alternative *BlockStatementTest) bool {
	if exp.TokenLiteral() != "if" {
		t.Errorf("exp.TokenLiteral() not 'if'. got=%q", exp.TokenLiteral())
		return false
	}

	ifExp, ok := exp.(*ast.IfExpression)
	if !ok {
		t.Errorf("exp not *ast.IfExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testExpression(t, ifExp.Condition, condition) {
		return false
	}

	if !testStatement(t, ifExp.Consequence, consequence) {
		return false
	}

	if alternative != nil && !testStatement(t, ifExp.Alternative, alternative) {
		return false
	}

	return true
}

func testExpression(t *testing.T, exp ast.Expression, expected ExpressionTest) bool {
	switch tt := expected.(type) {
	case IdentifierTest, IntegerLiteralTest, BooleanTest:
		return testLiteralExpression(t, exp, expected)
	case PrefixExpressionTest:
		return testPrefixExpression(t, exp, tt.operator, tt.rightValue)
	case InfixExpressionTest:
		return testInfixExpression(t, exp, tt.leftValue, tt.operator, tt.rightValue)
	case IfExpressionTest:
		return testIfExpression(t, exp, tt.condition, tt.consequence, tt.alternative)
	}
	t.Errorf("type of expected not handled. got=%T", expected)
	return false
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
