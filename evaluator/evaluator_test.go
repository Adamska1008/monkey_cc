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
	env := object.NewEnvironment()
	return Eval(program, env)
}

func assertInteger(t *testing.T, obj object.Object, expect int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("obj is not *object.Integer")
		return false
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
		return false
	}
	if result.Value != expect {
		t.Fatalf("obj: expect %t, found %t", expect, result.Value)
		return false
	}
	return true
}

func assertString(t *testing.T, obj object.Object, expect string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Fatalf("obj is not *object.Boolean")
		return false
	}
	if result.Value != expect {
		t.Fatalf("obj: expect %s, found %s", expect, result.Value)
		return false
	}
	return true
}

func assertNull(t *testing.T, obj object.Object) bool {
	_, ok := obj.(*object.Null)
	if !ok {
		t.Fatalf("obj is not *object.Null")
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
		{"5 + 5 + 5 - 10", 5},
		{"2 * 2 * 2", 8},
		{"5 * 2 + 10", 20},
		{"2 * (5 + 10)", 30},
		{"5 + 10 / 2", 10},
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
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"true && false", false},
		{"true || false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !assertBoolean(t, evaluated, tt.expect) {
			return
		}
	}
}

func TestEvalString(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{`"Hello, " + "world!"`, "Hello, world!"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !assertString(t, evaluated, tt.expect) {
			return
		}
	}
}

func TestEvalArray(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3];"
	evaluated := testEval(input)
	result, _ := evaluated.(*object.Array)

	if len(result.Elements) != 3 {
		t.Fatalf("expect len(result.Elements) to be %d, found %d", 3, len(result.Elements))
	}

	assertInteger(t, result.Elements[0], 1)
	assertInteger(t, result.Elements[1], 4)
	assertInteger(t, result.Elements[2], 6)
}

func TestEvalHash(t *testing.T) {
	input := `{"one": 1 + 9, 4: 4}`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("evaluated is not *object.Hash")
	}
	for _, v := range result.Pairs {
		switch v.Key.Type() {
		case object.STRING_OBJ:
			assertString(t, v.Key, "one")
			assertInteger(t, v.Value, 10)
		case object.INTEGER_OBJ:
			assertInteger(t, v.Key, 4)
			assertInteger(t, v.Value, 4)
		}
	}
}

func TestIndexExpression(t *testing.T) {
	input := `let arr = [1, 2, 3, 4]; arr[2];`
	evaluated := testEval(input)
	result, _ := evaluated.(*object.Integer)

	assertInteger(t, result, 3)
}

func TestHashIndexExpression(t *testing.T) {
	input := `let hash = {"hello": "world"}; hash["hello"]`
	evaluated := testEval(input)
	result, _ := evaluated.(*object.String)
	assertString(t, result, "world")
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

func TestIfElseExp(t *testing.T) {
	tests := []struct {
		input  string
		expect interface{}
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
		evaluated := testEval(tt.input)
		integer, ok := tt.expect.(int)
		if ok {
			assertInteger(t, evaluated, int64(integer))
		} else {
			assertNull(t, evaluated)
		}
	}
}

func TestReturnStmts(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`if (10 > 1) { if (10 > 1) { return 10; } } return 1;`, 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		assertInteger(t, evaluated, tt.expect)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"false + true",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
	}

	for i, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Fatalf("evaluated is not *object.Error")
		}
		if errObj.Message != tt.expectedMessage {
			t.Fatalf("tests %d:\nexpect error message:\n%s\nfound:\n%s\n", i, tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStmt(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		assertInteger(t, testEval(tt.input), tt.expect)
	}
}

func TestFunction(t *testing.T) {
	input := `fn(x) { x + 2; }`
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("evaluated is not *object.Funtion")
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("expect len(fn.Parameters): %d, found: %d\n", 1, len(fn.Parameters))
	}
	if fn.Parameters[0].Value != "x" {
		t.Fatalf("expect fn.Parameters[0].Value: %s, found: %s", "x", fn.Parameters[0].Value)
	}
	if fn.Body.String() != "{(x + 2);}" {
		t.Fatalf("expect fn.Body.String(): %s, found: %s", "{(x + 2);}", fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
	}
	for _, tt := range tests {
		assertInteger(t, testEval(tt.input), tt.expect)
	}
}

func TestBuiltInFunction(t *testing.T) {
	tests := []struct {
		input  string
		expect int64
	}{
		{`len("")`, 0},
		{`len("hello world")`, 11},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		assertInteger(t, evaluated, tt.expect)
	}
}
