package vm

import (
	"fmt"
	"monkey_cc/ast"
	"monkey_cc/compiler"
	"monkey_cc/lexer"
	"monkey_cc/object"
	"monkey_cc/parser"
	"testing"
)

const (
	COMPILER_ERROR     = "compiler error: %v"
	VM_ERROR           = "vm error: %v"
	INSTRUCTIONS_ERROR = "testInstructions failed: %v"
	CONSTANTS_ERROR    = "testConstants failed: %v"
	WRONG_LENGTH       = "%v has wrong length: expected=%v, found=%v"
	NOT_EXPECTED       = "expected %v to be %v, found %v"
)

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(val object.Object, expected int64) error {
	result, ok := val.(*object.Integer)
	if !ok {
		return fmt.Errorf(NOT_EXPECTED, "val.(type)", "*object.Integer", val.Type())
	}
	if result.Value != expected {
		return fmt.Errorf(NOT_EXPECTED, "result.Value", expected, result.Value)
	}
	return nil
}

func testObject(t *testing.T, expected interface{}, val object.Object) {
	t.Helper()
	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(val, int64(expected))
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	}
}

type vmTest struct {
	input    string
	expected interface{} // VM的运行结果为栈顶数据
}

func runTests(t *testing.T, tests []vmTest) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf(COMPILER_ERROR, err)
		}
		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf(VM_ERROR, err)
		}
		testObject(t, tt.expected, vm.StackTop())
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTest{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 2},
	}
	runTests(t, tests)
}
