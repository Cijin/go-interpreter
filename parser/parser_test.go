package parser

import (
	"testing"

	"github.com/cijin/go-interpreter/ast"
	"github.com/cijin/go-interpreter/lexer"
)

// Helper method
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("token literal is not let, got %s", s.TokenLiteral())
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("expected *ast.LetStatment, but got %T", s)
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

func TestLetStatements(t *testing.T) {
	input := `
	let x  5;
	let y  10;
	let foobar 838383;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatal("ParsePorgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("expected %d statements, got %d", 3, len(program.Statements))
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
		t.Fatalf("expected %d statements, got %d", 3, len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("expected *ast.ReturnStatement, but got %T", returnStmt)
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("expected token literal 'return', but got %T", returnStmt.TokenLiteral())
		}
	}
}
