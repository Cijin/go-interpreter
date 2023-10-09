package parser

import (
	"fmt"
	"testing"

	"github.com/cijin/go-interpreter/ast"
	"github.com/cijin/go-interpreter/lexer"
)

// Helper method
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("token literal is not let, got=%s", s.TokenLiteral())
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("expected *ast.LetStatment, but got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letstmt.Name.Value not eq to %s, got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Value.TokenLiteral() != name {
		t.Errorf("letstmt.Value.TokenLiteral() not eq to %s, got=%s", name, letStmt.Value.TokenLiteral())
		return false
	}

	return true
}

// Helper method
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %s", msg)
	}
	t.FailNow()
}

// Helper method
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expected *ast.IntegerLiteral, but got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("expected value to be %d, got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("expected token literal to be %s, got=%s", fmt.Sprintf("%d", value), integ.TokenLiteral())
		return false
	}

	return true
}

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatal("ParsePorgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("expected %d statements, got=%d", 3, len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]

		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 999893;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatal("ParseProgarm() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("expected %d statements, got=%d", 3, len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("expected *ast.ReturnStatement, but got=%T", returnStmt)
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("expected token literal 'return', but got=%T", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("expected %d statements but got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpresssionStatement)
	if !ok {
		t.Errorf("expected ast.ExpressionStatement but got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Errorf("expected ast.Identifier but got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("expected value to be %q got=%q", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("expected value to be %q got=%q", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("expected %d statements got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpresssionStatement)
	if !ok {
		t.Errorf("expected ast.ExpressionStatement got=%T", program.Statements[0])
	}

	if !testIntegerLiteral(t, stmt.Expression, 5) {
		return
	}
}
