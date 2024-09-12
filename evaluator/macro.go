package evaluator

import (
	"github.com/eugene-whitaker/monkey/ast"
	"github.com/eugene-whitaker/monkey/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}

	for i, stmt := range program.Statements {
		letStmt, ok := stmt.(*ast.LetStatement)
		if !ok {
			continue
		}

		macroExpr, ok := letStmt.Value.(*ast.MacroExpression)
		if !ok {
			continue
		}

		macro := &object.Macro{
			Parameters: macroExpr.Parameters,
			Body:       macroExpr.Body,
			Env:        env,
		}

		env.Set(letStmt.Name.Value, macro)
		definitions = append(definitions, i)
	}

	for i := len(definitions) - 1; i >= 0; i = i - 1 {
		idx := definitions[i]
		program.Statements = append(
			program.Statements[:idx],
			program.Statements[idx+1:]...,
		)
	}
}

func ExpandMacro(program *ast.Program, env *object.Environment) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		callExpr, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		ident, ok := callExpr.Function.(*ast.Identifier)
		if !ok {
			return node
		}

		obj, ok := env.Get(ident.Value)
		if !ok {
			return node
		}

		macro, ok := obj.(*object.Macro)
		if !ok {
			return node
		}

		args := []*object.Quote{}

		for _, arg := range callExpr.Arguments {
			args = append(args, &object.Quote{Node: arg})
		}

		enclosed := object.NewEnclosedEnvironment(macro.Env)

		for i, param := range macro.Parameters {
			enclosed.Set(param.Value, args[i])
		}

		eval := Eval(macro.Body, enclosed)

		quote, ok := eval.(*object.Quote)
		if !ok {
			return nil
		}

		return quote.Node
	})
}
