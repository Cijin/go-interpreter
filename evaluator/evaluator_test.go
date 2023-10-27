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
		in  string
		out int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, test := range tests {
		evaluated := testEval(test.in)
		testIntegerObject(t, evaluated, test.out)
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
		in  string
		out bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, test := range tests {
		evaluated := testEval(test.in)
		testBooleanObject(t, evaluated, test.out)
	}
}
