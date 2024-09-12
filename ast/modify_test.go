package ast

import (
	"reflect"
	"testing"
)

func TestModify(t *testing.T) {
	one := func() Expression { return &IntegerLiteral{Value: 1} }
	two := func() Expression { return &IntegerLiteral{Value: 2} }

	modifier := func(node Node) Node {
		integ, ok := node.(*IntegerLiteral)
		if !ok {
			return node
		}

		if integ.Value == 1 {
			integ.Value = 2
			return integ
		}

		return node
	}

	tests := []struct {
		input    Node
		expected Node
	}{
		{
			one(),
			two(),
		},
		{
			&Program{
				Statements: []Statement{
					&ExpressionStatement{
						Expression: one(),
					},
				},
			},
			&Program{
				Statements: []Statement{
					&ExpressionStatement{
						Expression: two(),
					},
				},
			},
		},
		{
			&InfixExpression{
				Left:     one(),
				Operator: "+",
				Right:    two(),
			},
			&InfixExpression{
				Left:     two(),
				Operator: "+",
				Right:    two(),
			},
		},
		{
			&PrefixExpression{
				Operator: "-",
				Right:    one(),
			},
			&PrefixExpression{
				Operator: "-",
				Right:    two(),
			},
		},
		{
			&IndexExpression{
				Struct: &ArrayLiteral{
					Elements: []Expression{
						one(),
					},
				},
				Index: one(),
			},
			&IndexExpression{
				Struct: &ArrayLiteral{
					Elements: []Expression{
						two(),
					},
				},
				Index: two(),
			},
		},
		{
			&IfExpression{
				Condition: one(),
				Consequence: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{
							Expression: one(),
						},
					},
				},
				Alternative: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{
							Expression: one(),
						},
					},
				},
			},
			&IfExpression{
				Condition: two(),
				Consequence: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{
							Expression: two(),
						},
					},
				},
				Alternative: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{
							Expression: two(),
						},
					},
				},
			},
		},
		{
			&ReturnStatement{
				ReturnValue: one(),
			},
			&ReturnStatement{
				ReturnValue: two(),
			},
		},
		{
			&LetStatement{
				Value: one(),
			},
			&LetStatement{
				Value: two(),
			},
		},
		{
			&FunctionLiteral{
				Parameters: []*Identifier{},
				Body: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{
							Expression: one(),
						},
					},
				},
			},
			&FunctionLiteral{
				Parameters: []*Identifier{},
				Body: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{
							Expression: two(),
						},
					},
				},
			},
		},
		{
			&ArrayLiteral{
				Elements: []Expression{
					one(),
					one(),
				},
			},
			&ArrayLiteral{
				Elements: []Expression{
					two(),
					two(),
				},
			},
		},
	}

	for i, test := range tests {
		actual := Modify(test.input, modifier)

		if !reflect.DeepEqual(test.expected, actual) {
			t.Errorf("test[%d] - Modify(test.input, modifier) ==> expected: %#v actual: %#v", i, test.expected, actual)
			continue
		}
	}

	input := &HashLiteral{
		Pairs: map[Expression]Expression{
			one(): one(),
			one(): one(),
		},
	}
	Modify(input, modifier)

	for key, val := range input.Pairs {
		key, _ := key.(*IntegerLiteral)
		if 2 != key.Value {
			t.Errorf("key.Value ==> expected: %d actual: %d", 2, key.Value)
		}

		val, _ := val.(*IntegerLiteral)
		if 2 != val.Value {
			t.Errorf("val.Value ==> expected: %d actual: %d", 2, val.Value)
		}
	}
}
