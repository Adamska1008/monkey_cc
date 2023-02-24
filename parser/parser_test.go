package parser

import (
	"fmt"
	"log"
	"monkey_cc/ast"
	"monkey_cc/lexer"
	"testing"
)

func assertNoError(t *testing.T, p *Parser) {
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

// 判断表达式是integer literal 类型，且值为value
func assertIntValue(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("is not *ast.IntegerLiteral")
		return false
	}
	if integer.Value != value {
		t.Errorf("value expect %d, found %d", value, integer.Value)
		return false
	}
	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("literal expect %d, found %s", value, integer.TokenLiteral())
		return false
	}
	return true
}

func assertBoolean(t *testing.T, b ast.Expression, value bool) bool {
	boolean, ok := b.(*ast.Boolean)
	if !ok {
		t.Errorf("is not *ast.Boolean")
		return false
	}
	if boolean.Value != value {
		t.Errorf("value expect %v, found %v", boolean.Value, value)
		return false
	}
	return true
}

func assertString(t *testing.T, s ast.Expression, value string) bool {
	str, ok := s.(*ast.StringLiteral)
	if !ok {
		t.Errorf("s is not *ast.StringLiteral")
		return false
	}
	if str.Value != value {
		t.Errorf("str.Value expect: %s, found: %s\n", value, str.Value)
		return false
	}
	return true
}

// 判断表达式是identifier 类型，且值为value
func assertIdentifier(t *testing.T, il ast.Expression, value string) bool {
	ident, ok := il.(*ast.Identifier)
	if !ok {
		t.Errorf("is not *ast.IntegerLiteral")
		return false
	}
	if ident.Value != value {
		t.Errorf("value expect %s, found %s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != fmt.Sprintf("%s", value) {
		t.Errorf("literal expect %s, found %s", value, ident.TokenLiteral())
		return false
	}
	return true
}

// 判断表达式为expect
func assertLiteralExp(t *testing.T, exp ast.Expression, expect interface{}) bool {
	switch v := expect.(type) {
	case int:
		return assertIntValue(t, exp, int64(v))
	case int64:
		return assertIntValue(t, exp, v)
	case bool:
		return assertBoolean(t, exp, v)
	case string:
		{
			switch exp.(type) {
			case *ast.StringLiteral:
				return assertString(t, exp, v)
			case *ast.Identifier:
				return assertIdentifier(t, exp, v)
			}
		}
	}
	t.Errorf("type of exp not handled: %T", exp)
	return false
}

func assertInfixExp(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)

	if !ok {
		t.Errorf("exp is not ast.InfixExpression")
		return false
	}

	if !assertLiteralExp(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator expected: %s, found %s", operator, opExp.Operator)
		return false
	}

	if !assertLiteralExp(t, opExp.Right, right) {
		return false
	}

	return true
}

func TestLetStmt(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 114514;`

	expect := []struct {
		Name  string
		Value int
	}{
		{"x", 5},
		{"y", 10},
		{"foobar", 114514},
	}

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assertNoError(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program does not contain 3 statements: got %d", len(program.Statements))
	}

	for i, st := range program.Statements {
		if st.TokenLiteral() != "let" {
			log.Fatalf("statement %d: literal is not 'let'", i)
		}
		st, ok := st.(*ast.LetStatement)
		if !ok {
			log.Fatalf("statement %d: is not let statement", i)
		}
		if st.Name.Value != expect[i].Name {
			log.Fatalf("statement %d: name of identifier is wrong: expected %s, found %s",
				i, expect[i].Name, st.Name.Value)
		}
		if st.Name.TokenLiteral() != expect[i].Name {
			log.Fatalf("statement %d: token literal of identifier is wrong: expected %s, found %s",
				i, expect[i].Name, st.Name.Value)
		}
		if !assertIntValue(t, st.Value, int64(expect[i].Value)) {
			return
		}
	}
}

func TestReturnStmt(t *testing.T) {
	input := `
	return 5;
	return 10;
	return add(15);`

	expectInt := []int{5, 10}

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assertNoError(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("expected %d statements: got %d", 3, len(program.Statements))
	}

	for i, stmt := range program.Statements {
		rs, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("statement %d: is not return statement", i)
		}
		if rs.TokenLiteral() != "return" {
			t.Fatalf("statement %d: token literal is wrong: expected return, found %s", i, rs.TokenLiteral())
		}
		if i < 2 {
			if !assertIntValue(t, rs.ReturnValue, int64(expectInt[i])) {
				return
			}
		} else {
			ce, ok := rs.ReturnValue.(*ast.CallExpression)
			if !ok {
				t.Fatalf("rs.ReturnValue is not *ast.CallExpression")
			}
			if !assertIdentifier(t, ce.Function, "add") {
				return
			}
			if !assertLiteralExp(t, ce.Arguments[0], 15) {
				return
			}
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assertNoError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected %d statements: got %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("is not expression statement")
	}

	assertLiteralExp(t, stmt.Exp, "foobar")
}

func TestIntExpression(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assertNoError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected %d statements: got %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("is not expression statement")
	}

	assertLiteralExp(t, stmt.Exp, 5)
}

func TestBoolean(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		assertNoError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected %d statements: got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("is not expression statement")
		}

		assertLiteralExp(t, stmt.Exp, tt.expect)
	}
}

func TestStringLiteral(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{`"Hello";`, "Hello"},
		{"World;", "World"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		assertNoError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected %d statements: got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("is not expression statement")
		}

		assertLiteralExp(t, stmt.Exp, tt.expect)
	}
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		assertNoError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected %d statements: got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("is not expression statement")
		}

		exp, ok := stmt.Exp.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("is not prefix expression")
		}

		if exp.Operator != tt.operator {
			t.Fatalf("expression operator is wrong: expect: %s, found: %s", tt.operator, exp.Operator)
		}

		if !assertLiteralExp(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"4 & 3", 4, "&", 3},
		{"4 | 3", 4, "|", 3},
		{"4 && 3", 4, "&&", 3},
		{"6 || 7", 6, "||", 7},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		assertNoError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected %d statements: got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("is not expression statement")
		}

		if !assertInfixExp(t, stmt.Exp, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x; }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assertNoError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected %d statements: got %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not *ast.ExpressionStatement")
	}

	exp, ok := stmt.Exp.(*ast.IfExpression)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression")
	}

	if !assertInfixExp(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("expected %d statements: got %d", 1, len(program.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Consequence.Statement[0] is not *ast.ExpressionStatement")
	}

	if !assertLiteralExp(t, consequence.Exp, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative is not nil")
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assertNoError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected %d statements: got %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statement[0] is not *ast.ExpressionStatement")
	}

	exp, ok := stmt.Exp.(*ast.IfExpression)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression")
	}

	if !assertInfixExp(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("expected %d statements: got %d", 1, len(program.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Consequence.Statement[0] is not *ast.ExpressionStatement")
	}

	if !assertLiteralExp(t, consequence.Exp, "x") {
		return
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Alternative.Statement[0] is not *ast.ExpressionStatement")
	}

	if !assertLiteralExp(t, alternative.Exp, "y") {
		return
	}
}

func TestFunctionLiteral(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assertNoError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected %d program.Statements: got %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Consequence.Statement[0] is not *ast.ExpressionStatement")
	}

	function, ok := stmt.Exp.(*ast.FunctionLiteral)

	if !ok {
		t.Fatalf("stmt.Exp is not *ast.ExpressionStatement")
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("expect %d parameters, found %d", 2, len(function.Parameters))
	}

	assertLiteralExp(t, function.Parameters[0], "x")
	assertLiteralExp(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("expected %d function.Body.Statements: got %d", 1, len(program.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("function.Body.Statement[0] is not *ast.BlockStatement")
	}

	assertInfixExp(t, bodyStmt.Exp, "x", "+", "y")
}

func TestCallExpression(t *testing.T) {
	input := `add(1, 2 * 3, 4 * 5);`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assertNoError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expect %d statements, found %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("function.Body.Statement[0] is not *ast.ExpressionStatement")
	}

	exp, ok := stmt.Exp.(*ast.CallExpression)

	if !ok {
		t.Fatalf("stmt.Exp is not *ast.CallExpression")
	}

	if !assertIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("expect %d exp.Arguments, found %d", 3, len(exp.Arguments))
	}

	assertLiteralExp(t, exp.Arguments[0], 1)
	assertInfixExp(t, exp.Arguments[1], 2, "*", 3)
	assertInfixExp(t, exp.Arguments[2], 4, "*", 5)
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"-a * b;", "((-a) * b);"},
		{"!-a;", "(!(-a));"},
		{"a + b + c;", "((a + b) + c);"},
		{"a + b / c;", "(a + (b / c));"},
		{"5 > 4 == 3 < 4;", "((5 > 4) == (3 < 4));"},
		{"a + b * c + d / e - f;", "(((a + (b * c)) + (d / e)) - f);"},
		{"3 > 5 == false;", "((3 > 5) == false);"},
		{"1 + (2 + 3) + 4;", "((1 + (2 + 3)) + 4);"},
		{"2 / (5 + 5);", "(2 / (5 + 5));"},
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		assertNoError(t, p)
		if program.String() != tt.expect {
			t.Fatalf("case %d:\nexpect: %s\nfound: %s\n", i, tt.expect, program.String())
		}
	}
}

func TestExpectTokenError(t *testing.T) {
	input := `
	let x 5;
	let = 10;
	let 838 383;`
	expect := []string{
		"expected next token to be =, found INT",
		"expected next token to be IDENT, found =",
		"expected next token to be IDENT, found INT",
	}
	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()
	errors := p.Errors()

	if len(errors) == 3 {
		return
	}

	t.Errorf("should detected 3 errors: found %d", len(errors))

	for i, msg := range errors {
		if msg != expect[i] {
			t.Errorf("expected %d error to be: %s\nfound: %s\n", i, expect[i], msg)
		}
	}
	t.FailNow()
}
