package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

var (
	ZERO  = &object.Integer{Value: 0}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.LetStatement:
		return evalLetStatement(node, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	case *ast.ExpressionStatement:
		return evalExpressionStatement(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node)
	case *ast.Boolean:
		return evalBoolean(node)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.FunctionLiteral:
		return evalFunctionLiteral(node, env)
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	}

	return nil
}

func evalProgram(node *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range node.Statements {
		result = Eval(stmt, env)

		if isError(result) {
			return result
		}

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

func evalLetStatement(node *ast.LetStatement, env *object.Environment) object.Object {
	value := Eval(node.Value, env)

	if isError(value) {
		return value
	}

	env.Set(node.Name.Value, value)

	return nil
}

func evalReturnStatement(node *ast.ReturnStatement, env *object.Environment) object.Object {
	value := Eval(node.ReturnValue, env)

	if isError(value) {
		return value
	}

	return &object.ReturnValue{
		Value: value,
	}
}

func evalExpressionStatement(node *ast.ExpressionStatement, env *object.Environment) object.Object {
	return Eval(node.Expression, env)
}

func evalBlockStatement(node *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range node.Statements {
		result = Eval(stmt, env)

		if isError(result) {
			return result
		}

		if result != nil && result.Type() == object.RETURN_VALUE_OBJECT {
			return result
		}
	}

	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if obj, ok := env.Get(node.Value); ok {
		return obj
	}

	return toErrorObject("unknown identifier: %s", node.Value)
}

func evalIntegerLiteral(node *ast.IntegerLiteral) object.Object {
	return toIntegerObject(node.Value)
}

func evalBoolean(node *ast.Boolean) object.Object {
	return toBooleanObject(node.Value)
}

func evalPrefixExpression(node *ast.PrefixExpression, env *object.Environment) object.Object {
	right := Eval(node.Right, env)

	if isError(right) {
		return right
	}

	switch {
	case right.Type() == object.INTEGER_OBJECT:
		return evalIntegerPrefixExpression(node.Operator, right.(*object.Integer))
	case right.Type() == object.BOOLEAN_OBJECT:
		return evalBooleanPrefixExpression(node.Operator, right.(*object.Boolean))
	case right.Type() == object.NULL_OBJECT:
		return evalNullPrefixExpression(node.Operator)
	default:
		return toErrorObject("unknown operation: %s%s", node.Operator, right.Type())
	}
}

func evalIntegerPrefixExpression(operator string, right *object.Integer) object.Object {
	switch operator {
	case "!":
		return toBooleanObject(right == ZERO)
	case "-":
		return toIntegerObject(-right.Value)
	default:
		return toErrorObject("unknown operation: %s%s", operator, object.INTEGER_OBJECT)
	}
}

func evalBooleanPrefixExpression(operator string, right *object.Boolean) object.Object {
	switch operator {
	case "!":
		return toBooleanObject(right == FALSE)
	default:
		return toErrorObject("unknown operation: %s%s", operator, object.BOOLEAN_OBJECT)
	}
}

func evalNullPrefixExpression(operator string) object.Object {
	switch operator {
	case "!":
		return TRUE
	default:
		return toErrorObject("unknown operation: %s%s", operator, object.NULL_OBJECT)
	}
}

func evalInfixExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)

	if isError(left) {
		return left
	}

	right := Eval(node.Right, env)

	if isError(right) {
		return right
	}

	switch {
	case left.Type() == object.INTEGER_OBJECT && right.Type() == object.INTEGER_OBJECT:
		return evalIntegerInfixExpression(node.Operator, left.(*object.Integer), right.(*object.Integer))
	case left.Type() == object.BOOLEAN_OBJECT && right.Type() == object.BOOLEAN_OBJECT:
		return evalBooleanInfixExpression(node.Operator, left.(*object.Boolean), right.(*object.Boolean))
	case left.Type() == object.NULL_OBJECT && right.Type() == object.NULL_OBJECT:
		return evalNullInfixExpression(node.Operator)
	default:
		return toErrorObject("unknown operation: %s %s %s", left.Type(), node.Operator, right.Type())
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
		return toErrorObject("unknown operation: %s %s %s", object.INTEGER_OBJECT, operator, object.INTEGER_OBJECT)
	}
}

func evalBooleanInfixExpression(operator string, left, right *object.Boolean) object.Object {
	switch operator {
	case "==":
		return toBooleanObject(left == right)
	case "!=":
		return toBooleanObject(left != right)
	default:
		return toErrorObject("unknown operation: %s %s %s", object.BOOLEAN_OBJECT, operator, object.BOOLEAN_OBJECT)
	}
}

func evalNullInfixExpression(operator string) object.Object {
	switch operator {
	case "==":
		return TRUE
	case "!=":
		return FALSE
	default:
		return toErrorObject("unknown operation: %s %s %s", object.NULL_OBJECT, operator, object.NULL_OBJECT)
	}
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	}

	return NULL
}

func evalFunctionLiteral(node *ast.FunctionLiteral, env *object.Environment) object.Object {
	return &object.Function{
		Parameters: node.Parameters,
		Body:       node.Body,
		Env:        env,
	}
}

func evalCallExpression(node *ast.CallExpression, env *object.Environment) object.Object {
	function := Eval(node.Function, env)

	if isError(function) {
		return function
	}

	args := []object.Object{}

	for _, arg := range node.Arguements {
		result := Eval(arg, env)

		if isError(result) {
			return result
		}

		args = append(args, result)
	}

	fn, ok := function.(*object.Function)
	if !ok {
		return toErrorObject("unknown call expression type: %s", function.Type())
	}

	enclosed := object.NewEnclosedEnvironment(fn.Env)

	for i, parameter := range fn.Parameters {
		enclosed.Set(parameter.Value, args[i])
	}

	result := Eval(fn.Body, enclosed)

	if returnValue, ok := result.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return result
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

func toErrorObject(format string, args ...any) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, args...),
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case ZERO, FALSE, NULL:
		return false
	default:
		return true
	}
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR_OBJECT
}
