package evaluator

import (
	"github.com/smith-30/go-monkey/ast"
	"github.com/smith-30/go-monkey/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// statement
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	// detail expression
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}
