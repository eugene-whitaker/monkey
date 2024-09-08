package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
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
				"unknown identifier: ident",
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
			"fn(x) { x; }(5)",
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
	}

	for i, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		env := object.NewEnvironment()
		eval := Eval(program, env)

		if !testObject(t, i, tt.input, eval, tt.test) {
			return
		}
	}
}

func testObject(t *testing.T, index int, input string, obj object.Object, test ObjectTest) bool {
	switch test := test.(type) {
	case IntegerTest:
		return testInteger(t, index, input, obj, int64(test))
	case BooleanTest:
		return testBoolean(t, index, input, obj, bool(test))
	case NullTest:
		return testNull(t, index, input, obj)
	case ErrorTest:
		return testError(t, index, input, obj, test.message)
	case FunctionTest:
		return testFunction(t, index, input, obj, test.parameters, test.body)
	}
	t.Errorf("test[%d] - %q ==> unexpected type. actual: %T", index, input, test)
	return false
}

func testInteger(t *testing.T, index int, input string, obj object.Object, value int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("test[%d] - %q - obj ==> unexpected type. expected: %T actual: %T", index, input, object.Integer{}, obj)
		return false
	}

	if value != result.Value {
		t.Errorf("test[%d] - %q - result.Value ==> expected: %d actual: %d", index, input, value, result.Value)
		return false
	}

	return true
}

func testBoolean(t *testing.T, index int, input string, obj object.Object, value bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("test[%d] - %q - obj ==> unexpected type. expected: %T actual: %T", index, input, object.Boolean{}, obj)
		return false
	}

	if value != result.Value {
		t.Errorf("test[%d] - %q - result.Value ==> expected: %t actual: %t", index, input, value, result.Value)
		return false
	}

	return true
}

func testNull(t *testing.T, index int, input string, obj object.Object) bool {
	if NULL != obj {
		t.Errorf("test[%d] - %q - obj ==> unexpected type. expected: %T actual: %T", index, input, object.Null{}, obj)
		return false
	}

	return true
}

func testError(t *testing.T, index int, input string, obj object.Object, msg string) bool {
	result, ok := obj.(*object.Error)
	if !ok {
		t.Errorf("test[%d] - %q - obj ==> unexpected type. expected: %T actual: %T", index, input, object.Error{}, obj)
		return false
	}

	if msg != result.Message {
		t.Errorf("test[%d] - %q - result.Value ==> expected: %q actual: %q", index, input, msg, result.Message)
		return false
	}

	return true
}

func testFunction(t *testing.T, index int, input string, obj object.Object, parameters []string, body string) bool {
	result, ok := obj.(*object.Function)
	if !ok {
		t.Errorf("test[%d] - %q - obj ==> unexpected type. expected: %T actual: %T", index, input, object.Function{}, obj)
		return false
	}

	if len(parameters) != len(result.Parameters) {
		t.Errorf("test[%d] - %q - len(result.Parameters) ==> expected: %d actual: %d", index, input, len(parameters), len(result.Parameters))
		return false
	}

	for i, parameter := range result.Parameters {
		if parameters[i] != parameter.String() {
			t.Errorf("test[%d] - %q - result.Parameters[%d].String() ==> expected: %q actual: %q", index, input, i, parameters[i], parameter)
			return false
		}
	}

	actual := result.Body.String()
	if body != actual {
		t.Errorf("test[%d] - %q - result.Body.String() ==> expected: %q actual: %q", index, input, body, actual)
		return false
	}

	return true
}
