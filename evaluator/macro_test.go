package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

type EnvTest interface {
	ident()
}

type DefaultTest string

func (dt DefaultTest) ident() {}

type MacroTest struct {
	name       string
	parameters []string
	body       string
}

func (mt MacroTest) ident() {}

func TestDefineMacros(t *testing.T) {
	tests := []struct {
		input string
		tests []EnvTest
	}{
		{
			"let number = 1; let func = fn(x, y) { x + y; }; let sum = macro(x, y) { x + y; };",
			[]EnvTest{
				DefaultTest("number"),
				DefaultTest("func"),
				MacroTest{
					"sum",
					[]string{
						"x",
						"y",
					},
					"(x + y)",
				},
			},
		},
	}

	for i, test := range tests {
		if !testEnv(t, i, test.input, test.tests) {
			continue
		}
	}
}

func TestExpandMacros(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"let infix = macro() { quote(1 + 2); }; infix();",
			"(1 + 2)",
		},
		{
			"let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); }; reverse(2 + 2, 10 - 5);",
			"((10 - 5) - (2 + 2))",
		},
		{
			"let unless = macro(condition, consequence, alternative) { quote(if (!(unquote(condition))) { unquote(consequence); } else { unquote(alternative); }); }; unless(10 > 5, puts(\"false\"), puts(\"true\"));",
			"if (!(10 > 5)) puts(false) else puts(true)",
		},
	}

	for i, test := range tests {
		l := lexer.NewLexer(test.input)
		p := parser.NewParser(l)
		program := p.ParseProgram()

		env := object.NewEnvironment()
		DefineMacros(program, env)

		expanded := ExpandMacro(program, env)

		actual := expanded.String()
		if test.expected != actual {
			t.Errorf("test[%d] - %q - expanded.String() ==> expected: %q actual: %q", i, test.input, test.expected, actual)
			continue
		}
	}
}

func testEnv(t *testing.T, idx int, input string, tests []EnvTest) bool {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	env := object.NewEnvironment()
	DefineMacros(program, env)

	length := 0
	for _, test := range tests {
		if _, ok := test.(MacroTest); !ok {
			length = length + 1
		}
	}

	if length != len(program.Statements) {
		t.Errorf("test[%d] - %q - len(program.Statements) ==> expected: %d actual: %d", idx, input, length, len(program.Statements))
		return false
	}

	for _, test := range tests {
		switch test := test.(type) {
		case MacroTest:
			obj, ok := env.Get(test.name)
			if !ok {
				t.Errorf("test[%d] - %q - env.Get(test.name) ==> expected: not <nil>", idx, input)
				return false
			}

			return testMacro(t, idx, input, obj, test.parameters, test.body)
		default:
			t.Logf("test: %#v env: %#v", test, env)
			name, ok := test.(DefaultTest)
			if !ok {
				t.Errorf("test[%d] - %q - test.(DefaultTest) ==> unexpected type. expected: %T actual: %T", idx, input, DefaultTest(""), test)
				return false
			}

			obj, ok := env.Get(string(name))
			if ok {
				t.Errorf("test[%d] - %q - env.Get(string(name)) ==> expected: <nil> actual %#v", idx, input, obj)
				return false
			}
		}
	}

	return true
}

func testMacro(t *testing.T, idx int, input string, obj object.Object, params []string, body string) bool {
	result, ok := obj.(*object.Macro)
	if !ok {
		t.Errorf("test[%d] - %q - obj.(*object.Macro) ==> unexpected type. expected: %T actual: %T", idx, input, object.Macro{}, obj)
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
