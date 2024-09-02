package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		value int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testIntegerObject(t, evaluated, tt.value) {
			return
		}
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		value bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testBooleanObject(t, evaluated, tt.value) {
			return
		}
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input string
		value bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testBooleanObject(t, evaluated, tt.value) {
			return
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, value int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("obj ==> unexpected type. expected: %T actual: %T", &object.Integer{}, obj)
		return false
	}

	if value != result.Value {
		t.Errorf("result.Value ==> expected: %d actual: %d", value, result.Value)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, value bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("obj ==> unexpected type. expected: %T actual: %T", &object.Boolean{}, obj)
		return false
	}

	if value != result.Value {
		t.Errorf("result.Value ==> expected: %t actual: %t", value, result.Value)
		return false
	}

	return true
}
