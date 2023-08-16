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
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1;2;",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 - 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 * 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "2 / 1",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
		},
	}
	runTests(t, tests)
}

func TestBooleanExp(t *testing.T) {
	tests := []compilerTest{
		{
			input:             "true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "false",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 < 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpLess),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 > 2",
			expectedConstants: []interface{}{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpLess),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 != 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1 == 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
	}
	runTests(t, tests)
}

func TestPrefixExp(t *testing.T) {
	tests := []compilerTest{
		{
			input:             "-1",
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "!true",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
	}
	runTests(t, tests)
}

func TestIfExp(t *testing.T) {
	tests := []compilerTest{
		{
			input:             "if(true){10}3333;",
			expectedConstants: []interface{}{10, 3333},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpTrue),
				// 0001
				code.Make(code.OpJumpNotTruthy, 10),
				// 0004
				code.Make(code.OpConstant, 0),
				// 0007
				code.Make(code.OpJump, 11),
				// 0010
				code.Make(code.OpNull),
				// 0011
				code.Make(code.OpPop),
				// 0012
				code.Make(code.OpConstant, 1),
				// 0015
				code.Make(code.OpPop),
			},
		},
		{
			input:             "if(true){10}else{20};3333;",
			expectedConstants: []interface{}{10, 20, 3333},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpTrue),
				// 0001
				code.Make(code.OpJumpNotTruthy, 10),
				// 0004
				code.Make(code.OpConstant, 0),
				// 0007
				code.Make(code.OpJump, 13),
				// 00010
				code.Make(code.OpConstant, 1),
				// 0013
				code.Make(code.OpPop),
				// 0014
				code.Make(code.OpConstant, 2),
				// 0017
				code.Make(code.OpPop),
			},
		},
	}
	runTests(t, tests)
}

func TestGlobalLetStmt(t *testing.T) {
	tests := []compilerTest{
		{
			input: `let one = 1;
					let two = 2;`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
			},
		},
		{
			input: `let one = 1;
					one;`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
	}
	runTests(t, tests)
}

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Scope: GlobalScope, Index: 1},
	}
	global := NewSymbolTable()
	a := global.Define("a")
	if a != expected["a"] {
		t.Errorf(NOT_EXPECTED, "a", expected["a"], a)
	}
	b := global.Define("b")
	if b != expected["b"] {
		t.Errorf(NOT_EXPECTED, "b", expected["b"], b)
	}
}

func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")
	expect := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "b", Scope: GlobalScope, Index: 1},
	}
	for _, sym := range expect {
		result, ok := global.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
		}
		if result != sym {
			t.Errorf(NOT_EXPECTED, sym.Name, sym, result)
		}
	}
}
