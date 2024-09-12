package evaluator

import (
	"testing"

	"github.com/eugene-whitaker/writing-an-interpreter-in-go/lexer"
	"github.com/eugene-whitaker/writing-an-interpreter-in-go/object"
	"github.com/eugene-whitaker/writing-an-interpreter-in-go/parser"
)

type ObjectTest interface {
	object()
}

type IntegerTest int64

func (iot IntegerTest) object() {}

type BooleanTest bool

func (bt BooleanTest) object() {}

type NullTest struct{}

func (nt NullTest) object() {}

type ErrorTest struct {
	message string
}

func (et ErrorTest) object() {}

type FunctionTest struct {
	parameters []string
	body       string
}

func (ft FunctionTest) object() {}

type StringTest string

func (st StringTest) object() {}

type ArrayTest []ObjectTest

func (at ArrayTest) object() {}

type HashTest map[object.HashKey]ObjectTest

func (ht HashTest) object() {}

type QuoteTest struct {
	node string
}

func (qt QuoteTest) object() {}

func TestEval(t *testing.T) {
	tests := []struct {
		input string
		test  ObjectTest
	}{
		{
			"5;",
			IntegerTest(5),
		},
		{
			"10;",
			IntegerTest(10),
		},
		{
			"-5;",
			IntegerTest(-5),
		},
		{
			"-10;",
			IntegerTest(-10),
		},
		{
			"5 + 5 + 5 + 5 - 10;",
			IntegerTest(10),
		},
		{
			"2 * 2 * 2 * 2 * 2;",
			IntegerTest(32),
		},
		{
			"-50 + 100 + -50;",
			IntegerTest(0),
		},
		{
			"5 * 2 + 10;",
			IntegerTest(20),
		},
		{
			"5 + 2 * 10;",
			IntegerTest(25),
		},
		{
			"20 + 2 * -10;",
			IntegerTest(0),
		},
		{
			"50 / 2 * 2 + 10;",
			IntegerTest(60),
		},
		{
			"2 * (5 + 10);",
			IntegerTest(30),
		},
		{
			"3 * 3 * 3 + 10;",
			IntegerTest(37),
		},
		{
			"3 * (3 * 3) + 10;",
			IntegerTest(37),
		},
		{
			"(5 + 10 * 2 + 15 / 3) * 2 + -10;",
			IntegerTest(50),
		},
		{
			"true;",
			BooleanTest(true),
		},
		{
			"false;",
			BooleanTest(false),
		},
		{
			"1 < 2;",
			BooleanTest(true),
		},
		{
			"1 > 2;",
			BooleanTest(false),
		},
		{
			"1 < 1;",
			BooleanTest(false),
		},
		{
			"1 > 1;",
			BooleanTest(false),
		},
		{
			"1 == 1;",
			BooleanTest(true),
		},
		{
			"1 != 1;",
			BooleanTest(false),
		},
		{
			"1 == 2;",
			BooleanTest(false),
		},
		{
			"1 != 2;",
			BooleanTest(true),
		},
		{
			"true == true;",
			BooleanTest(true),
		},
		{
			"false == false;",
			BooleanTest(true),
		},
		{
			"true == false;",
			BooleanTest(false),
		},
		{
			"true != false;",
			BooleanTest(true),
		},
		{
			"false != true;",
			BooleanTest(true),
		},
		{
			"(1 < 2) == true;",
			BooleanTest(true),
		},
		{
			"(1 < 2) == false;",
			BooleanTest(false),
		},
		{
			"(1 > 2) == true;",
			BooleanTest(false),
		},
		{
			"(1 > 2) == false;",
			BooleanTest(true),
		},
		{
			"!true;",
			BooleanTest(false),
		},
		{
			"!false;",
			BooleanTest(true),
		},
		{
			"!5;",
			BooleanTest(false),
		},
		{
			"!!true;",
			BooleanTest(true),
		},
		{
			"!!false;",
			BooleanTest(false),
		},
		{
			"!!5;",
			BooleanTest(true),
		},
		{
			"if (true) { 10 };",
			IntegerTest(10),
		},
		{
			"if (false) { 10 };",
			NullTest{},
		},
		{
			"if (1) { 10 };",
			IntegerTest(10),
		},
		{
			"if (1 < 2) { 10 };",
			IntegerTest(10),
		},
		{
			"if (1 > 2) { 10 };",
			NullTest{},
		},
		{
			"if (1 > 2) { 10 } else { 20 };",
			IntegerTest(20),
		},
		{
			"if (1 < 2) { 10 } else { 20 };",
			IntegerTest(10),
		},
		{
			"return 10;",
			IntegerTest(10),
		},
		{
			"return 10; 9;",
			IntegerTest(10),
		},
		{
			"return 2 * 5; 9;",
			IntegerTest(10),
		},
		{
			"9; return 2 * 5; 9;",
			IntegerTest(10),
		},
		{
			"if (10 > 1) {  return 10; };",
			IntegerTest(10),
		},
		{
			"if (10 > 1) { if (10 > 1) { return 10; }; return 1; };",
			IntegerTest(10),
		},
		{
			"let f = fn(x) { return x; x + 10; }; f(10);",
			IntegerTest(10),
		},
		{
			"let f = fn(x) { let result = x + 10; return result; return 10; }; f(10);",
			IntegerTest(20),
		},
		{
			"5 + true;",
			ErrorTest{
				"unknown operation: INTEGER + BOOLEAN",
			},
		},
		{
			"5 + true; 5;",
			ErrorTest{
				"unknown operation: INTEGER + BOOLEAN",
			},
		},
		{
			"-true;",
			ErrorTest{
				"unknown operation: -BOOLEAN",
			},
		},
		{
			"true + false;",
			ErrorTest{
				"unknown operation: BOOLEAN + BOOLEAN",
			},
		},
		{
			"5; true + false; 5;",
			ErrorTest{
				"unknown operation: BOOLEAN + BOOLEAN",
			},
		},
		{
			"if (10 > 1) { true + false; };",
			ErrorTest{
				"unknown operation: BOOLEAN + BOOLEAN",
			},
		},
		{
			"if (10 > 1) { if (10 > 1) { return true + false; } return 1; };",
			ErrorTest{
				"unknown operation: BOOLEAN + BOOLEAN",
			},
		},
		{
			"ident",
			ErrorTest{
				"undefined reference: ident",
			},
		},
		{
			"\"hello\" - \"world\"",
			ErrorTest{
				"unknown operation: STRING - STRING",
			},
		},
		{
			"{\"key\": \"value\"}[fn(x) { x; }];",
			ErrorTest{
				"unknown operation: HASH[FUNCTION]",
			},
		},
		{
			"999[1];",
			ErrorTest{
				"unknown operation: INTEGER[INTEGER]",
			},
		},
		{
			"let a = 5; a;",
			IntegerTest(5),
		},
		{
			"let a = 5 * 5; a;",
			IntegerTest(25),
		},
		{
			"let a = 5; let b = a; b;",
			IntegerTest(5),
		},
		{
			"let a = 5; let b = a; let c = a + b + 5; c;",
			IntegerTest(15),
		},
		{
			"fn(x) { x + 2; };",
			FunctionTest{
				[]string{
					"x",
				},
				"(x + 2)",
			},
		},
		{
			"let identity = fn(x) { x; }; identity(5);",
			IntegerTest(5),
		},
		{
			"let identity = fn(x) { return x; }; identity(5);",
			IntegerTest(5),
		},
		{
			"let double = fn(x) { x * 2; }; double(5);",
			IntegerTest(10),
		},
		{
			"let add = fn(x, y) { x + y; }; add(5, 5);",
			IntegerTest(10),
		},
		{
			"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));",
			IntegerTest(20),
		},
		{
			"fn(x) { x; }(5);",
			IntegerTest(5),
		},
		{
			"let adder = fn(x) { fn(y) { x + y }; }; let addTwo = adder(2); addTwo(2);",
			IntegerTest(4),
		},
		{
			"let first = 10; let second = 10; let third = 10; let func = fn(first) { let second = 20; first + second + third; }; func(20) + first + second;",
			IntegerTest(70),
		},
		{
			"\"hello world\";",
			StringTest("hello world"),
		},
		{
			"\"hello\" + \" \" + \"world\";",
			StringTest("hello world"),
		},
		{
			"len(\"\");",
			IntegerTest(0),
		},
		{
			"len(\"four\");",
			IntegerTest(4),
		},
		{
			"len(\"hello world\");",
			IntegerTest(11),
		},
		{
			"len(1)",
			ErrorTest{
				"invalid argument types in call to `len`: found (INTEGER) want (STRING) or (ARRAY)",
			},
		},
		{
			"len(\"one\", \"two\")",
			ErrorTest{
				"invalid argument count in call to `len`: found (STRING, STRING) want (STRING) or (ARRAY)",
			},
		},
		{
			"len([1, 2, 3])",
			IntegerTest(3),
		},
		{
			"len([])",
			IntegerTest(0),
		},
		{
			"first([1, 2, 3])",
			IntegerTest(1),
		},
		{
			"first([])",
			NullTest{},
		},
		{
			"first(1)",
			ErrorTest{
				"invalid argument types in call to `first`: found (INTEGER) want (ARRAY)",
			},
		},
		{
			"last([1, 2, 3])",
			IntegerTest(3),
		},
		{
			"last([])",
			NullTest{},
		},
		{
			"last(1)",
			ErrorTest{
				"invalid argument types in call to `last`: found (INTEGER) want (ARRAY)",
			},
		},
		{
			"rest([1, 2, 3])",
			ArrayTest(
				[]ObjectTest{
					IntegerTest(2),
					IntegerTest(3),
				},
			),
		},
		{
			"rest([])",
			NullTest{},
		},
		{
			"push([], 1)",
			ArrayTest(
				[]ObjectTest{
					IntegerTest(1),
				},
			),
		},
		{
			"push(1, 1)",
			ErrorTest{
				"invalid argument types in call to `push`: found (INTEGER, INTEGER) want (ARRAY, ANY)",
			},
		},
		{
			"[1, 2 * 2, 3 + 3];",
			ArrayTest(
				[]ObjectTest{
					IntegerTest(1),
					IntegerTest(4),
					IntegerTest(6),
				},
			),
		},
		{
			"[1, 2, 3][0];",
			IntegerTest(1),
		},
		{
			"[1, 2, 3][1]",
			IntegerTest(2),
		},
		{
			"[1, 2, 3][2]",
			IntegerTest(3),
		},
		{
			"let i = 0; [1][i];",
			IntegerTest(1),
		},
		{
			"[1, 2, 3][1 + 1];",
			IntegerTest(3),
		},
		{
			"let array = [1, 2, 3]; array[2];",
			IntegerTest(3),
		},
		{
			"let array = [1, 2, 3]; array[0] + array[1] + array[2];",
			IntegerTest(6),
		},
		{
			"let array = [1, 2, 3]; let i = array[0]; array[i]",
			IntegerTest(2),
		},
		{
			"[1, 2, 3][3]",
			NullTest{},
		},
		{
			"[1, 2, 3][-1]",
			NullTest{},
		},
		{
			"let two = \"two\"; {\"one\": 10 - 9, two: 1 + 1, \"thr\" + \"ee\": 6 / 2, 4: 4, true: 5, false: 6}",
			HashTest(
				map[object.HashKey]ObjectTest{
					(&object.String{Value: "one"}).HashKey():   IntegerTest(1),
					(&object.String{Value: "two"}).HashKey():   IntegerTest(2),
					(&object.String{Value: "three"}).HashKey(): IntegerTest(3),
					(&object.Integer{Value: 4}).HashKey():      IntegerTest(4),
					(&object.Boolean{Value: true}).HashKey():   IntegerTest(5),
					(&object.Boolean{Value: false}).HashKey():  IntegerTest(6),
				},
			),
		},
		{
			"{\"foo\": 5}[\"foo\"]",
			IntegerTest(5),
		},
		{
			"{\"foo\": 5}[\"bar\"]",
			NullTest{},
		},
		{
			"let key = \"foo\"; {\"foo\": 5}[key]",
			IntegerTest(5),
		},
		{
			"{}[\"foo\"]",
			NullTest{},
		},
		{
			"{5: 5}[5]",
			IntegerTest(5),
		},
		{
			"{true: 5}[true]",
			IntegerTest(5),
		},
		{
			"quote(5)",
			QuoteTest{
				"5",
			},
		},
		{
			"quote(5 + 8)",
			QuoteTest{
				"(5 + 8)",
			},
		},
		{
			"quote(ident)",
			QuoteTest{
				"ident",
			},
		},
		{
			"quote(ident + ident)",
			QuoteTest{
				"(ident + ident)",
			},
		},
		{
			"quote(unquote(4))",
			QuoteTest{
				"4",
			},
		},
		{
			"quote(unquote(4 + 4))",
			QuoteTest{
				"8",
			},
		},
		{
			"quote(8 + unquote(4 + 4))",
			QuoteTest{
				"(8 + 8)",
			},
		},
		{
			"quote(unquote(4 + 4) + 8)",
			QuoteTest{
				"(8 + 8)",
			},
		},
		{
			"let ident = 8; quote(ident)",
			QuoteTest{
				"ident",
			},
		},
		{
			"let ident = 8; quote(unquote(ident))",
			QuoteTest{
				"8",
			},
		},
		{
			"quote(unquote(true))",
			QuoteTest{
				"true",
			},
		},
		{
			"quote(unquote(true == false))",
			QuoteTest{
				"false",
			},
		},
		{
			"quote(unquote(quote(4 + 4)))",
			QuoteTest{
				"(4 + 4)",
			},
		},
		{
			"let quoted = quote(4 + 4); quote(unquote(4 + 4) + unquote(quoted))",
			QuoteTest{
				"(8 + (4 + 4))",
			},
		},
	}

	for i, test := range tests {
		l := lexer.NewLexer(test.input)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		eval := Eval(program, env)

		if !testObject(t, i, test.input, eval, test.test) {
			continue
		}
	}
}

func testObject(t *testing.T, idx int, input string, obj object.Object, test ObjectTest) bool {
	switch test := test.(type) {
	case IntegerTest:
		return testInteger(t, idx, input, obj, int64(test))
	case BooleanTest:
		return testBoolean(t, idx, input, obj, bool(test))
	case NullTest:
		return testNull(t, idx, input, obj)
	case ErrorTest:
		return testError(t, idx, input, obj, test.message)
	case FunctionTest:
		return testFunction(t, idx, input, obj, test.parameters, test.body)
	case StringTest:
		return testString(t, idx, input, obj, string(test))
	case ArrayTest:
		return testArray(t, idx, input, obj, []ObjectTest(test))
	case HashTest:
		return testHash(t, idx, input, obj, map[object.HashKey]ObjectTest(test))
	case QuoteTest:
		return testQuote(t, idx, input, obj, test.node)
	}
	t.Errorf("test[%d] - %q ==> unexpected type. actual: %T", idx, input, test)
	return false
}

func testInteger(t *testing.T, idx int, input string, obj object.Object, value int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("test[%d] - %q - obj ==> unexpected type. expected: %T actual: %T", idx, input, object.Integer{}, obj)
		return false
	}

	if value != result.Value {
		t.Errorf("test[%d] - %q - result.Value ==> expected: %d actual: %d", idx, input, value, result.Value)
		return false
	}

	return true
}

func testBoolean(t *testing.T, idx int, input string, obj object.Object, value bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("test[%d] - %q - obj ==> unexpected type. expected: %T actual: %T", idx, input, object.Boolean{}, obj)
		return false
	}

	if value != result.Value {
		t.Errorf("test[%d] - %q - result.Value ==> expected: %t actual: %t", idx, input, value, result.Value)
		return false
	}

	return true
}

func testNull(t *testing.T, idx int, input string, obj object.Object) bool {
	if NULL != obj {
		t.Errorf("test[%d] - %q - obj ==> unexpected type. expected: %T actual: %T", idx, input, object.Null{}, obj)
		return false
	}

	return true
}

func testError(t *testing.T, idx int, input string, obj object.Object, message string) bool {
	result, ok := obj.(*object.Error)
	if !ok {
		t.Errorf("test[%d] - %q - obj ==> unexpected type. expected: %T actual: %T", idx, input, object.Error{}, obj)
		return false
	}

	if message != result.Message {
		t.Errorf("test[%d] - %q - result.Message ==> expected: %q actual: %q", idx, input, message, result.Message)
		return false
	}

	return true
}

func testFunction(t *testing.T, idx int, input string, obj object.Object, params []string, body string) bool {
	result, ok := obj.(*object.Function)
	if !ok {
		t.Errorf("test[%d] - %q - obj ==> unexpected type. expected: %T actual: %T", idx, input, object.Function{}, obj)
		return false
	}

	if len(params) != len(result.Parameters) {
		t.Errorf("test[%d] - %q - len(result.Parameters) ==> expected: %d actual: %d", idx, input, len(params), len(result.Parameters))
		return false
	}

	for i, param := range result.Parameters {
		if params[i] != param.String() {
			t.Errorf("test[%d] - %q - result.Parameters[%d].String() ==> expected: %q actual: %q", idx, input, i, params[i], param)
			return false
		}
	}

	actual := result.Body.String()
	if body != actual {
		t.Errorf("test[%d] - %q - result.Body.String() ==> expected: %q actual: %q", idx, input, body, actual)
		return false
	}

	return true
}

func testString(t *testing.T, idx int, input string, obj object.Object, value string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("test[%d] - %q - obj ==> unexpected type. expected: %T actual: %T", idx, input, object.String{}, obj)
		return false
	}

	if value != result.Value {
		t.Errorf("test[%d] - %q - result.Value ==> expected: %q actual: %q", idx, input, value, result.Value)
		return false
	}

	return true
}

func testArray(t *testing.T, idx int, input string, obj object.Object, tests []ObjectTest) bool {
	result, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("test[%d] - %q - obj ==> unexpected type. expected: %T actual: %T", idx, input, object.Array{}, obj)
		return false
	}

	if len(tests) != len(result.Elements) {
		t.Errorf("test[%d] - %q - len(result.Elements) ==> expected: %d actual: %d", idx, input, len(tests), len(result.Elements))
		return false
	}

	for i, elem := range result.Elements {
		if !testObject(t, idx, input, elem, tests[i]) {
			return false
		}
	}

	return true
}

func testHash(t *testing.T, idx int, input string, obj object.Object, tests map[object.HashKey]ObjectTest) bool {
	result, ok := obj.(*object.Hash)
	if !ok {
		t.Errorf("test[%d] - %q - obj ==> unexpected type. expected: %T actual: %T", idx, input, object.Hash{}, obj)
		return false
	}

	if len(tests) != len(result.Pairs) {
		t.Errorf("test[%d] - %q - len(result.Pairs) ==> expected: %d actual: %d", idx, input, len(tests), len(result.Pairs))
		return false
	}

	for key, value := range result.Pairs {
		test, ok := tests[key]
		if !ok {
			t.Errorf("test[%d] - %q - %T(%+v) ==> expected: not <nil>", idx, input, value.Key, value.Key)
			return false
		}

		if !testObject(t, idx, input, value.Value, test) {
			return false
		}
	}

	return true
}

func testQuote(t *testing.T, idx int, input string, obj object.Object, expected string) bool {
	result, ok := obj.(*object.Quote)
	if !ok {
		t.Errorf("test[%d] - %q - obj ==> unexpected type. expected: %T actual: %T", idx, input, object.Quote{}, obj)
		return false
	}

	actual := result.Node.String()
	if expected != actual {
		t.Errorf("test[%d] - %q - result.Node.String() ==> expected: %q actual: %q", idx, input, expected, actual)
		return false
	}

	return true
}
