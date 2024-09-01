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
	name string
}

func (lst LetStatementTest) statement() {}

type ReturnStatementTest struct{}

func (rst ReturnStatementTest) statement() {}

type ExpressionStatementTest struct {
	test       ExpressionTest
	precedence string
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

type FunctionLiteralTest struct {
	parameters []string
	body       *BlockStatementTest
}

func (flt FunctionLiteralTest) expression() {}

func TestParseProgram(t *testing.T) {
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
		{
			"foobar;",
			[]StatementTest{
				ExpressionStatementTest{
					IdentifierTest("foobar"),
					"foobar",
				},
			},
		},
		{
			"5;",
			[]StatementTest{
				ExpressionStatementTest{
					IntegerLiteralTest(5),
					"5",
				},
			},
		},
		{
			"true;",
			[]StatementTest{
				ExpressionStatementTest{
					BooleanTest(true),
					"true",
				},
			},
		},
		{
			"false;",
			[]StatementTest{
				ExpressionStatementTest{
					BooleanTest(false),
					"false",
				},
			},
		},
		{
			"!5;",
			[]StatementTest{
				ExpressionStatementTest{
					PrefixExpressionTest{
						"!",
						IntegerLiteralTest(5),
					},
					"(!5)",
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
					"(-15)",
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
					"(!true)",
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
					"(!false)",
				},
			},
		},
		{
			"5 + 5;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IntegerLiteralTest(5),
						"+",
						IntegerLiteralTest(5),
					},
					"(5 + 5)",
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
					"(5 - 5)",
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
					"(5 * 5)",
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
					"(5 / 5)",
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
					"(5 > 5)",
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
					"(5 == 5)",
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
					"(5 != 5)",
				},
			},
		},
		{
			"true == true;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						BooleanTest(true),
						"==",
						BooleanTest(true),
					},
					"(true == true)",
				},
			},
		},
		{
			"true != false;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						BooleanTest(true),
						"!=",
						BooleanTest(false),
					},
					"(true != false)",
				},
			},
		},
		{
			"false == false;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						BooleanTest(false),
						"==",
						BooleanTest(false),
					},
					"(false == false)",
				},
			},
		},
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
								ExpressionStatementTest{
									IdentifierTest("x"),
									"x",
								},
							},
						},
						nil,
					},
					"if (x < y) x",
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
								ExpressionStatementTest{
									IdentifierTest("x"),
									"x",
								},
							},
						},
						&BlockStatementTest{
							[]StatementTest{
								ExpressionStatementTest{
									IdentifierTest("y"),
									"y",
								},
							},
						},
					},
					"if (x < y) x else y",
				},
			},
		},
		{
			"fn(x, y) { x + y; };",
			[]StatementTest{
				ExpressionStatementTest{
					FunctionLiteralTest{
						[]string{
							"x",
							"y",
						},
						&BlockStatementTest{
							[]StatementTest{
								ExpressionStatementTest{
									InfixExpressionTest{
										IdentifierTest("x"),
										"+",
										IdentifierTest("y"),
									},
									"(x + y)",
								},
							},
						},
					},
					"fn(x, y)(x + y)",
				},
			},
		},
		{
			"fn() {};",
			[]StatementTest{
				ExpressionStatementTest{
					FunctionLiteralTest{
						[]string{},
						&BlockStatementTest{
							[]StatementTest{},
						},
					},
					"fn()",
				},
			},
		},
		{
			"fn(x) {};",
			[]StatementTest{
				ExpressionStatementTest{
					FunctionLiteralTest{
						[]string{
							"x",
						},
						&BlockStatementTest{
							[]StatementTest{},
						},
					},
					"fn(x)",
				},
			},
		},
		{
			"fn(x, y, z) {};",
			[]StatementTest{
				ExpressionStatementTest{
					FunctionLiteralTest{
						[]string{
							"x",
							"y",
							"z",
						},
						&BlockStatementTest{
							[]StatementTest{},
						},
					},
					"fn(x, y, z)",
				},
			},
		},
		{
			"-a * b;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						PrefixExpressionTest{
							"-",
							IdentifierTest("a"),
						},
						"*",
						IdentifierTest("b"),
					},
					"((-a) * b)",
				},
			},
		},
		{
			"!-a;",
			[]StatementTest{
				ExpressionStatementTest{
					PrefixExpressionTest{
						"!",
						PrefixExpressionTest{
							"-",
							IdentifierTest("a"),
						},
					},
					"(!(-a))",
				},
			},
		},
		{
			"a + b + c;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						InfixExpressionTest{
							IdentifierTest("a"),
							"+",
							IdentifierTest("b"),
						},
						"+",
						IdentifierTest("c"),
					},
					"((a + b) + c)",
				},
			},
		},
		{
			"a + b - c;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						InfixExpressionTest{
							IdentifierTest("a"),
							"+",
							IdentifierTest("b"),
						},
						"-",
						IdentifierTest("c"),
					},
					"((a + b) - c)",
				},
			},
		},
		{
			"a * b * c;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						InfixExpressionTest{
							IdentifierTest("a"),
							"*",
							IdentifierTest("b"),
						},
						"*",
						IdentifierTest("c"),
					},
					"((a * b) * c)",
				},
			},
		},
		{
			"a * b / c;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						InfixExpressionTest{
							IdentifierTest("a"),
							"*",
							IdentifierTest("b"),
						},
						"/",
						IdentifierTest("c"),
					},
					"((a * b) / c)",
				},
			},
		},
		{
			"a + b / c;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IdentifierTest("a"),
						"+",
						InfixExpressionTest{
							IdentifierTest("b"),
							"/",
							IdentifierTest("c"),
						},
					},
					"(a + (b / c))",
				},
			},
		},
		{
			"a + b * c + d / e - f;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						InfixExpressionTest{
							InfixExpressionTest{
								IdentifierTest("a"),
								"+",
								InfixExpressionTest{
									IdentifierTest("b"),
									"*",
									IdentifierTest("c"),
								},
							},
							"+",
							InfixExpressionTest{
								IdentifierTest("d"),
								"/",
								IdentifierTest("e"),
							},
						},
						"-",
						IdentifierTest("f"),
					},
					"(((a + (b * c)) + (d / e)) - f)",
				},
			},
		},
		{
			"3 + 4; -5 * 5;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IntegerLiteralTest(3),
						"+",
						IntegerLiteralTest(4),
					},
					"(3 + 4)",
				},
				ExpressionStatementTest{
					InfixExpressionTest{
						PrefixExpressionTest{
							"-",
							IntegerLiteralTest(5),
						},
						"*",
						IntegerLiteralTest(5),
					},
					"((-5) * 5)",
				},
			},
		},
		{
			"5 > 4 == 3 < 4;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						InfixExpressionTest{
							IntegerLiteralTest(5),
							">",
							IntegerLiteralTest(4),
						},
						"==",
						InfixExpressionTest{
							IntegerLiteralTest(3),
							"<",
							IntegerLiteralTest(4),
						},
					},
					"((5 > 4) == (3 < 4))",
				},
			},
		},
		{
			"5 < 4 != 3 > 4;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						InfixExpressionTest{
							IntegerLiteralTest(5),
							"<",
							IntegerLiteralTest(4),
						},
						"!=",
						InfixExpressionTest{
							IntegerLiteralTest(3),
							">",
							IntegerLiteralTest(4),
						},
					},
					"((5 < 4) != (3 > 4))",
				},
			},
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						InfixExpressionTest{
							IntegerLiteralTest(3),
							"+",
							InfixExpressionTest{
								IntegerLiteralTest(4),
								"*",
								IntegerLiteralTest(5),
							},
						},
						"==",
						InfixExpressionTest{
							InfixExpressionTest{
								IntegerLiteralTest(3),
								"*",
								IntegerLiteralTest(1),
							},
							"+",
							InfixExpressionTest{
								IntegerLiteralTest(4),
								"*",
								IntegerLiteralTest(5),
							},
						},
					},
					"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
				},
			},
		},
		{
			"true;",
			[]StatementTest{
				ExpressionStatementTest{
					BooleanTest(true),
					"true",
				},
			},
		},
		{
			"false;",
			[]StatementTest{
				ExpressionStatementTest{
					BooleanTest(false),
					"false",
				},
			},
		},
		{
			"3 > 5 == false;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						InfixExpressionTest{
							IntegerLiteralTest(3),
							">",
							IntegerLiteralTest(5),
						},
						"==",
						BooleanTest(false),
					},
					"((3 > 5) == false)",
				},
			},
		},
		{
			"3 < 5 == true;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						InfixExpressionTest{
							IntegerLiteralTest(3),
							"<",
							IntegerLiteralTest(5),
						},
						"==",
						BooleanTest(true),
					},
					"((3 < 5) == true)",
				},
			},
		},
		{
			"1 + (2 + 3) + 4;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						InfixExpressionTest{
							IntegerLiteralTest(1),
							"+",
							InfixExpressionTest{
								IntegerLiteralTest(2),
								"+",
								IntegerLiteralTest(3),
							},
						},
						"+",
						IntegerLiteralTest(4),
					},
					"((1 + (2 + 3)) + 4)",
				},
			},
		},
		{
			"(5 + 5) * 2;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						InfixExpressionTest{
							IntegerLiteralTest(5),
							"+",
							IntegerLiteralTest(5),
						},
						"*",
						IntegerLiteralTest(2),
					},
					"((5 + 5) * 2)",
				},
			},
		},
		{
			"2 / (5 + 5);",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IntegerLiteralTest(2),
						"/",
						InfixExpressionTest{
							IntegerLiteralTest(5),
							"+",
							IntegerLiteralTest(5),
						},
					},
					"(2 / (5 + 5))",
				},
			},
		},
		{
			"-(5 + 5);",
			[]StatementTest{
				ExpressionStatementTest{
					PrefixExpressionTest{
						"-",
						InfixExpressionTest{
							IntegerLiteralTest(5),
							"+",
							IntegerLiteralTest(5),
						},
					},
					"(-(5 + 5))",
				},
			},
		},
		{
			"!(true == true)",
			[]StatementTest{
				ExpressionStatementTest{
					PrefixExpressionTest{
						"!",
						InfixExpressionTest{
							BooleanTest(true),
							"==",
							BooleanTest(true),
						},
					},
					"(!(true == true))",
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

func testStatement(t *testing.T, stmt ast.Statement, test StatementTest) bool {
	switch tt := test.(type) {
	case LetStatementTest:
		return testLetStatement(t, stmt, tt.name)
	case ReturnStatementTest:
		return testReturnStatement(t, stmt)
	case ExpressionStatementTest:
		actual := stmt.String()
		if actual != tt.precedence {
			t.Errorf("expected=%q, got=%q", tt.precedence, actual)
			return false
		}
		return testExpressionStatement(t, stmt, tt.test)
	case *BlockStatementTest:
		return testBlockStatement(t, stmt, tt.tests)
	}
	t.Errorf("type of test not handled. got=%T", test)
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

func testLiteralExpression(t *testing.T, exp ast.Expression, test ExpressionTest) bool {
	switch v := test.(type) {
	case IdentifierTest:
		return testIdentifier(t, exp, string(v))
	case IntegerLiteralTest:
		return testIntegerLiteral(t, exp, int64(v))
	case BooleanTest:
		return testBoolean(t, exp, bool(v))
	}
	t.Errorf("type of test not handled. got=%T", test)
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

	if !testBlockStatement(t, ifExp.Consequence, consequence.tests) {
		return false
	}

	if alternative != nil && !testBlockStatement(t, ifExp.Alternative, alternative.tests) {
		return false
	}

	return true
}

func testFunctionLiteral(t *testing.T, exp ast.Expression, parameters []string, body *BlockStatementTest) bool {
	if exp.TokenLiteral() != "fn" {
		t.Errorf("exp.TokenLiteral() not 'fn'. got=%q", exp.TokenLiteral())
		return false
	}

	fnExp, ok := exp.(*ast.FunctionLiteral)
	if !ok {
		t.Errorf("exp not *ast.FunctionLiteral. got=%T(%s)", exp, exp)
		return false
	}

	if len(fnExp.Parameters) != len(parameters) {
		t.Errorf("fnExp.Parameters does not contain %d statements. got=%d", len(parameters), len(fnExp.Parameters))
		return false
	}

	for i, p := range fnExp.Parameters {
		if !testIdentifier(t, p, parameters[i]) {
			return false
		}
	}

	if !testBlockStatement(t, fnExp.Body, body.tests) {
		return false
	}

	return true
}

func testExpression(t *testing.T, exp ast.Expression, test ExpressionTest) bool {
	switch tt := test.(type) {
	case IdentifierTest, IntegerLiteralTest, BooleanTest:
		return testLiteralExpression(t, exp, test)
	case PrefixExpressionTest:
		return testPrefixExpression(t, exp, tt.operator, tt.rightValue)
	case InfixExpressionTest:
		return testInfixExpression(t, exp, tt.leftValue, tt.operator, tt.rightValue)
	case IfExpressionTest:
		return testIfExpression(t, exp, tt.condition, tt.consequence, tt.alternative)
	case FunctionLiteralTest:
		return testFunctionLiteral(t, exp, tt.parameters, tt.body)
	}
	t.Errorf("type of test not handled. got=%T", test)
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
