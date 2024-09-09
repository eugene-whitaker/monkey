package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

var (
	ZERO = &object.Integer{
		Value: 0,
	}
	TRUE = &object.Boolean{
		Value: true,
	}
	FALSE = &object.Boolean{
		Value: false,
	}
	NULL         = &object.Null{}
	EMPTY_STRING = &object.String{
		Value: "",
	}
	EMPTY_ARRAY = &object.Array{
		Elements: []object.Object{},
	}
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
	case *ast.BooleanLiteral:
		return evalBooleanLiteral(node)
	case *ast.FunctionLiteral:
		return evalFunctionLiteral(node, env)
	case *ast.StringLiteral:
		return evalStringLiteral(node)
	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, env)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	case *ast.IndexExpression:
		return evalIndexExpression(node, env)
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

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return toErrorObject("undefined reference: %s", node.Value)
}

func evalIntegerLiteral(node *ast.IntegerLiteral) object.Object {
	return toIntegerObject(node.Value)
}

func evalBooleanLiteral(node *ast.BooleanLiteral) object.Object {
	return toBooleanObject(node.Value)
}

func evalFunctionLiteral(node *ast.FunctionLiteral, env *object.Environment) object.Object {
	return &object.Function{
		Parameters: node.Parameters,
		Body:       node.Body,
		Env:        env,
	}
}

func evalStringLiteral(node *ast.StringLiteral) object.Object {
	return toStringObject(node.Value)
}

func evalArrayLiteral(node *ast.ArrayLiteral, env *object.Environment) object.Object {
	elems := []object.Object{}

	for _, elem := range node.Elements {
		result := Eval(elem, env)

		if isError(result) {
			return result
		}

		elems = append(elems, result)
	}

	return toArrayObject(elems)
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
	case right.Type() == object.STRING_OBJECT:
		return evalStringPrefixExpression(node.Operator, right.(*object.String))
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

func evalStringPrefixExpression(operator string, right *object.String) object.Object {
	switch operator {
	case "!":
		return toBooleanObject(right == EMPTY_STRING)
	default:
		return toErrorObject("unknown operation: %s%s", operator, object.STRING_OBJECT)
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
	case left.Type() == object.STRING_OBJECT && right.Type() == object.STRING_OBJECT:
		return evalStringInfixExpression(node.Operator, left.(*object.String), right.(*object.String))
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

func evalStringInfixExpression(operator string, left, right *object.String) object.Object {
	switch operator {
	case "+":
		return toStringObject(left.Value + right.Value)
	case "==":
		return toBooleanObject(left.Value == right.Value)
	case "!=":
		return toBooleanObject(left.Value != right.Value)
	default:
		return toErrorObject("unknown operation: %s %s %s", object.STRING_OBJECT, operator, object.STRING_OBJECT)
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

	switch fn := function.(type) {
	case *object.Function:
		enclosed := object.NewEnclosedEnvironment(fn.Env)

		for i, parameter := range fn.Parameters {
			enclosed.Set(parameter.Value, args[i])
		}

		result := Eval(fn.Body, enclosed)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}

		return result
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return toErrorObject("unknown operation: %s()", function.Type())
	}
}

func evalIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object {
	array := Eval(node.Array, env)

	if isError(array) {
		return array
	}

	index := Eval(node.Index, env)

	if isError(index) {
		return index
	}

	switch {
	case array.Type() == object.ARRAY_OBJECT && index.Type() == object.INTEGER_OBJECT:
		arr := array.(*object.Array)
		idx := index.(*object.Integer).Value

		if idx < 0 || idx >= int64(len(arr.Elements)) {
			return NULL
		}

		return arr.Elements[idx]
	default:
		return toErrorObject("unknown operation: %s[%s]", array.Type(), index.Type())
	}
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

func toErrorObject(format string, args ...any) object.Object {
	return &object.Error{
		Message: fmt.Sprintf(format, args...),
	}
}

func toStringObject(val string) object.Object {
	if val == "" {
		return EMPTY_STRING
	}
	return &object.String{
		Value: val,
	}
}

func toArrayObject(elems []object.Object) object.Object {
	if len(elems) == 0 {
		return EMPTY_ARRAY
	}
	return &object.Array{
		Elements: elems,
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case ZERO, FALSE, NULL, EMPTY_STRING:
		return false
	default:
		return true
	}
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR_OBJECT
}
