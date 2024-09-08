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

func TestEval(t *testing.T) {
	tests := []struct {
		input string
		test  ObjectTest
	}{
		{
			"5",
			IntegerTest(5),
		},
		{
			"10",
			IntegerTest(10),
		},
		{
			"-5",
			IntegerTest(-5),
		},
		{
			"-10",
			IntegerTest(-10),
		},
		{
			"5 + 5 + 5 + 5 - 10",
			IntegerTest(10),
		},
		{
			"2 * 2 * 2 * 2 * 2",
			IntegerTest(32),
		},
		{
			"-50 + 100 + -50",
			IntegerTest(0),
		},
		{
			"5 * 2 + 10",
			IntegerTest(20),
		},
		{
			"5 + 2 * 10",
			IntegerTest(25),
		},
		{
			"20 + 2 * -10",
			IntegerTest(0),
		},
		{
			"50 / 2 * 2 + 10",
			IntegerTest(60),
		},
		{
			"2 * (5 + 10)",
			IntegerTest(30),
		},
		{
			"3 * 3 * 3 + 10",
			IntegerTest(37),
		},
		{
			"3 * (3 * 3) + 10",
			IntegerTest(37),
		},
		{
			"(5 + 10 * 2 + 15 / 3) * 2 + -10",
			IntegerTest(50),
		},
		{
			"true",
			BooleanTest(true),
		},
		{
			"false",
			BooleanTest(false),
		},
		{
			"true == true",
			BooleanTest(true),
		},
		{
			"false == false",
			BooleanTest(true),
		},
		{
			"true == false",
			BooleanTest(false),
		},
		{
			"true != false",
			BooleanTest(true),
		},
		{
			"false != true",
			BooleanTest(true),
		},
		{
			"(1 < 2) == true",
			BooleanTest(true),
		},
		{
			"(1 < 2) == false",
			BooleanTest(false),
		},
		{
			"(1 > 2) == true",
			BooleanTest(false),
		},
		{
			"(1 > 2) == false",
			BooleanTest(true),
		},
		{
			"!true",
			BooleanTest(false),
		},
		{
			"1 < 2",
			BooleanTest(true),
		},
		{
			"1 > 2",
			BooleanTest(false),
		},
		{
			"1 < 1",
			BooleanTest(false),
		},
		{
			"1 > 1",
			BooleanTest(false),
		},
		{
			"1 == 1",
			BooleanTest(true),
		},
		{
			"1 != 1",
			BooleanTest(false),
		},
		{
			"1 == 2",
			BooleanTest(false),
		},
		{
			"1 != 2",
			BooleanTest(true),
		},
		{
			"!false",
			BooleanTest(true),
		},
		{
			"!5",
			BooleanTest(false),
		},
		{
			"!!true",
			BooleanTest(true),
		},
		{
			"!!false",
			BooleanTest(false),
		},
		{
			"!!5",
			BooleanTest(true),
		},
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		program := p.ParseProgram()
		eval := Eval(program)

		if !testObject(t, i, eval, tt.test) {
			return
		}
	}
}

func testObject(t *testing.T, index int, obj object.Object, test ObjectTest) bool {
	switch test := test.(type) {
	case IntegerTest:
		return testInteger(t, index, obj, int64(test))
	case BooleanTest:
		return testBoolean(t, index, obj, bool(test))
	}
	t.Errorf("test[%d] - test ==> unexpected type. actual: %T", index, test)
	return false
}

func testInteger(t *testing.T, index int, obj object.Object, value int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("test[%d] - obj ==> unexpected type. expected: %T actual: %T", index, object.Integer{}, obj)
		return false
	}

	if value != result.Value {
		t.Errorf("test[%d] - result.Value ==> expected: %q actual: %q", index, value, result.Value)
		return false
	}

	return true
}

func testBoolean(t *testing.T, index int, obj object.Object, value bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("test[%d] - obj ==> unexpected type. expected: %T actual: %T", index, object.Boolean{}, obj)
		return false
	}

	if value != result.Value {
		t.Errorf("test[%d] - result.Value ==> expected: %t actual: %t", index, value, result.Value)
		return false
	}

	return true
}
