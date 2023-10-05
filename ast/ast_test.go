package ast

import (
	"testing"

	"github.com/cijin/go-interpreter/token"
)

func TestString(t *testing.T) {
	expected := "let x = y;"
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{
					Type:    token.LET,
					Literal: "let",
				},
				Name: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "x",
					},
					Value: "y",
				},
				Value: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "y",
					},
				},
			},
		},
	}

	if program.String() != expected {
		t.Errorf("Expected %s, got %s\n", expected, program.String())
	}
}
