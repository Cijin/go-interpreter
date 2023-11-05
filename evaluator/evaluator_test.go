package evaluator

import (
	"testing"

	"github.com/cijin/go-interpreter/lexer"
	"github.com/cijin/go-interpreter/object"
	"github.com/cijin/go-interpreter/parser"
)

func testEval(in string) object.Object {
	l := lexer.New(in)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, value int64) {
	o, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("expected evaluated to be Integer, got=%T (+%v)", obj, obj)
	}

	if o.Value != value {
		t.Errorf("expected value to be %d, got=%d", value, o.Value)
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not null, got=%T (%+v)", obj, obj)
		return false
	}

	return true
}

func TestIntegerExpression(t *testing.T) {
	tests := []struct {
		in       string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, test := range tests {
		evaluated := testEval(test.in)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, value bool) {
	o, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("expected evaluated to be Integer, got=%T (+%v)", obj, obj)
	}

	if o.Value != value {
		t.Errorf("expected value to be %t, got=%t", value, o.Value)
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		in       string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, test := range tests {
		evaluated := testEval(test.in)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func TestBangPrefixExpressions(t *testing.T) {
	tests := []struct {
		in       string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.in)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfExpressions(t *testing.T) {
	tests := []struct {
		in       string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.in)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		in       string
		expected int
	}{
		{"return 10", 10},
		{"return 10; 9", 10},
		{"9;return 10; 9", 10},
		{"return 2 * 5; 9", 10},
		{"9;return 2 * 5; 9", 10},
		{`
			if (10 > 1) {
				if (10 > 1) {
					return 10
				}

				return 2
			}
		`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.in)
		testIntegerObject(t, evaluated, int64(tt.expected))
	}
}
