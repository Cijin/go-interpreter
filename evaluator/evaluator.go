package evaluator

import (
	"github.com/cijin/go-interpreter/ast"
	"github.com/cijin/go-interpreter/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}

func nativeBoolToBooleanObject(v bool) object.Object {
	if v {
		return TRUE
	}

	return FALSE
}

func Eval(node ast.Node) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalStatements(n.Statements)

	case *ast.ExpressionStatement:
		return Eval(n.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: n.Value,
		}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(n.Value)
	}

	return nil
}
