package ast

type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {
	case *Program:
		for i, stmt := range node.Statements {
			node.Statements[i], _ = Modify(stmt, modifier).(Statement)
		}
	case *LetStatement:
		node.Name, _ = Modify(node.Name, modifier).(*Identifier)
		node.Value, _ = Modify(node.Value, modifier).(Expression)
	case *ReturnStatement:
		node.ReturnValue, _ = Modify(node.ReturnValue, modifier).(Expression)
	case *ExpressionStatement:
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)
	case *BlockStatement:
		for i, stmt := range node.Statements {
			node.Statements[i], _ = Modify(stmt, modifier).(Statement)
		}
	case *FunctionLiteral:
		for i, ident := range node.Parameters {
			node.Parameters[i], _ = Modify(ident, modifier).(*Identifier)
		}
		node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)
	case *ArrayLiteral:
		for i, elem := range node.Elements {
			node.Elements[i], _ = Modify(elem, modifier).(Expression)
		}
	case *HashLiteral:
		pairs := make(map[Expression]Expression)
		for key, value := range node.Pairs {
			key, _ = Modify(key, modifier).(Expression)
			value, _ = Modify(value, modifier).(Expression)

			pairs[key] = value
		}
		node.Pairs = pairs
	case *PrefixExpression:
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *InfixExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *IfExpression:
		node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStatement)
		if node.Alternative != nil {
			node.Alternative, _ = Modify(node.Alternative, modifier).(*BlockStatement)
		}
	case *CallExpression:
		node.Function, _ = Modify(node.Function, modifier).(Expression)
		for i, arg := range node.Arguments {
			node.Arguments[i], _ = Modify(arg, modifier).(Expression)
		}
	case *IndexExpression:
		node.Struct, _ = Modify(node.Struct, modifier).(Expression)
		node.Index, _ = Modify(node.Index, modifier).(Expression)
	}

	return modifier(node)
}
