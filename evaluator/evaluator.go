package evaluator

import (
	"github.com/cijin/go-interpreter/ast"
	"github.com/cijin/go-interpreter/object"
	"github.com/cijin/go-interpreter/token"
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

func evalBangPrefixExpressionOperator(operand object.Object) object.Object {
	switch operand {
	case TRUE:
		return FALSE

	case FALSE:
		return TRUE

	case NULL:
		return TRUE

	default:
		return FALSE
	}
}

func evalMinusPrefixExpressionOperator(operand object.Object) object.Object {
	if operand.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := operand.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case token.MINUS:
		return evalMinusPrefixExpressionOperator(right)

	case token.BANG:
		return evalBangPrefixExpressionOperator(right)

	default:
		return NULL
	}
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

	case *ast.PrefixExpression:
		right := Eval(n.Right)
		return evalPrefixExpression(n.Operator, right)
	}

	return nil
}
