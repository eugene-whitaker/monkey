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
	name  string
	value ExpressionTest
}

func (lst LetStatementTest) statement() {}

type ReturnStatementTest struct {
	returnValue ExpressionTest
}

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

func (iet IfExpressionTest) expression() {}

type FunctionLiteralTest struct {
	parameters []string
	body       *BlockStatementTest
}

func (flt FunctionLiteralTest) expression() {}

type CallExpressionTest struct {
	function  ExpressionTest
	arguments []ExpressionTest
}

func (cet CallExpressionTest) expression() {}

func TestParseProgram(t *testing.T) {
	tests := []struct {
		input      string
		precedence string
		tests      []StatementTest
	}{
		{
			"let x = 5;",
			"let x = 5;",
			[]StatementTest{
				LetStatementTest{
					"x",
					IntegerLiteralTest(5),
				},
			},
		},
		{
			"let y = true;",
			"let y = true;",
			[]StatementTest{
				LetStatementTest{
					"y",
					BooleanTest(true),
				},
			},
		},
		{
			"let z = y;",
			"let z = y;",
			[]StatementTest{
				LetStatementTest{
					"z",
					IdentifierTest("y"),
				},
			},
		},
		{
			"return 5;",
			"return 5;",
			[]StatementTest{
				ReturnStatementTest{
					IntegerLiteralTest(5),
				},
			},
		},
		{
			"return true;",
			"return true;",
			[]StatementTest{
				ReturnStatementTest{
					BooleanTest(true),
				},
			},
		},
		{
			"return x;",
			"return x;",
			[]StatementTest{
				ReturnStatementTest{
					IdentifierTest("x"),
				},
			},
		},
		{
			"x;",
			"x",
			[]StatementTest{
				ExpressionStatementTest{
					IdentifierTest("x"),
				},
			},
		},
		{
			"5;",
			"5",
			[]StatementTest{
				ExpressionStatementTest{
					IntegerLiteralTest(5),
				},
			},
		},
		{
			"true;",
			"true",
			[]StatementTest{
				ExpressionStatementTest{
					BooleanTest(true),
				},
			},
		},
		{
			"false;",
			"false",
			[]StatementTest{
				ExpressionStatementTest{
					BooleanTest(false),
				},
			},
		},
		{
			"!5;",
			"(!5)",
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
			"(-15)",
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
			"(!true)",
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
			"(!false)",
			[]StatementTest{
				ExpressionStatementTest{
					PrefixExpressionTest{
						"!",
						BooleanTest(false),
					},
				},
			},
		},
		{
			"5 + 5;",
			"(5 + 5)",
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
			"(5 - 5)",
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
			"(5 * 5)",
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
			"(5 / 5)",
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
			"(5 > 5)",
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
			"(5 == 5)",
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
			"(5 != 5)",
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
			"true == true;",
			"(true == true)",
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
			"true != false;",
			"(true != false)",
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
			"false == false;",
			"(false == false)",
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
		{
			"if (x < y) { x; };",
			"if (x < y) x",
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
								},
							},
						},
						nil,
					},
				},
			},
		},
		{
			"if (x < y) { x; } else { y; };",
			"if (x < y) x else y",
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
								},
							},
						},
						&BlockStatementTest{
							[]StatementTest{
								ExpressionStatementTest{
									IdentifierTest("y"),
								},
							},
						},
					},
				},
			},
		},
		{
			"fn(x, y) { x + y; };",
			"fn(x, y)(x + y)",
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
								},
							},
						},
					},
				},
			},
		},
		{
			"fn() {};",
			"fn()",
			[]StatementTest{
				ExpressionStatementTest{
					FunctionLiteralTest{
						[]string{},
						&BlockStatementTest{
							[]StatementTest{},
						},
					},
				},
			},
		},
		{
			"fn(x) {};",
			"fn(x)",
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
				},
			},
		},
		{
			"fn(x, y, z) {};",
			"fn(x, y, z)",
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
				},
			},
		},
		{
			"call(1, 2 * 3, 4 + 5)",
			"call(1, (2 * 3), (4 + 5))",
			[]StatementTest{
				ExpressionStatementTest{
					CallExpressionTest{
						IdentifierTest("call"),
						[]ExpressionTest{
							IntegerLiteralTest(1),
							InfixExpressionTest{
								IntegerLiteralTest(2),
								"*",
								IntegerLiteralTest(3),
							},
							InfixExpressionTest{
								IntegerLiteralTest(4),
								"+",
								IntegerLiteralTest(5),
							},
						},
					},
				},
			},
		},
		{
			"-a * b;",
			"((-a) * b)",
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
				},
			},
		},
		{
			"!-a;",
			"(!(-a))",
			[]StatementTest{
				ExpressionStatementTest{
					PrefixExpressionTest{
						"!",
						PrefixExpressionTest{
							"-",
							IdentifierTest("a"),
						},
					},
				},
			},
		},
		{
			"a + b + c;",
			"((a + b) + c)",
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
				},
			},
		},
		{
			"a + b - c;",
			"((a + b) - c)",
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
				},
			},
		},
		{
			"a * b * c;",
			"((a * b) * c)",
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
				},
			},
		},
		{
			"a * b / c;",
			"((a * b) / c)",
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
				},
			},
		},
		{
			"a + b / c;",
			"(a + (b / c))",
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
				},
			},
		},
		{
			"a + b * c + d / e - f;",
			"(((a + (b * c)) + (d / e)) - f)",
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
				},
			},
		},
		{
			"3 + 4; -5 * 5;",
			"(3 + 4)((-5) * 5)",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IntegerLiteralTest(3),
						"+",
						IntegerLiteralTest(4),
					},
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
				},
			},
		},
		{
			"5 > 4 == 3 < 4;",
			"((5 > 4) == (3 < 4))",
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
				},
			},
		},
		{
			"5 < 4 != 3 > 4;",
			"((5 < 4) != (3 > 4))",
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
				},
			},
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5;",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
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
				},
			},
		},
		{
			"3 > 5 == false;",
			"((3 > 5) == false)",
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
				},
			},
		},
		{
			"3 < 5 == true;",
			"((3 < 5) == true)",
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
				},
			},
		},
		{
			"1 + (2 + 3) + 4;",
			"((1 + (2 + 3)) + 4)",
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
				},
			},
		},
		{
			"(5 + 5) * 2;",
			"((5 + 5) * 2)",
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
				},
			},
		},
		{
			"2 / (5 + 5);",
			"(2 / (5 + 5))",
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
				},
			},
		},
		{
			"-(5 + 5);",
			"(-(5 + 5))",
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
				},
			},
		},
		{
			"!(true == true)",
			"(!(true == true))",
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
				},
			},
		},
		{
			"a + call(b * c) + d;",
			"((a + call((b * c))) + d)",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						InfixExpressionTest{
							IdentifierTest("a"),
							"+",
							CallExpressionTest{
								IdentifierTest("call"),
								[]ExpressionTest{
									InfixExpressionTest{
										IdentifierTest("b"),
										"*",
										IdentifierTest("c"),
									},
								},
							},
						},
						"+",
						IdentifierTest("d"),
					},
				},
			},
		},
		{
			"call(a, b, 1, 2 * 3, 4 + 5, call(6, 7 * 8));",
			"call(a, b, 1, (2 * 3), (4 + 5), call(6, (7 * 8)))",
			[]StatementTest{
				ExpressionStatementTest{
					CallExpressionTest{
						IdentifierTest("call"),
						[]ExpressionTest{
							IdentifierTest("a"),
							IdentifierTest("b"),
							IntegerLiteralTest(1),
							InfixExpressionTest{
								IntegerLiteralTest(2),
								"*",
								IntegerLiteralTest(3),
							},
							InfixExpressionTest{
								IntegerLiteralTest(4),
								"+",
								IntegerLiteralTest(5),
							},
							CallExpressionTest{
								IdentifierTest("call"),
								[]ExpressionTest{
									IntegerLiteralTest(6),
									InfixExpressionTest{
										IntegerLiteralTest(7),
										"*",
										IntegerLiteralTest(8),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"call(a + b + c * d / f + g);",
			"call((((a + b) + ((c * d) / f)) + g))",
			[]StatementTest{
				ExpressionStatementTest{
					CallExpressionTest{
						IdentifierTest("call"),
						[]ExpressionTest{
							InfixExpressionTest{
								InfixExpressionTest{
									InfixExpressionTest{
										IdentifierTest("a"),
										"+",
										IdentifierTest("b"),
									},
									"+",
									InfixExpressionTest{
										InfixExpressionTest{
											IdentifierTest("c"),
											"*",
											IdentifierTest("d"),
										},
										"/",
										IdentifierTest("f"),
									},
								},
								"+",
								IdentifierTest("g"),
							},
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		if !testProgram(t, i, tt.input, tt.precedence, tt.tests) {
			return
		}
	}
}

func testProgram(t *testing.T, index int, input string, precedence string, tests []StatementTest) bool {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if 0 != len(p.Errors()) {
		t.Errorf("test[%d] - p.Errors() ==> expected: 0 actual: %d", index, len(p.Errors()))
		for i, msg := range p.Errors() {
			t.Errorf("test[%d] - p.Errors()[%d]: %q", index, i, msg)
		}
		t.FailNow()
	}

	if program == nil {
		t.Errorf("test[%d] - ParserProgram() ==> expected: not <nil>", index)
		return false
	}

	actual := program.String()
	if precedence != actual {
		t.Errorf("test[%d] - program.String() ==> expected: %q actual: %q", index, precedence, actual)
		return false
	}

	if len(tests) != len(program.Statements) {
		t.Errorf("test[%d] - len(program.Statements) ==> expected: %d actual: %d", index, len(tests), len(program.Statements))
		return false
	}

	for i, s := range program.Statements {
		if !testStatement(t, index, s, tests[i]) {
			return false
		}
	}

	return true
}

func testStatement(t *testing.T, index int, stmt ast.Statement, test StatementTest) bool {
	switch test := test.(type) {
	case LetStatementTest:
		return testLetStatement(t, index, stmt, test.name, test.value)
	case ReturnStatementTest:
		return testReturnStatement(t, index, stmt, test.returnValue)
	case ExpressionStatementTest:
		return testExpressionStatement(t, index, stmt, test.test)
	case *BlockStatementTest:
		return testBlockStatement(t, index, stmt, test.tests)
	}
	t.Errorf("test[%d] - test ==> unexpected type. actual: %T", index, test)
	return false
}

func testLetStatement(t *testing.T, index int, stmt ast.Statement, name string, value ExpressionTest) bool {
	if "let" != stmt.TokenLiteral() {
		t.Errorf("test[%d] - stmt.TokenLiteral() ==> expected: 'let' actual: %q", index, stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("test[%d] - stmt ==> unexpected type. expected: %T actual: %T", index, &ast.LetStatement{}, stmt)
		return false
	}

	if !testIdentifier(t, index, letStmt.Name, name) {
		return false
	}

	if !testExpression(t, index, letStmt.Value, value) {
		return false
	}

	return true
}

func testReturnStatement(t *testing.T, index int, stmt ast.Statement, returnValue ExpressionTest) bool {
	if "return" != stmt.TokenLiteral() {
		t.Errorf("test[%d] - stmt.TokenLiteral() ==> expected: 'return' actual: %q", index, stmt.TokenLiteral())
		return false
	}

	returnStmt, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("test[%d] - stmt ==> unexpected type. expected: %T actual: %T", index, &ast.ReturnStatement{}, stmt)
		return false
	}

	if !testExpression(t, index, returnStmt.ReturnValue, returnValue) {
		return false
	}

	return true
}

func testExpressionStatement(t *testing.T, index int, stmt ast.Statement, test ExpressionTest) bool {
	expStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("test[%d] - stmt ==> unexpected type. expected: %T actual: %T", index, &ast.ExpressionStatement{}, stmt)
		return false
	}

	if !testExpression(t, index, expStmt.Expression, test) {
		return false
	}

	return true
}

func testBlockStatement(t *testing.T, index int, stmt ast.Statement, tests []StatementTest) bool {
	if "{" != stmt.TokenLiteral() {
		t.Errorf("test[%d] - stmt.TokenLiteral() ==> expected: '{' actual: %q", index, stmt.TokenLiteral())
		return false
	}

	blockStmt, ok := stmt.(*ast.BlockStatement)
	if !ok {
		t.Errorf("test[%d] - stmt ==> unexpected type. expected: %T actual: %T", index, &ast.BlockStatement{}, stmt)
		return false
	}

	if len(tests) != len(blockStmt.Statements) {
		t.Errorf("test[%d] - len(blockStmt.Statements) ==> expected: %d actual: %d", index, len(tests), len(blockStmt.Statements))
		return false
	}

	for i, s := range blockStmt.Statements {
		if !testStatement(t, index, s, tests[i]) {
			return false
		}
	}

	return true
}

func testExpression(t *testing.T, index int, exp ast.Expression, test ExpressionTest) bool {
	switch test := test.(type) {
	case IdentifierTest:
		return testIdentifier(t, index, exp, string(test))
	case IntegerLiteralTest:
		return testIntegerLiteral(t, index, exp, int64(test))
	case BooleanTest:
		return testBoolean(t, index, exp, bool(test))
	case PrefixExpressionTest:
		return testPrefixExpression(t, index, exp, test.operator, test.rightValue)
	case InfixExpressionTest:
		return testInfixExpression(t, index, exp, test.leftValue, test.operator, test.rightValue)
	case IfExpressionTest:
		return testIfExpression(t, index, exp, test.condition, test.consequence, test.alternative)
	case FunctionLiteralTest:
		return testFunctionLiteral(t, index, exp, test.parameters, test.body)
	case CallExpressionTest:
		return testCallExpression(t, index, exp, test.function, test.arguments)
	}
	t.Errorf("test ==> unexpected type. actual: %T", test)
	return false
}

func testIdentifier(t *testing.T, index int, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("test[%d] - exp ==> unexpected type. expected: %T actual: %T", index, &ast.Identifier{}, exp)
		return false
	}

	if value != ident.Value {
		t.Errorf("test[%d] - ident.Value ==> expected: %q actual: %q", index, value, ident.Value)
		return false
	}

	if value != ident.TokenLiteral() {
		t.Errorf("test[%d] - ident.TokenLiteral() ==> expected: %q actual: %q", index, value, ident.TokenLiteral())
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, index int, exp ast.Expression, value int64) bool {
	integ, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("test[%d] - exp ==> unexpected type. expected: %T actual: %T", index, &ast.IntegerLiteral{}, exp)
		return false
	}

	if value != integ.Value {
		t.Errorf("test[%d] - integ.Value ==> expected: %d actual: %d", index, value, integ.Value)
		return false
	}

	if fmt.Sprintf("%d", value) != integ.TokenLiteral() {
		t.Errorf("test[%d] - integ.TokenLiteral() ==> expected: %d actual: %q", index, value, integ.TokenLiteral())
		return false
	}

	return true
}

func testBoolean(t *testing.T, index int, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("test[%d] - exp ==> unexpected type. expected: %T actual: %T", index, &ast.Boolean{}, exp)
		return false
	}

	if value != boolean.Value {
		t.Errorf("test[%d] - boolean.Value ==> expected: %t actual: %t", index, value, boolean.Value)
		return false
	}

	if fmt.Sprintf("%t", value) != boolean.TokenLiteral() {
		t.Errorf("test[%d] - boolean.TokenLiteral() ==> expected: %t actual: %q", index, value, boolean.TokenLiteral())
		return false
	}

	return true
}

func testPrefixExpression(t *testing.T, index int, exp ast.Expression, operator string, right ExpressionTest) bool {
	opExp, ok := exp.(*ast.PrefixExpression)
	if !ok {
		t.Errorf("test[%d] - exp ==> unexpected type. expected: %T actual: %T", index, &ast.PrefixExpression{}, exp)
		return false
	}

	if operator != opExp.Operator {
		t.Errorf("test[%d] - opExp.Operator ==> expected: %q actual: %q", index, operator, opExp.Operator)
		return false
	}

	if !testExpression(t, index, opExp.Right, right) {
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, index int, exp ast.Expression, left ExpressionTest, operator string, right ExpressionTest) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("test[%d] - exp ==> unexpected type. expected: %T actual: %T", index, &ast.InfixExpression{}, exp)
		return false
	}

	if !testExpression(t, index, opExp.Left, left) {
		return false
	}

	if operator != opExp.Operator {
		t.Errorf("test[%d] - opExp.Operator ==> expected: %q actual: %q", index, operator, opExp.Operator)
		return false
	}

	if !testExpression(t, index, opExp.Right, right) {
		return false
	}

	return true
}

func testIfExpression(t *testing.T, index int, exp ast.Expression, condition ExpressionTest, consequence *BlockStatementTest, alternative *BlockStatementTest) bool {
	if "if" != exp.TokenLiteral() {
		t.Errorf("test[%d] - exp.TokenLiteral() ==> expected: 'if' actual: %q", index, exp.TokenLiteral())
		return false
	}

	ifExp, ok := exp.(*ast.IfExpression)
	if !ok {
		t.Errorf("test[%d] - exp ==> unexpected type. expected: %T actual: %T", index, &ast.IfExpression{}, exp)
		return false
	}

	if !testExpression(t, index, ifExp.Condition, condition) {
		return false
	}

	if !testBlockStatement(t, index, ifExp.Consequence, consequence.tests) {
		return false
	}

	if alternative != nil && !testBlockStatement(t, index, ifExp.Alternative, alternative.tests) {
		return false
	}

	return true
}

func testFunctionLiteral(t *testing.T, index int, exp ast.Expression, parameters []string, body *BlockStatementTest) bool {
	if "fn" != exp.TokenLiteral() {
		t.Errorf("test[%d] - exp.TokenLiteral() ==> expected: 'fn' actual: %q", index, exp.TokenLiteral())
		return false
	}

	fnExp, ok := exp.(*ast.FunctionLiteral)
	if !ok {
		t.Errorf("test[%d] - exp ==> unexpected type. expected: %T actual: %T", index, &ast.FunctionLiteral{}, exp)
		return false
	}

	if len(parameters) != len(fnExp.Parameters) {
		t.Errorf("test[%d] - len(fnExp.Parameters) ==> expected: %d actual: %d", index, len(parameters), len(fnExp.Parameters))
		return false
	}

	for i, parameter := range fnExp.Parameters {
		if !testIdentifier(t, index, parameter, parameters[i]) {
			return false
		}
	}

	if !testBlockStatement(t, index, fnExp.Body, body.tests) {
		return false
	}

	return true
}

func testCallExpression(t *testing.T, index int, exp ast.Expression, function ExpressionTest, arguments []ExpressionTest) bool {
	if "(" != exp.TokenLiteral() {
		t.Errorf("test[%d] - exp.TokenLiteral() ==> expected: '(' actual: %q", index, exp.TokenLiteral())
		return false
	}

	callExp, ok := exp.(*ast.CallExpression)
	if !ok {
		t.Errorf("test[%d] - exp ==> unexpected type. expected: %T actual: %T", index, &ast.CallExpression{}, exp)
		return false
	}

	switch test := function.(type) {
	case IdentifierTest:
		if !testIdentifier(t, index, callExp.Function, string(test)) {
			return false
		}
	case FunctionLiteralTest:
		if !testFunctionLiteral(t, index, callExp.Function, test.parameters, test.body) {
			return false
		}
	default:
		t.Errorf("test[%d] - test ==> unexpected type. actual: %T", index, test)
		return false
	}

	if len(arguments) != len(callExp.Arguements) {
		t.Errorf("test[%d] - len(callExp.Arguements) ==> expected: %d actual: %d", index, len(arguments), len(callExp.Arguements))
		return false
	}

	for i, argument := range callExp.Arguements {
		if !testExpression(t, index, argument, arguments[i]) {
			return false
		}
	}

	return true
}
