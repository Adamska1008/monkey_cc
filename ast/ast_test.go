package ast

import (
	"monkey_cc/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: *token.New(token.LET, "let"),
				Name: &Identifier{
					Token: *token.New(token.IDENT, "myVar"),
					Value: "myVar",
				},
				Value: &Identifier{
					Token: *token.New(token.IDENT, "anotherVar"),
					Value: "anotherVar",
				},
			},
		},
	}
	expect := "let myVar = anotherVar;"
	if program.String() != expect {
		t.Errorf("program.String() is wrong:\nexpect: %s\nfound:%s\n", expect, program.String())
	}
}
