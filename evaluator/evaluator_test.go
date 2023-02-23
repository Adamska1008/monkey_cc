package evaluator

import (
	"monkey_cc/lexer"
	"monkey_cc/object"
	"monkey_cc/parser"
	"testing"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}

func assertInteger(t *testing.T, obj object.Object, expect int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("obj is not *object.Integer")
	}
	if result.Value != expect {
		t.Fatalf("obj: expect %d, found %d", expect, result.Value)
		return false
	}
	return true
}

func assertBoolean(t *testing.T, obj object.Object, expect bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Fatalf("obj is not *object.Boolean")
	}
	if result.Value != expect {
		t.Fatalf("obj: expect %t, found %t", expect, result.Value)
		return false
	}
	return true
}

func TestEvalInteger(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !assertInteger(t, evaluated, tt.expect) {
			return
		}
	}
}

func TestEvalBoolean(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !assertBoolean(t, evaluated, tt.expect) {
			return
		}
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"!false", true},
		{"!true", false},
		{"!5", false},
		{"!!true", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		assertBoolean(t, evaluated, tt.expect)
	}
}
