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
		input string
		tests []StatementTest
	}{
		{
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
			[]StatementTest{
				LetStatementTest{
					"y",
					BooleanTest(true),
				},
			},
		},
		{
			"let foobar = y;",
			[]StatementTest{
				LetStatementTest{
					"foobar",
					IdentifierTest("y"),
				},
			},
		},
		{
			"return 5;",
			[]StatementTest{
				ReturnStatementTest{
					IntegerLiteralTest(5),
				},
			},
		},
		{
			"return true;",
			[]StatementTest{
				ReturnStatementTest{
					BooleanTest(true),
				},
			},
		},
		{
			"return y;",
			[]StatementTest{
				ReturnStatementTest{
					IdentifierTest("y"),
				},
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
			"add(1, 2 * 3, 4 + 5)",
			[]StatementTest{
				ExpressionStatementTest{
					CallExpressionTest{
						IdentifierTest("add"),
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
					"add(1, (2 * 3), (4 + 5))",
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
		{
			"a + add(b * c) + d;",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						InfixExpressionTest{
							IdentifierTest("a"),
							"+",
							CallExpressionTest{
								IdentifierTest("add"),
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
					"((a + add((b * c))) + d)",
				},
			},
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8));",
			[]StatementTest{
				ExpressionStatementTest{
					CallExpressionTest{
						IdentifierTest("add"),
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
								IdentifierTest("add"),
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
					"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
				},
			},
		},
		{
			"add(a + b + c * d / f + g);",
			[]StatementTest{
				ExpressionStatementTest{
					CallExpressionTest{
						IdentifierTest("add"),
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
					"add((((a + b) + ((c * d) / f)) + g))",
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
	testParserErrors(t, p)

	if program == nil {
		t.Errorf("ParserProgram() ==> expected: not <nil>")
		return false
	}

	if len(tests) != len(program.Statements) {
		t.Errorf("len(program.Statements) ==> expected: %d actual: %d", len(tests), len(program.Statements))
		return false
	}

	for i, s := range program.Statements {
		if !testStatement(t, s, tests[i]) {
			return false
		}
	}

	return true
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string, value ExpressionTest) bool {
	if "let" != stmt.TokenLiteral() {
		t.Errorf("stmt.TokenLiteral() ==> expected: 'let' actual: %q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt ==> unexpected type. expected: %T actual: %T", &ast.LetStatement{}, stmt)
		return false
	}

	if !testIdentifier(t, letStmt.Name, name) {
		return false
	}

	if !testExpression(t, letStmt.Value, value) {
		return false
	}

	return true
}

func testReturnStatement(t *testing.T, stmt ast.Statement, returnValue ExpressionTest) bool {
	if "return" != stmt.TokenLiteral() {
		t.Errorf("stmt.TokenLiteral() ==> expected: 'return' actual: %q", stmt.TokenLiteral())
		return false
	}

	returnStmt, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("stmt ==> unexpected type. expected: %T actual: %T", &ast.ReturnStatement{}, stmt)
		return false
	}

	if !testExpression(t, returnStmt.ReturnValue, returnValue) {
		return false
	}

	return true
}

func testExpressionStatement(t *testing.T, stmt ast.Statement, test ExpressionTest) bool {
	expStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("stmt ==> unexpected type. expected: %T actual: %T", &ast.ExpressionStatement{}, stmt)
		return false
	}

	if !testExpression(t, expStmt.Expression, test) {
		return false
	}

	return true
}

func testBlockStatement(t *testing.T, stmt ast.Statement, tests []StatementTest) bool {
	if "{" != stmt.TokenLiteral() {
		t.Errorf("stmt.TokenLiteral() ==> expected: '{' actual: %q", stmt.TokenLiteral())
		return false
	}

	blockStmt, ok := stmt.(*ast.BlockStatement)
	if !ok {
		t.Errorf("stmt ==> unexpected type. expected: %T actual: %T", &ast.BlockStatement{}, stmt)
		return false
	}

	if len(tests) != len(blockStmt.Statements) {
		t.Errorf("len(blockStmt.Statements) ==> expected: %d actual: %d", len(tests), len(blockStmt.Statements))
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
		return testLetStatement(t, stmt, tt.name, tt.value)
	case ReturnStatementTest:
		return testReturnStatement(t, stmt, tt.returnValue)
	case ExpressionStatementTest:
		actual := stmt.String()
		if tt.precedence != actual {
			t.Errorf("stmt.String() ==> expected: %q actual: %q", tt.precedence, actual)
			return false
		}
		return testExpressionStatement(t, stmt, tt.test)
	case *BlockStatementTest:
		return testBlockStatement(t, stmt, tt.tests)
	}
	t.Errorf("test ==> unexpected type. actual: %T", test)
	return false
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp ==> unexpected type. expected: %T actual: %T", &ast.Identifier{}, exp)
		return false
	}

	if value != ident.Value {
		t.Errorf("ident.Value ==> expected: %q actual: %q", value, ident.Value)
		return false
	}

	if value != ident.TokenLiteral() {
		t.Errorf("ident.TokenLiteral() ==> expected: %q actual: %q", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	integ, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("exp ==> unexpected type. expected: %T actual: %T", &ast.IntegerLiteral{}, exp)
		return false
	}

	if value != integ.Value {
		t.Errorf("integ.Value ==> expected: %d actual: %d", value, integ.Value)
		return false
	}

	if fmt.Sprintf("%d", value) != integ.TokenLiteral() {
		t.Errorf("integ.TokenLiteral() ==> expected: %d actual: %q", value, integ.TokenLiteral())
		return false
	}

	return true
}

func testBoolean(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp ==> unexpected type. expected: %T actual: %T", &ast.Boolean{}, exp)
		return false
	}

	if value != boolean.Value {
		t.Errorf("boolean.Value ==> expected: %t actual: %t", value, boolean.Value)
		return false
	}

	if fmt.Sprintf("%t", value) != boolean.TokenLiteral() {
		t.Errorf("boolean.TokenLiteral() ==> expected: %t actual: %q", value, boolean.TokenLiteral())
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
	t.Errorf("test ==> unexpected type. actual: %T", test)
	return false
}

func testPrefixExpression(t *testing.T, exp ast.Expression, operator string, right ExpressionTest) bool {
	opExp, ok := exp.(*ast.PrefixExpression)
	if !ok {
		t.Errorf("exp ==> unexpected type. expected: %T actual: %T", &ast.PrefixExpression{}, exp)
		return false
	}

	if operator != opExp.Operator {
		t.Errorf("opExp.Operator ==> expected: %q actual: %q", operator, opExp.Operator)
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
		t.Errorf("exp ==> unexpected type. expected: %T actual: %T", &ast.InfixExpression{}, exp)
		return false
	}

	if !testExpression(t, opExp.Left, left) {
		return false
	}

	if operator != opExp.Operator {
		t.Errorf("opExp.Operator ==> expected: %q actual: %q", operator, opExp.Operator)
		return false
	}

	if !testExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testIfExpression(t *testing.T, exp ast.Expression, condition ExpressionTest, consequence *BlockStatementTest, alternative *BlockStatementTest) bool {
	if "if" != exp.TokenLiteral() {
		t.Errorf("exp.TokenLiteral() ==> expected: 'if' actual: %q", exp.TokenLiteral())
		return false
	}

	ifExp, ok := exp.(*ast.IfExpression)
	if !ok {
		t.Errorf("exp ==> unexpected type. expected: %T actual: %T", &ast.IfExpression{}, exp)
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
	if "fn" != exp.TokenLiteral() {
		t.Errorf("exp.TokenLiteral() ==> expected: 'fn' actual: %q", exp.TokenLiteral())
		return false
	}

	fnExp, ok := exp.(*ast.FunctionLiteral)
	if !ok {
		t.Errorf("exp ==> unexpected type. expected: %T actual: %T", &ast.FunctionLiteral{}, exp)
		return false
	}

	if len(parameters) != len(fnExp.Parameters) {
		t.Errorf("len(fnExp.Parameters) ==> expected: %d actual: %d", len(parameters), len(fnExp.Parameters))
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

func testCallExpression(t *testing.T, exp ast.Expression, function ExpressionTest, arguments []ExpressionTest) bool {
	if "(" != exp.TokenLiteral() {
		t.Errorf("exp.TokenLiteral() ==> expected: '(' actual: %q", exp.TokenLiteral())
		return false
	}

	callExp, ok := exp.(*ast.CallExpression)
	if !ok {
		t.Errorf("exp ==> unexpected type. expected: %T actual: %T", &ast.CallExpression{}, exp)
		return false
	}

	switch test := function.(type) {
	case IdentifierTest:
		if !testIdentifier(t, callExp.Function, string(test)) {
			return false
		}
	case FunctionLiteralTest:
		if !testFunctionLiteral(t, callExp.Function, test.parameters, test.body) {
			return false
		}
	default:
		t.Errorf("test ==> unexpected type. actual: %T", test)
		return false
	}

	if len(arguments) != len(callExp.Arguements) {
		t.Errorf("len(callExp.Arguements) ==> expected: %d actual: %d", len(arguments), len(callExp.Arguements))
		return false
	}

	for i, a := range callExp.Arguements {
		if !testExpression(t, a, arguments[i]) {
			return false
		}
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
	case CallExpressionTest:
		return testCallExpression(t, exp, tt.function, tt.arguments)
	}
	t.Errorf("test ==> unexpected type. actual: %T", test)
	return false
}

func testParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if 0 != len(errors) {
		t.Errorf("p.Errors() ==> expected: 0 actual: %d", len(errors))
		for i, msg := range errors {
			t.Errorf("p.Errors()[%d]: %q", i, msg)
		}
		t.FailNow()
	}
}
