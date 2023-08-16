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

func testBooleanObject(val object.Object, expected bool) error {
	result, ok := val.(*object.Boolean)
	if !ok {
		return fmt.Errorf(NOT_EXPECTED, "val.(type)", "*object.Boolean", val.Type())
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
	case bool:
		err := testBooleanObject(val, bool(expected))
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	case *object.Null:
		if val != Null {
			t.Errorf("object is not Null: %T (%+v)", val, val)
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
		testObject(t, tt.expected, vm.LastPopped())
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTest{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"-1", -1},
		{"(1 + 3 * 5) * (-1)", -16},
	}
	runTests(t, tests)
}

func TestBooleanExp(t *testing.T) {
	tests := []vmTest{
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
		{"true == false", false},
		{"(1 < 2) == true", true},
		{"!true", false},
		{"!false", true},
	}
	runTests(t, tests)
}

func TestIfExp(t *testing.T) {
	tests := []vmTest{
		{"if (true) {10}", 10},
		{"if (true) {10} else {20}", 10},
		{"if (1 > 2) {10} else {20}", 20},
		{"if (1 > 2) {10}", Null},
	}
	runTests(t, tests)
}

func TestGlobalLetStmt(t *testing.T) {
	tests := []vmTest{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one + two", 3},
	}
	runTests(t, tests)
}
