package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
	ZERO  = &object.Integer{Value: 0}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ReturnStatement:
		return evalReturnStatement(node)
	case *ast.ExpressionStatement:
		return evalExpressionStatement(node)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node)
	case *ast.Boolean:
		return evalBoolean(node)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node)
	case *ast.InfixExpression:
		return evalInfixExpression(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	}

	return nil
}

func evalProgram(node *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range node.Statements {
		result = Eval(stmt)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

func evalReturnStatement(node *ast.ReturnStatement) object.Object {
	return &object.ReturnValue{
		Value: Eval(node.ReturnValue),
	}
}

func evalExpressionStatement(node *ast.ExpressionStatement) object.Object {
	return Eval(node.Expression)
}

func evalBlockStatement(node *ast.BlockStatement) object.Object {
	var result object.Object

	for _, stmt := range node.Statements {
		result = Eval(stmt)

		if result != nil && result.Type() == object.RETURN_VALUE_OBJECT {
			return result
		}
	}

	return result
}

func evalIntegerLiteral(node *ast.IntegerLiteral) object.Object {
	return toIntegerObject(node.Value)
}

func evalBoolean(node *ast.Boolean) object.Object {
	return toBooleanObject(node.Value)
}

func evalPrefixExpression(node *ast.PrefixExpression) object.Object {
	right := Eval(node.Right)

	switch {
	case right.Type() == object.INTEGER_OBJECT:
		return evalIntegerPrefixExpression(node.Operator, right.(*object.Integer))
	case right.Type() == object.BOOLEAN_OBJECT:
		return evalBooleanPrefixExpression(node.Operator, right.(*object.Boolean))
	case right.Type() == object.NULL_OBJECT:
		return evalNullPrefixExpression(node.Operator)
	default:
		return NULL
	}
}

func evalIntegerPrefixExpression(operator string, right *object.Integer) object.Object {
	switch operator {
	case "!":
		return toBooleanObject(right == ZERO)
	case "-":
		return toIntegerObject(-right.Value)
	default:
		return NULL
	}
}

func evalBooleanPrefixExpression(operator string, right *object.Boolean) object.Object {
	switch operator {
	case "!":
		return toBooleanObject(right == FALSE)
	default:
		return NULL
	}
}

func evalNullPrefixExpression(operator string) object.Object {
	switch operator {
	case "!":
		return TRUE
	default:
		return NULL
	}
}

func evalInfixExpression(node *ast.InfixExpression) object.Object {
	left := Eval(node.Left)
	right := Eval(node.Right)

	switch {
	case left.Type() == object.INTEGER_OBJECT && right.Type() == object.INTEGER_OBJECT:
		return evalIntegerInfixExpression(node.Operator, left.(*object.Integer), right.(*object.Integer))
	case left.Type() == object.BOOLEAN_OBJECT && right.Type() == object.BOOLEAN_OBJECT:
		return evalBooleanInfixExpression(node.Operator, left.(*object.Boolean), right.(*object.Boolean))
	case left.Type() == object.NULL_OBJECT && right.Type() == object.NULL_OBJECT:
		return evalNullInfixExpression(node.Operator)
	default:
		return NULL
	}
}

func evalIntegerInfixExpression(operator string, left, right *object.Integer) object.Object {
	switch operator {
	case "+":
		return toIntegerObject(left.Value + right.Value)
	case "-":
		return toIntegerObject(left.Value - right.Value)
	case "*":
		return toIntegerObject(left.Value * right.Value)
	case "/":
		return toIntegerObject(left.Value / right.Value)
	case "<":
		return toBooleanObject(left.Value < right.Value)
	case ">":
		return toBooleanObject(left.Value > right.Value)
	case "==":
		return toBooleanObject(left.Value == right.Value)
	case "!=":
		return toBooleanObject(left.Value != right.Value)
	default:
		return NULL
	}
}

func evalBooleanInfixExpression(operator string, left, right *object.Boolean) object.Object {
	switch operator {
	case "==":
		return toBooleanObject(left == right)
	case "!=":
		return toBooleanObject(left != right)
	default:
		return NULL
	}
}

func evalNullInfixExpression(operator string) object.Object {
	switch operator {
	case "==":
		return TRUE
	case "!=":
		return FALSE
	default:
		return NULL
	}
}

func evalIfExpression(node *ast.IfExpression) object.Object {
	condition := Eval(node.Condition)

	if isTruthy(condition) {
		return Eval(node.Consequence)
	} else if node.Alternative != nil {
		return Eval(node.Alternative)
	}

	return NULL
}

func toIntegerObject(val int64) object.Object {
	if val == 0 {
		return ZERO
	}
	return &object.Integer{
		Value: val,
	}
}

func toBooleanObject(val bool) object.Object {
	if val {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case ZERO, FALSE, NULL:
		return false
	default:
		return true
	}
}
