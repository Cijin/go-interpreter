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
