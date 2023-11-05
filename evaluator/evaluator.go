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

func evalProgram(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)

		if r, ok := result.(*object.ReturnValue); ok {
			return r.Value
		}
	}

	return result
}

/*
 * This method is seperated to ensure that the innermost block statement
 * return gets returned and the execution stops. If you unwrap at the block
 * statement level (return obj), execution keeps continuing as `evalProgram`
 * executes the next statement.
 *
 * Prior to seperating evalBlockStatements & evalProgram, execution for statements
 * was handled by a generic evalStatements
 */
func evalBlockStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)

		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
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

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}

	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)

	default:
		return NULL
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	case operator == "==":
		return nativeBoolToBooleanObject(left == right)

	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	default:
		return NULL
	}
}

func isTruthy(condition object.Object) bool {
	switch condition {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		return true
	}
}

func evalIfExpression(n *ast.IfExpression) object.Object {
	condition := Eval(n.Condition)

	if isTruthy(condition) {
		return Eval(n.Consequence)
	} else if n.Alternative != nil {
		return Eval(n.Alternative)
	}

	return NULL
}

func Eval(node ast.Node) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n.Statements)

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

	case *ast.InfixExpression:
		left := Eval(n.Left)
		right := Eval(n.Right)

		return evalInfixExpression(n.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatements(n.Statements)

	case *ast.IfExpression:
		return evalIfExpression(n)

	case *ast.ReturnStatement:
		v := Eval(n.ReturnValue)
		return &object.ReturnValue{Value: v}

	}

	return nil
}
