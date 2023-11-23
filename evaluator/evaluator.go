package evaluator

import (
	"fmt"

	"github.com/cijin/go-interpreter/ast"
	"github.com/cijin/go-interpreter/object"
	"github.com/cijin/go-interpreter/token"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("too many args for len, expected=1, got=%d", len(args))}
			}

			switch a := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(a.Value))}

			default:
				return &object.Error{Message: fmt.Sprintf("invalid arg type for len, expected=STRING, got=%s", a.Type())}
			}
		},
	},
}

func newErrorf(format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

func evalProgram(stmts []ast.Statement, env *object.Enviornment) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value

		case *object.Error:
			return result
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
func evalBlockStatements(stmts []ast.Statement, env *object.Enviornment) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)

		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
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
		return newErrorf("operator '-' not defined on %s", operand.Type())
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
		return newErrorf("unknown operator: %s", operator)
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
		return newErrorf("unknown operator: %s", operator)
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftValue + rightValue}
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)

	default:
		return newErrorf("operartor %s not supported on type string", operator)
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newErrorf("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)

	case operator == "==":
		return nativeBoolToBooleanObject(left == right)

	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	default:
		return newErrorf("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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

func evalIfExpression(n *ast.IfExpression, env *object.Enviornment) object.Object {
	condition := Eval(n.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(n.Consequence, env)
	} else if n.Alternative != nil {
		return Eval(n.Alternative, env)
	}

	return NULL
}

func evalIdentifier(ident *ast.Identifier, env *object.Enviornment) object.Object {
	i, ok := env.Get(ident.Value)
	if ok {
		return i
	}

	i, ok = builtins[ident.Value]
	if ok {
		return i
	}

	return newErrorf("identifier is undefined: %s", ident.Value)
}

func evalExpressions(exprs []ast.Expression, env *object.Enviornment) []object.Object {
	var result []object.Object

	for _, arg := range exprs {
		evaluated := Eval(arg, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Enviornment {
	env := object.NewEnclosedEnviornment(fn.Env)

	for i, argAst := range fn.Args {
		env.Set(argAst.Value, args[i])
	}

	return env
}

func unwrapReturn(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if ok {
		extendedEnv := extendFunctionEnv(function, args)
		evaluated := Eval(function.Body, extendedEnv)
		return unwrapReturn(evaluated)
	}

	builtin, ok := fn.(*object.Builtin)
	if ok {
		return builtin.Fn(args...)
	}

	return newErrorf("not a function: %s", fn.Type())
}

func Eval(node ast.Node, env *object.Enviornment) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: n.Value,
		}

	case *ast.StringLiteral:
		return &object.String{
			Value: n.Value,
		}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(n.Value)

	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(n.Operator, right)

	case *ast.InfixExpression:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(n.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatements(n.Statements, env)

	case *ast.IfExpression:
		return evalIfExpression(n, env)

	case *ast.ReturnStatement:
		v := Eval(n.ReturnValue, env)
		if isError(v) {
			return v
		}
		return &object.ReturnValue{Value: v}

	case *ast.LetStatement:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}

		env.Set(n.Name.TokenLiteral(), val)

	case *ast.Identifier:
		return evalIdentifier(n, env)

	case *ast.FunctionLiteral:
		args := n.Parameters
		body := n.Body

		// not sure about the evn
		return &object.Function{Args: args, Body: body, Env: env}

	case *ast.CallExpression:
		fn := Eval(n.Function, env)
		if isError(fn) {
			return fn
		}

		args := evalExpressions(n.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(fn, args)
	}

	return nil
}
