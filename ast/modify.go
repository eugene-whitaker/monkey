package ast

type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {
	case *Program:
		modifyProgram(node, modifier)
	case *LetStatement:
		modifyLetStatement(node, modifier)
	case *ReturnStatement:
		modifyReturnStatement(node, modifier)
	case *ExpressionStatement:
		modifyExpressionStatement(node, modifier)
	case *BlockStatement:
		modifyBlockStatement(node, modifier)
	case *FunctionLiteral:
		modifyFunctionLiteral(node, modifier)
	case *ArrayLiteral:
		modifyArrayLiteral(node, modifier)
	case *HashLiteral:
		modifyHashLiteral(node, modifier)
	case *PrefixExpression:
		modifyPrefixExpression(node, modifier)
	case *InfixExpression:
		modifyInfixExpression(node, modifier)
	case *IfExpression:
		modifyIfExpression(node, modifier)
	case *CallExpression:
		modifyCallExpression(node, modifier)
	case *IndexExpression:
		modifyIndexExpression(node, modifier)
	}

	return modifier(node)
}

func modifyProgram(node *Program, modifier ModifierFunc) {
	for i, stmt := range node.Statements {
		node.Statements[i], _ = Modify(stmt, modifier).(Statement)
	}
}

func modifyLetStatement(node *LetStatement, modifier ModifierFunc) {
	node.Name, _ = Modify(node.Name, modifier).(*Identifier)
	node.Value, _ = Modify(node.Value, modifier).(Expression)
}

func modifyReturnStatement(node *ReturnStatement, modifier ModifierFunc) {
	node.ReturnValue, _ = Modify(node.ReturnValue, modifier).(Expression)
}

func modifyExpressionStatement(node *ExpressionStatement, modifier ModifierFunc) {
	node.Expression, _ = Modify(node.Expression, modifier).(Expression)
}

func modifyBlockStatement(node *BlockStatement, modifier ModifierFunc) {
	for i, stmt := range node.Statements {
		node.Statements[i], _ = Modify(stmt, modifier).(Statement)
	}
}

func modifyFunctionLiteral(node *FunctionLiteral, modifier ModifierFunc) {
	for i, ident := range node.Parameters {
		node.Parameters[i], _ = Modify(ident, modifier).(*Identifier)
	}
	node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)
}

func modifyArrayLiteral(node *ArrayLiteral, modifier ModifierFunc) {
	for i, elem := range node.Elements {
		node.Elements[i], _ = Modify(elem, modifier).(Expression)
	}
}

func modifyHashLiteral(node *HashLiteral, modifier ModifierFunc) {
	pairs := make(map[Expression]Expression)
	for key, value := range node.Pairs {
		key, _ = Modify(key, modifier).(Expression)
		value, _ = Modify(value, modifier).(Expression)

		pairs[key] = value
	}
	node.Pairs = pairs
}

func modifyPrefixExpression(node *PrefixExpression, modifier ModifierFunc) {
	node.Right, _ = Modify(node.Right, modifier).(Expression)
}

func modifyInfixExpression(node *InfixExpression, modifier ModifierFunc) {
	node.Left, _ = Modify(node.Left, modifier).(Expression)
	node.Right, _ = Modify(node.Right, modifier).(Expression)
}

func modifyIfExpression(node *IfExpression, modifier ModifierFunc) {
	node.Condition, _ = Modify(node.Condition, modifier).(Expression)
	node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStatement)
	if node.Alternative != nil {
		node.Alternative, _ = Modify(node.Alternative, modifier).(*BlockStatement)
	}
}

func modifyCallExpression(node *CallExpression, modifier ModifierFunc) {
	node.Function, _ = Modify(node.Function, modifier).(Expression)
	for i, arg := range node.Arguments {
		node.Arguments[i], _ = Modify(arg, modifier).(Expression)
	}
}

func modifyIndexExpression(node *IndexExpression, modifier ModifierFunc) {
	node.Struct, _ = Modify(node.Struct, modifier).(Expression)
	node.Index, _ = Modify(node.Index, modifier).(Expression)
}
