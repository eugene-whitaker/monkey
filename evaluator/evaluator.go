package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
	"monkey/token"
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
	EMPTY_HASH = &object.Hash{
		Pairs: make(map[object.HashKey]object.HashPair),
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
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
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

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for key, value := range node.Pairs {
		key := Eval(key, env)

		if isError(key) {
			return key
		}

		hashkey, ok := key.(object.Hashable)
		if !ok {
			return toErrorObject("invalid type: %s is not hashable", key.Type())
		}

		value := Eval(value, env)

		if isError(value) {
			return value
		}

		pairs[hashkey.HashKey()] = object.HashPair{Key: key, Value: value}
	}

	return toHashObject(pairs)
}

func evalPrefixExpression(node *ast.PrefixExpression, env *object.Environment) object.Object {
	right := Eval(node.Right, env)

	if isError(right) {
		return right
	}

	switch right.Type() {
	case object.INTEGER_OBJECT:
		return evalIntegerPrefixExpression(node.Operator, right.(*object.Integer))
	case object.BOOLEAN_OBJECT:
		return evalBooleanPrefixExpression(node.Operator, right.(*object.Boolean))
	case object.NULL_OBJECT:
		return evalNullPrefixExpression(node.Operator)
	case object.STRING_OBJECT:
		return evalStringPrefixExpression(node.Operator, right.(*object.String))
	case object.ARRAY_OBJECT:
		return evalArrayPrefixExpression(node.Operator, right.(*object.Array))
	case object.HASH_OBJECT:
		return evalHashPrefixExpression(node.Operator, right.(*object.Hash))
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

func evalArrayPrefixExpression(operator string, right *object.Array) object.Object {
	switch operator {
	case "!":
		return toBooleanObject(right == EMPTY_ARRAY)
	default:
		return toErrorObject("unknown operation: %s%s", operator, object.ARRAY_OBJECT)
	}
}

func evalHashPrefixExpression(operator string, right *object.Hash) object.Object {
	switch operator {
	case "!":
		return toBooleanObject(right == EMPTY_HASH)
	default:
		return toErrorObject("unknown operation: %s%s", operator, object.HASH_OBJECT)
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
	case left.Type() == object.FUNCTION_OBJECT && right.Type() == object.FUNCTION_OBJECT:
		return evalFunctionInfixExpression(node.Operator, left.(*object.Function), right.(*object.Function))
	case left.Type() == object.STRING_OBJECT && right.Type() == object.STRING_OBJECT:
		return evalStringInfixExpression(node.Operator, left.(*object.String), right.(*object.String))
	case left.Type() == object.ARRAY_OBJECT && right.Type() == object.ARRAY_OBJECT:
		return evalArrayInfixExpression(node.Operator, left.(*object.Array), right.(*object.Array))
	case left.Type() == object.HASH_OBJECT && right.Type() == object.HASH_OBJECT:
		return evalHashInfixExpression(node.Operator, left.(*object.Hash), right.(*object.Hash))
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

func evalFunctionInfixExpression(operator string, left, right *object.Function) object.Object {
	switch operator {
	case "==":
		return toBooleanObject(left == right)
	case "!=":
		return toBooleanObject(left != right)
	default:
		return toErrorObject("unknown operation: %s %s %s", object.FUNCTION_OBJECT, operator, object.FUNCTION_OBJECT)
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

func evalArrayInfixExpression(operator string, left, right *object.Array) object.Object {
	switch operator {
	case "==":
		return toBooleanObject(left == right)
	case "!=":
		return toBooleanObject(left != right)
	default:
		return toErrorObject("unknown operation: %s %s %s", object.ARRAY_OBJECT, operator, object.ARRAY_OBJECT)
	}
}

func evalHashInfixExpression(operator string, left, right *object.Hash) object.Object {
	switch operator {
	case "==":
		return toBooleanObject(left == right)
	case "!=":
		return toBooleanObject(left != right)
	default:
		return toErrorObject("unknown operation: %s %s %s", object.HASH_OBJECT, operator, object.HASH_OBJECT)
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
	if node.Function.TokenLexeme() == "quote" {
		return toQuoteObject(node.Arguments[0], env)
	}

	function := Eval(node.Function, env)

	if isError(function) {
		return function
	}

	args := []object.Object{}

	for _, arg := range node.Arguments {
		result := Eval(arg, env)

		if isError(result) {
			return result
		}

		args = append(args, result)
	}

	switch fn := function.(type) {
	case *object.Function:
		return evalFunctionCallExpression(fn, args)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return toErrorObject("unknown operation: %s()", function.Type())
	}
}

func evalFunctionCallExpression(fn *object.Function, args []object.Object) object.Object {
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

func evalUnquoteCallExpression(node ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(node, func(node ast.Node) ast.Node {
		callExpr, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		if callExpr.Function.TokenLexeme() != "unquote" {
			return node
		}

		if len(callExpr.Arguments) != 1 {
			return node
		}

		return toASTNode(Eval(callExpr.Arguments[0], env))
	})
}

func evalIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object {
	indexable := Eval(node.Struct, env)

	if isError(indexable) {
		return indexable
	}

	index := Eval(node.Index, env)

	if isError(index) {
		return index
	}

	switch {
	case indexable.Type() == object.ARRAY_OBJECT && index.Type() == object.INTEGER_OBJECT:
		return evalArrayIndexExpression(indexable.(*object.Array), index.(*object.Integer))
	case indexable.Type() == object.HASH_OBJECT:
		return evalHashIndexExpression(indexable.(*object.Hash), index)
	default:
		return toErrorObject("unknown operation: %s[%s]", indexable.Type(), index.Type())
	}
}

func evalArrayIndexExpression(array *object.Array, index *object.Integer) object.Object {
	if index.Value >= 0 && index.Value < int64(len(array.Elements)) {
		return array.Elements[index.Value]
	}
	return NULL
}

func evalHashIndexExpression(hash *object.Hash, index object.Object) object.Object {
	hashable, ok := index.(object.Hashable)
	if !ok {
		return toErrorObject("unknown operation: HASH[%s]", index.Type())
	}

	if result, ok := hash.Pairs[hashable.HashKey()]; ok {
		return result.Value
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

func toHashObject(pairs map[object.HashKey]object.HashPair) object.Object {
	if len(pairs) == 0 {
		return EMPTY_HASH
	}
	return &object.Hash{
		Pairs: pairs,
	}
}

func toQuoteObject(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquoteCallExpression(node, env)
	return &object.Quote{
		Node: node,
	}
}

func toASTNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		return &ast.IntegerLiteral{
			Token: token.Token{
				Type:   token.INT,
				Lexeme: fmt.Sprintf("%d", obj.Value),
			},
			Value: obj.Value,
		}
	default:
		return nil
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case ZERO, FALSE, NULL, EMPTY_STRING, EMPTY_ARRAY, EMPTY_HASH:
		return false
	default:
		return true
	}
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR_OBJECT
}
