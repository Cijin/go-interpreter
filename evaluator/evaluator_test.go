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
	env := object.NewEnviornment()

	return Eval(program, env)
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

func TestStringLiteral(t *testing.T) {
	input := `"hello world"`
	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Errorf("expected *object.String, got=%T", evaluated)
	}

	if str.Value != "hello world" {
		t.Errorf("expected value to be 'hello world', got=%s", str.Value)
	}
}

func TestStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello" + "world";`, "helloworld"},
		{`"hello" + " "  + "world";`, "hello world"},
	}

	for i, tc := range tests {
		evaluated := testEval(tc.input)

		str, ok := evaluated.(*object.String)
		if !ok {
			t.Errorf("[tc %d]:expected *object.String, got=%T", i, evaluated)
		}

		if str.Value != tc.expected {
			t.Errorf("[tc %d]:expected %s, got=%s", i, tc.expected, tc.input)
		}
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
		{`"hello" == "world"`, false},
		{`"hello" == "hello"`, true},
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

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		in       string
		expected string
	}{
		{"true + 5", "type mismatch: BOOLEAN + INTEGER"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "operator '-' not defined on BOOLEAN"},
		{"true + false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}

				return 1;
			}
			`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{"foobar", "identifier is undefined: foobar"},
		{`"hello" - "world"`, "operartor - not supported on type string"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.in)

		err, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("expected object.Error to be returned, got=%T", evaluated)
		}

		if err.Message != tt.expected {
			t.Errorf("expected error to be %s, got=%s", err.Message, tt.expected)
		}
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		in       string
		expected int64
	}{
		{"let x = 5;x;", 5},
		{"let x = 5 * 5;x;", 25},
		{"let x = 5, let y = x;y;", 5},
		{"let x = 5, let y = x; let z = x + y + 5;z;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.in), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 5 }"
	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Errorf("expected evaluated to be of *object.Function type, got=%T", evaluated)
	}

	if len(fn.Args) != 1 {
		t.Errorf("expected 1 argument, got=%d", len(fn.Args))
	}

	if fn.Args[0].TokenLiteral() != "x" {
		t.Errorf("expected arg to be 'x', got=%s", fn.Args[0].TokenLiteral())
	}

	expectedBody := "(x + 5)"
	if fn.Body.String() != expectedBody {
		t.Errorf("expected body to be %s, got=%s", expectedBody, fn.Body.String())
	}
}

func TestFunctionCall(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let x = fn(y) { y }; x(5);", 5},
		{"let x = fn(y) { return y }; x(5);", 5},
		{"let double = fn(x) { return x * 2 }; double(5);", 10},
		{"let add = fn(x, y) { return x + y }; add(5, 5);", 10},
		{"let add = fn(x, y) { return x + y }; add(5, add(5, 5));", 15},
		{"fn(x, y) { return x + y }(5, 5);", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
	let x = 10;
	let y = 10;

	let addTwo = fn(x) {
		return fn(y) {
			return x + y;
		};
	};

	let add = addTwo(2);
	add(4);
	`

	var expected int64 = 6

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, expected)
}
