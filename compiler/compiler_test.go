package compiler

import (
	"fmt"
	"monkey_cc/ast"
	"monkey_cc/code"
	"monkey_cc/lexer"
	"monkey_cc/object"
	"monkey_cc/parser"
	"testing"
)

const (
	COMPILER_ERROR     = "compiler error: %v"
	INSTRUCTIONS_ERROR = "testInstructions failed: %v"
	CONSTANTS_ERROR    = "testConstants failed: %v"
	WRONG_LENGTH       = "%v has wrong length: expected=%v, found=%v"
	NOT_EXPECTED       = "expected %v to be %v, found %v"
)

// 测试生成的常量池和指令流
type compilerTest struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

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

func testInstructions(val code.Instructions, expected []code.Instructions) error {
	concatted := code.Instructions{}
	for _, ins := range expected {
		concatted = append(concatted, ins...)
	}
	if len(val) != len(concatted) {
		return fmt.Errorf(WRONG_LENGTH, "len(val)", len(concatted), len(val))
	}
	for i := range concatted {
		if val[i] != concatted[i] {
			return fmt.Errorf(NOT_EXPECTED, "val[i]", concatted[i], val[i])
		}
	}
	return nil
}

func testConstants(val []object.Object, expected []interface{}, t *testing.T) error {
	if len(val) != len(expected) {
		return fmt.Errorf(WRONG_LENGTH, "val", len(expected), len(val))
	}
	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(val[i], int64(constant))
			if err != nil {
				return fmt.Errorf(CONSTANTS_ERROR, err)
			}
		}
	}
	return nil
}

func runTests(t *testing.T, tests []compilerTest) {
	// 辅助函数，用于标记函数使其在输出日志时不出现
	t.Helper()
	for _, tt := range tests {
		program := parse(tt.input)
		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf(COMPILER_ERROR, err)
		}
		// 使用bytecode进行测试，分别测试指令字节和常量池
		bytecode := compiler.Bytecode()
		err = testInstructions(bytecode.Instructions, tt.expectedInstructions)
		if err != nil {
			t.Fatalf(INSTRUCTIONS_ERROR, err)
		}
		err = testConstants(bytecode.Constants, tt.expectedConstants, t)
		if err != nil {
			t.Fatalf(CONSTANTS_ERROR, err)
		}
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTest{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
			},
		},
	}
	runTests(t, tests)
}
