package parser

import (
	"fmt"
	"testing"

	"github.com/eugene-whitaker/writing-an-interpreter-in-go/ast"
	"github.com/eugene-whitaker/writing-an-interpreter-in-go/lexer"
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

type BooleanLiteralTest bool

func (blt BooleanLiteralTest) expression() {}

type FunctionLiteralTest struct {
	parameters []string
	body       *BlockStatementTest
}

func (flt FunctionLiteralTest) expression() {}

type StringLiteralTest string

func (slt StringLiteralTest) expression() {}

type ArrayLiteralTest []ExpressionTest

func (alt ArrayLiteralTest) expression() {}

type HashLiteralTest map[any]ExpressionTest

func (hlt HashLiteralTest) expression() {}

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

type CallExpressionTest struct {
	function  ExpressionTest
	arguments []ExpressionTest
}

func (cet CallExpressionTest) expression() {}

type IndexExpressionTest struct {
	array ExpressionTest
	index ExpressionTest
}

func (iet IndexExpressionTest) expression() {}

type MacroExpressionTest struct {
	parameters []string
	body       *BlockStatementTest
}

func (met MacroExpressionTest) expression() {}

func TestParseProgram(t *testing.T) {
	tests := []struct {
		input      string
		precedence string
		tests      []StatementTest
	}{
		{
			"let ident = 5;",
			"let ident = 5;",
			[]StatementTest{
				LetStatementTest{
					"ident",
					IntegerLiteralTest(5),
				},
			},
		},
		{
			"let ident = true;",
			"let ident = true;",
			[]StatementTest{
				LetStatementTest{
					"ident",
					BooleanLiteralTest(true),
				},
			},
		},
		{
			"let ident = value;",
			"let ident = value;",
			[]StatementTest{
				LetStatementTest{
					"ident",
					IdentifierTest("value"),
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
					BooleanLiteralTest(true),
				},
			},
		},
		{
			"return ident;",
			"return ident;",
			[]StatementTest{
				ReturnStatementTest{
					IdentifierTest("ident"),
				},
			},
		},
		{
			"ident;",
			"ident",
			[]StatementTest{
				ExpressionStatementTest{
					IdentifierTest("ident"),
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
					BooleanLiteralTest(true),
				},
			},
		},
		{
			"false;",
			"false",
			[]StatementTest{
				ExpressionStatementTest{
					BooleanLiteralTest(false),
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
			"!ident;",
			"(!ident)",
			[]StatementTest{
				ExpressionStatementTest{
					PrefixExpressionTest{
						"!",
						IdentifierTest("ident"),
					},
				},
			},
		},
		{
			"-ident;",
			"(-ident)",
			[]StatementTest{
				ExpressionStatementTest{
					PrefixExpressionTest{
						"-",
						IdentifierTest("ident"),
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
						BooleanLiteralTest(true),
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
						BooleanLiteralTest(false),
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
			"5 < 5;",
			"(5 < 5)",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IntegerLiteralTest(5),
						"<",
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
			"ident + ident;",
			"(ident + ident)",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IdentifierTest("ident"),
						"+",
						IdentifierTest("ident"),
					},
				},
			},
		},
		{
			"ident - ident;",
			"(ident - ident)",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IdentifierTest("ident"),
						"-",
						IdentifierTest("ident"),
					},
				},
			},
		},
		{
			"ident * ident;",
			"(ident * ident)",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IdentifierTest("ident"),
						"*",
						IdentifierTest("ident"),
					},
				},
			},
		},
		{
			"ident / ident;",
			"(ident / ident)",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IdentifierTest("ident"),
						"/",
						IdentifierTest("ident"),
					},
				},
			},
		},
		{
			"ident > ident;",
			"(ident > ident)",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IdentifierTest("ident"),
						">",
						IdentifierTest("ident"),
					},
				},
			},
		},
		{
			"ident < ident;",
			"(ident < ident)",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IdentifierTest("ident"),
						"<",
						IdentifierTest("ident"),
					},
				},
			},
		},
		{
			"ident == ident;",
			"(ident == ident)",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IdentifierTest("ident"),
						"==",
						IdentifierTest("ident"),
					},
				},
			},
		},
		{
			"ident != ident;",
			"(ident != ident)",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						IdentifierTest("ident"),
						"!=",
						IdentifierTest("ident"),
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
						BooleanLiteralTest(true),
						"==",
						BooleanLiteralTest(true),
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
						BooleanLiteralTest(true),
						"!=",
						BooleanLiteralTest(false),
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
						BooleanLiteralTest(false),
						"==",
						BooleanLiteralTest(false),
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
			"call()",
			"call()",
			[]StatementTest{
				ExpressionStatementTest{
					CallExpressionTest{
						IdentifierTest("call"),
						[]ExpressionTest{},
					},
				},
			},
		},
		{
			"call(1)",
			"call(1)",
			[]StatementTest{
				ExpressionStatementTest{
					CallExpressionTest{
						IdentifierTest("call"),
						[]ExpressionTest{
							IntegerLiteralTest(1),
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
						BooleanLiteralTest(false),
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
						BooleanLiteralTest(true),
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
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						InfixExpressionTest{
							InfixExpressionTest{
								IntegerLiteralTest(5),
								"+",
								IntegerLiteralTest(5),
							},
							"*",
							IntegerLiteralTest(2),
						},
						"*",
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
							BooleanLiteralTest(true),
							"==",
							BooleanLiteralTest(true),
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
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
			[]StatementTest{
				ExpressionStatementTest{
					InfixExpressionTest{
						InfixExpressionTest{
							IdentifierTest("a"),
							"*",
							IndexExpressionTest{
								ArrayLiteralTest(
									[]ExpressionTest{
										IntegerLiteralTest(1),
										IntegerLiteralTest(2),
										IntegerLiteralTest(3),
										IntegerLiteralTest(4),
									},
								),
								InfixExpressionTest{
									IdentifierTest("b"),
									"*",
									IdentifierTest("c"),
								},
							},
						},
						"*",
						IdentifierTest("d"),
					},
				},
			},
		},
		{
			"call(a * b[2], b[1], 2 * [1, 2][1])",
			"call((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
			[]StatementTest{
				ExpressionStatementTest{
					CallExpressionTest{
						IdentifierTest("call"),
						[]ExpressionTest{
							InfixExpressionTest{
								IdentifierTest("a"),
								"*",
								IndexExpressionTest{
									IdentifierTest("b"),
									IntegerLiteralTest(2),
								},
							},
							IndexExpressionTest{
								IdentifierTest("b"),
								IntegerLiteralTest(1),
							},
							InfixExpressionTest{
								IntegerLiteralTest(2),
								"*",
								IndexExpressionTest{
									ArrayLiteralTest(
										[]ExpressionTest{
											IntegerLiteralTest(1),
											IntegerLiteralTest(2),
										},
									),
									IntegerLiteralTest(1),
								},
							},
						},
					},
				},
			},
		},
		{
			"\"hello world\";",
			"hello world",
			[]StatementTest{
				ExpressionStatementTest{
					StringLiteralTest("hello world"),
				},
			},
		},
		{
			"[1, 2 * 2, 3 + 3];",
			"[1, (2 * 2), (3 + 3)]",
			[]StatementTest{
				ExpressionStatementTest{
					ArrayLiteralTest(
						[]ExpressionTest{
							IntegerLiteralTest(1),
							InfixExpressionTest{
								IntegerLiteralTest(2),
								"*",
								IntegerLiteralTest(2),
							},
							InfixExpressionTest{
								IntegerLiteralTest(3),
								"+",
								IntegerLiteralTest(3),
							},
						},
					),
				},
			},
		},
		{
			"[];",
			"[]",
			[]StatementTest{
				ExpressionStatementTest{
					ArrayLiteralTest([]ExpressionTest{}),
				},
			},
		},
		{
			"array[1 + 1];",
			"(array[(1 + 1)])",
			[]StatementTest{
				ExpressionStatementTest{
					IndexExpressionTest{
						IdentifierTest("array"),
						InfixExpressionTest{
							IntegerLiteralTest(1),
							"+",
							IntegerLiteralTest(1),
						},
					},
				},
			},
		},
		{
			"{\"one\": 1, \"two\": 2, \"three\": 3};",
			"{one:1, two:2, three:3}",
			[]StatementTest{
				ExpressionStatementTest{
					HashLiteralTest{
						"one":   IntegerLiteralTest(1),
						"two":   IntegerLiteralTest(2),
						"three": IntegerLiteralTest(3),
					},
				},
			},
		},
		{
			"{};",
			"{}",
			[]StatementTest{
				ExpressionStatementTest{
					HashLiteralTest{},
				},
			},
		},
		{
			"{\"one\": 0 + 1, \"two\": 10 - 8, \"three\": 15 / 5};",
			"{one:(0 + 1), two:(10 - 8), three:(15 / 5)}",
			[]StatementTest{
				ExpressionStatementTest{
					HashLiteralTest{
						"one": InfixExpressionTest{
							IntegerLiteralTest(0),
							"+",
							IntegerLiteralTest(1),
						},
						"two": InfixExpressionTest{
							IntegerLiteralTest(10),
							"-",
							IntegerLiteralTest(8),
						},
						"three": InfixExpressionTest{
							IntegerLiteralTest(15),
							"/",
							IntegerLiteralTest(5),
						},
					},
				},
			},
		},
		{
			"{1: 1, 2: 2, 3: 3};",
			"{1:1, 2:2, 3:3}",
			[]StatementTest{
				ExpressionStatementTest{
					HashLiteralTest{
						1: IntegerLiteralTest(1),
						2: IntegerLiteralTest(2),
						3: IntegerLiteralTest(3),
					},
				},
			},
		},
		{
			"{true: 1, false: 2};",
			"{true:1, false:2}",
			[]StatementTest{
				ExpressionStatementTest{
					HashLiteralTest{
						true:  IntegerLiteralTest(1),
						false: IntegerLiteralTest(2),
					},
				},
			},
		},
		{
			"macro(x, y) { x + y; };",
			"macro(x, y)(x + y)",
			[]StatementTest{
				ExpressionStatementTest{
					MacroExpressionTest{
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
	}

	for i, test := range tests {
		if !testProgram(t, i, test.input, test.precedence, test.tests) {
			continue
		}
	}
}

func testProgram(t *testing.T, idx int, input string, precedence string, tests []StatementTest) bool {
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if 0 != len(p.Errors()) {
		t.Errorf("test[%d] - %q - len(p.Errors()) ==> expected: 0 actual: %d", idx, input, len(p.Errors()))
		for i, msg := range p.Errors() {
			t.Errorf("test[%d] - p.Errors()[%d]: %q", idx, i, msg)
		}
		t.FailNow()
	}

	if program == nil {
		t.Errorf("test[%d] - %q - p.ParseProgram() ==> expected: not <nil>", idx, input)
		return false
	}

	actual := program.String()
	if precedence != actual {
		t.Logf("test[%d] - %q - program.String() ==> expected: %q actual: %q", idx, input, precedence, actual)
		return false
	}

	if len(tests) != len(program.Statements) {
		t.Errorf("test[%d] - %q - len(program.Statements) ==> expected: %d actual: %d", idx, input, len(tests), len(program.Statements))
		return false
	}

	for i, stmt := range program.Statements {
		if !testStatement(t, idx, input, stmt, tests[i]) {
			return false
		}
	}

	return true
}

func testStatement(t *testing.T, idx int, input string, stmt ast.Statement, test StatementTest) bool {
	switch test := test.(type) {
	case LetStatementTest:
		return testLetStatement(t, idx, input, stmt, test.name, test.value)
	case ReturnStatementTest:
		return testReturnStatement(t, idx, input, stmt, test.returnValue)
	case ExpressionStatementTest:
		return testExpressionStatement(t, idx, input, stmt, test.test)
	case *BlockStatementTest:
		return testBlockStatement(t, idx, input, stmt, test.tests)
	}
	t.Errorf("test[%d] - %q ==> unexpected type. actual: %T", idx, input, test)
	return false
}

func testLetStatement(t *testing.T, idx int, input string, stmt ast.Statement, name string, value ExpressionTest) bool {
	if "let" != stmt.TokenLexeme() {
		t.Errorf("test[%d] - %q - stmt.TokenLexeme() ==> expected: 'let' actual: %q", idx, input, stmt.TokenLexeme())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("test[%d] - %q - stmt.(*ast.LetStatement) ==> unexpected type. expected: %T actual: %T", idx, input, &ast.LetStatement{}, stmt)
		return false
	}

	if !testIdentifier(t, idx, input, letStmt.Name, name) {
		return false
	}

	if !testExpression(t, idx, input, letStmt.Value, value) {
		return false
	}

	return true
}

func testReturnStatement(t *testing.T, idx int, input string, stmt ast.Statement, returnValue ExpressionTest) bool {
	if "return" != stmt.TokenLexeme() {
		t.Errorf("test[%d] - %q - stmt.TokenLexeme() ==> expected: 'return' actual: %q", idx, input, stmt.TokenLexeme())
		return false
	}

	returnStmt, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("test[%d] - %q - stmt.(*ast.ReturnStatement) ==> unexpected type. expected: %T actual: %T", idx, input, &ast.ReturnStatement{}, stmt)
		return false
	}

	if !testExpression(t, idx, input, returnStmt.ReturnValue, returnValue) {
		return false
	}

	return true
}

func testExpressionStatement(t *testing.T, idx int, input string, stmt ast.Statement, test ExpressionTest) bool {
	exprStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("test[%d] - %q - stmt.(*ast.ExpressionStatement) ==> unexpected type. expected: %T actual: %T", idx, input, &ast.ExpressionStatement{}, stmt)
		return false
	}

	if !testExpression(t, idx, input, exprStmt.Expression, test) {
		return false
	}

	return true
}

func testBlockStatement(t *testing.T, idx int, input string, stmt ast.Statement, tests []StatementTest) bool {
	if "{" != stmt.TokenLexeme() {
		t.Errorf("test[%d] - %q - stmt.TokenLexeme() ==> expected: '{' actual: %q", idx, input, stmt.TokenLexeme())
		return false
	}

	blockStmt, ok := stmt.(*ast.BlockStatement)
	if !ok {
		t.Errorf("test[%d] - %q - stmt.(*ast.BlockStatement) ==> unexpected type. expected: %T actual: %T", idx, input, &ast.BlockStatement{}, stmt)
		return false
	}

	if len(tests) != len(blockStmt.Statements) {
		t.Errorf("test[%d] - %q - len(blockStmt.Statements) ==> expected: %d actual: %d", idx, input, len(tests), len(blockStmt.Statements))
		return false
	}

	for i, stmt := range blockStmt.Statements {
		if !testStatement(t, idx, input, stmt, tests[i]) {
			return false
		}
	}

	return true
}

func testExpression(t *testing.T, idx int, input string, exp ast.Expression, test ExpressionTest) bool {
	switch test := test.(type) {
	case IdentifierTest:
		return testIdentifier(t, idx, input, exp, string(test))
	case IntegerLiteralTest:
		return testIntegerLiteral(t, idx, input, exp, int64(test))
	case BooleanLiteralTest:
		return testBooleanLiteral(t, idx, input, exp, bool(test))
	case FunctionLiteralTest:
		return testFunctionLiteral(t, idx, input, exp, test.parameters, test.body)
	case StringLiteralTest:
		return testStringLiteral(t, idx, input, exp, string(test))
	case ArrayLiteralTest:
		return testArrayLiteral(t, idx, input, exp, []ExpressionTest(test))
	case HashLiteralTest:
		return testHashLiteral(t, idx, input, exp, map[any]ExpressionTest(test))
	case PrefixExpressionTest:
		return testPrefixExpression(t, idx, input, exp, test.operator, test.rightValue)
	case InfixExpressionTest:
		return testInfixExpression(t, idx, input, exp, test.leftValue, test.operator, test.rightValue)
	case IfExpressionTest:
		return testIfExpression(t, idx, input, exp, test.condition, test.consequence, test.alternative)
	case CallExpressionTest:
		return testCallExpression(t, idx, input, exp, test.function, test.arguments)
	case IndexExpressionTest:
		return testIndexExpression(t, idx, input, exp, test.array, test.index)
	case MacroExpressionTest:
		return testMacroExpression(t, idx, input, exp, test.parameters, test.body)
	}
	t.Errorf("test[%d] - %q ==> unexpected type. actual: %T", idx, input, test)
	return false
}

func testIdentifier(t *testing.T, idx int, input string, expr ast.Expression, value string) bool {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Errorf("test[%d] - %q - exp.(*ast.Identifier) ==> unexpected type. expected: %T actual: %T", idx, input, &ast.Identifier{}, expr)
		return false
	}

	if value != ident.Value {
		t.Errorf("test[%d] - %q - ident.Value ==> expected: %q actual: %q", idx, input, value, ident.Value)
		return false
	}

	if value != ident.TokenLexeme() {
		t.Errorf("test[%d] - %q - ident.TokenLexeme() ==> expected: %q actual: %q", idx, input, value, ident.TokenLexeme())
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, idx int, input string, expr ast.Expression, value int64) bool {
	integ, ok := expr.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("test[%d] - %q - exp.(*ast.IntegerLiteral) ==> unexpected type. expected: %T actual: %T", idx, input, &ast.IntegerLiteral{}, expr)
		return false
	}

	if value != integ.Value {
		t.Errorf("test[%d] - %q - integ.Value ==> expected: %d actual: %d", idx, input, value, integ.Value)
		return false
	}

	if fmt.Sprintf("%d", value) != integ.TokenLexeme() {
		t.Errorf("test[%d] - %q - integ.TokenLexeme() ==> expected: %d actual: %q", idx, input, value, integ.TokenLexeme())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, idx int, input string, expr ast.Expression, value bool) bool {
	boolean, ok := expr.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("test[%d] - %q - exp.(*ast.BooleanLiteral) ==> unexpected type. expected: %T actual: %T", idx, input, &ast.BooleanLiteral{}, expr)
		return false
	}

	if value != boolean.Value {
		t.Errorf("test[%d] - %q - boolean.Value ==> expected: %t actual: %t", idx, input, value, boolean.Value)
		return false
	}

	if fmt.Sprintf("%t", value) != boolean.TokenLexeme() {
		t.Errorf("test[%d] - %q - boolean.TokenLexeme() ==> expected: %t actual: %q", idx, input, value, boolean.TokenLexeme())
		return false
	}

	return true
}

func testFunctionLiteral(t *testing.T, idx int, input string, expr ast.Expression, params []string, body *BlockStatementTest) bool {
	if "fn" != expr.TokenLexeme() {
		t.Errorf("test[%d] - %q - exp.TokenLexeme() ==> expected: 'fn' actual: %q", idx, input, expr.TokenLexeme())
		return false
	}

	fn, ok := expr.(*ast.FunctionLiteral)
	if !ok {
		t.Errorf("test[%d] - %q - exp.(*ast.FunctionLiteral) ==> unexpected type. expected: %T actual: %T", idx, input, &ast.FunctionLiteral{}, expr)
		return false
	}

	if len(params) != len(fn.Parameters) {
		t.Errorf("test[%d] - %q - len(fn.Parameters) ==> expected: %d actual: %d", idx, input, len(params), len(fn.Parameters))
		return false
	}

	for i, param := range fn.Parameters {
		if !testIdentifier(t, idx, input, param, params[i]) {
			return false
		}
	}

	if !testBlockStatement(t, idx, input, fn.Body, body.tests) {
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, idx int, input string, expr ast.Expression, value string) bool {
	str, ok := expr.(*ast.StringLiteral)
	if !ok {
		t.Errorf("test[%d] - %q - exp.(*ast.StringLiteral) ==> unexpected type. expected: %T actual: %T", idx, input, &ast.StringLiteral{}, expr)
		return false
	}

	if value != str.Value {
		t.Errorf("test[%d] - %q - str.Value ==> expected: %q actual: %q", idx, input, value, str.Value)
		return false
	}

	if value != str.TokenLexeme() {
		t.Errorf("test[%d] - %q - str.TokenLexeme() ==> expected: %q actual: %q", idx, input, value, str.TokenLexeme())
		return false
	}

	return true
}

func testArrayLiteral(t *testing.T, idx int, input string, expr ast.Expression, tests []ExpressionTest) bool {
	arr, ok := expr.(*ast.ArrayLiteral)
	if !ok {
		t.Errorf("test[%d] - %q - exp.(*ast.ArrayLiteral) ==> unexpected type. expected: %T actual: %T", idx, input, &ast.ArrayLiteral{}, expr)
		return false
	}

	if len(tests) != len(arr.Elements) {
		t.Errorf("test[%d] - %q - len(arr.Elements) ==> expected: %d actual: %d", idx, input, len(tests), len(arr.Elements))
		return false
	}

	for i, elem := range arr.Elements {
		if !testExpression(t, idx, input, elem, tests[i]) {
			return false
		}
	}

	return true
}

func testHashLiteral(t *testing.T, idx int, input string, expr ast.Expression, tests map[any]ExpressionTest) bool {
	hash, ok := expr.(*ast.HashLiteral)
	if !ok {
		t.Errorf("test[%d] - %q - exp.(*ast.HashLiteral) ==> unexpected type. expected: %T actual: %T", idx, input, &ast.HashLiteral{}, expr)
		return false
	}

	if len(tests) != len(hash.Pairs) {
		t.Errorf("test[%d] - %q - len(hash.Pairs) ==> expected: %d actual: %d", idx, input, len(tests), len(hash.Pairs))
		return false
	}

	for key, value := range hash.Pairs {
		switch key := key.(type) {
		case *ast.IntegerLiteral:
			expected, ok := tests[int(key.Value)]
			if !ok {
				t.Errorf("test[%d] - %q - tests[key.Value] ==> expected: not <nil>", idx, input)
				return false
			}
			if !testExpression(t, idx, input, value, expected) {
				return false
			}
		case *ast.BooleanLiteral:
			expected, ok := tests[key.Value]
			if !ok {
				t.Errorf("test[%d] - %q - tests[key.Value] ==> expected: not <nil>", idx, input)
				return false
			}
			if !testExpression(t, idx, input, value, expected) {
				return false
			}
		case *ast.StringLiteral:
			expected, ok := tests[key.Value]
			if !ok {
				t.Errorf("test[%d] - %q - tests[key.Value] ==> expected: not <nil>", idx, input)
				return false
			}
			if !testExpression(t, idx, input, value, expected) {
				return false
			}
		default:
			t.Errorf("test[%d] - %q ==> unexpected type. actual: %T", idx, input, key)
			return false
		}
	}

	return true
}

func testPrefixExpression(t *testing.T, idx int, input string, expr ast.Expression, operator string, right ExpressionTest) bool {
	opExpr, ok := expr.(*ast.PrefixExpression)
	if !ok {
		t.Errorf("test[%d] - %q - exp.(*ast.PrefixExpression) ==> unexpected type. expected: %T actual: %T", idx, input, &ast.PrefixExpression{}, expr)
		return false
	}

	if operator != opExpr.Operator {
		t.Errorf("test[%d] - %q - opExp.Operator ==> expected: %q actual: %q", idx, input, operator, opExpr.Operator)
		return false
	}

	if !testExpression(t, idx, input, opExpr.Right, right) {
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, idx int, input string, expr ast.Expression, left ExpressionTest, operator string, right ExpressionTest) bool {
	opExpr, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Errorf("test[%d] - %q - exp.(*ast.InfixExpression) ==> unexpected type. expected: %T actual: %T", idx, input, &ast.InfixExpression{}, expr)
		return false
	}

	if !testExpression(t, idx, input, opExpr.Left, left) {
		return false
	}

	if operator != opExpr.Operator {
		t.Errorf("test[%d] - %q - opExp.Operator ==> expected: %q actual: %q", idx, input, operator, opExpr.Operator)
		return false
	}

	if !testExpression(t, idx, input, opExpr.Right, right) {
		return false
	}

	return true
}

func testIfExpression(t *testing.T, idx int, input string, expr ast.Expression, condition ExpressionTest, consequence *BlockStatementTest, alternative *BlockStatementTest) bool {
	if "if" != expr.TokenLexeme() {
		t.Errorf("test[%d] - %q - exp.TokenLexeme() ==> expected: 'if' actual: %q", idx, input, expr.TokenLexeme())
		return false
	}

	ifExpr, ok := expr.(*ast.IfExpression)
	if !ok {
		t.Errorf("test[%d] - %q - exp.(*ast.IfExpression) ==> unexpected type. expected: %T actual: %T", idx, input, &ast.IfExpression{}, expr)
		return false
	}

	if !testExpression(t, idx, input, ifExpr.Condition, condition) {
		return false
	}

	if !testBlockStatement(t, idx, input, ifExpr.Consequence, consequence.tests) {
		return false
	}

	if alternative != nil && !testBlockStatement(t, idx, input, ifExpr.Alternative, alternative.tests) {
		return false
	}

	return true
}

func testCallExpression(t *testing.T, idx int, input string, expr ast.Expression, function ExpressionTest, arguments []ExpressionTest) bool {
	if "(" != expr.TokenLexeme() {
		t.Errorf("test[%d] - %q - exp.TokenLexeme() ==> expected: '(' actual: %q", idx, input, expr.TokenLexeme())
		return false
	}

	callExpr, ok := expr.(*ast.CallExpression)
	if !ok {
		t.Errorf("test[%d] - %q - exp.(*ast.CallExpression) ==> unexpected type. expected: %T actual: %T", idx, input, &ast.CallExpression{}, expr)
		return false
	}

	switch test := function.(type) {
	case IdentifierTest:
		if !testIdentifier(t, idx, input, callExpr.Function, string(test)) {
			return false
		}
	case FunctionLiteralTest:
		if !testFunctionLiteral(t, idx, input, callExpr.Function, test.parameters, test.body) {
			return false
		}
	default:
		t.Errorf("test[%d] - %q - test ==> unexpected type. actual: %T", idx, input, test)
		return false
	}

	if len(arguments) != len(callExpr.Arguments) {
		t.Errorf("test[%d] - %q - len(callExp.Arguements) ==> expected: %d actual: %d", idx, input, len(arguments), len(callExpr.Arguments))
		return false
	}

	for i, arg := range callExpr.Arguments {
		if !testExpression(t, idx, input, arg, arguments[i]) {
			return false
		}
	}

	return true
}

func testIndexExpression(t *testing.T, idx int, input string, expr ast.Expression, array ExpressionTest, index ExpressionTest) bool {
	if "[" != expr.TokenLexeme() {
		t.Errorf("test[%d] - %q - exp.TokenLexeme() ==> expected: '[' actual: %q", idx, input, expr.TokenLexeme())
		return false
	}

	indexExpr, ok := expr.(*ast.IndexExpression)
	if !ok {
		t.Errorf("test[%d] - %q - exp.(*ast.IndexExpression) ==> unexpected type. expected: %T actual: %T", idx, input, &ast.IndexExpression{}, expr)
		return false
	}

	switch test := array.(type) {
	case IdentifierTest:
		if !testIdentifier(t, idx, input, indexExpr.Struct, string(test)) {
			return false
		}
	case ArrayLiteralTest:
		if !testArrayLiteral(t, idx, input, indexExpr.Struct, []ExpressionTest(test)) {
			return false
		}
	default:
		t.Errorf("test[%d] - %q ==> unexpected type. actual: %T", idx, input, test)
		return false
	}

	if !testExpression(t, idx, input, indexExpr.Index, index) {
		return false
	}

	return true
}

func testMacroExpression(t *testing.T, idx int, input string, expr ast.Expression, params []string, body *BlockStatementTest) bool {
	if "macro" != expr.TokenLexeme() {
		t.Errorf("test[%d] - %q - exp.TokenLexeme() ==> expected: 'macro' actual: %q", idx, input, expr.TokenLexeme())
		return false
	}

	macroExpr, ok := expr.(*ast.MacroExpression)
	if !ok {
		t.Errorf("test[%d] - %q - exp.(*ast.FunctionLiteral) ==> unexpected type. expected: %T actual: %T", idx, input, &ast.FunctionLiteral{}, expr)
		return false
	}

	if len(params) != len(macroExpr.Parameters) {
		t.Errorf("test[%d] - %q - len(macroExpr.Parameters) ==> expected: %d actual: %d", idx, input, len(params), len(macroExpr.Parameters))
		return false
	}

	for i, param := range macroExpr.Parameters {
		if !testIdentifier(t, idx, input, param, params[i]) {
			return false
		}
	}

	if !testBlockStatement(t, idx, input, macroExpr.Body, body.tests) {
		return false
	}

	return true
}
