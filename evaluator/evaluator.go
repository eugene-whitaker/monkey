package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return evalExpressionStatement(node)
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node)
	case *ast.Boolean:
		return evalBoolean(node)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node)
	}

	return nil
}

func evalProgram(node *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range node.Statements {
		result = Eval(stmt)
	}

	return result
}

func evalExpressionStatement(node *ast.ExpressionStatement) object.Object {
	return Eval(node.Expression)
}

func evalIntegerLiteral(node *ast.IntegerLiteral) object.Object {
	return &object.Integer{
		Value: node.Value,
	}
}

func evalBoolean(node *ast.Boolean) object.Object {
	if node.Value {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(node *ast.PrefixExpression) object.Object {
	right := Eval(node.Right)

	switch node.Operator {
	case "!":
		switch right.Type() {
		case object.INTEGER_OBJ:
			if right.(*object.Integer).Value != 0 {
				return FALSE
			}
		case object.BOOLEAN_OBJ:
			if right == TRUE {
				return FALSE
			}
		case object.NULL_OBJ:
		default:
			return NULL
		}
		return TRUE
	case "-":
		if right.Type() != object.INTEGER_OBJ {
			return NULL
		}

		value := right.(*object.Integer).Value
		return &object.Integer{
			Value: -value,
		}
	default:
		return NULL
	}
}
