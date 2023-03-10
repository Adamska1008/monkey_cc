package parser

import (
	"fmt"
	"monkey_cc/ast"
	"monkey_cc/lexer"
	"monkey_cc/token"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// 优先级
const (
	_ int = iota
	LOWEST
	OR
	AND
	BIT_OR
	BIT_AND
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

var precedences = map[token.TokenType]int{
	token.AND:      AND,
	token.OR:       OR,
	token.BIT_AND:  BIT_AND,
	token.BIT_OR:   BIT_OR,
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type Parser struct {
	l      *lexer.Lexer
	errors []string
	peek   *token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New 初始化Parser
// 注册解析函数
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         []string{},
		peek:           nil,
		prefixParseFns: map[token.TokenType]prefixParseFn{},
		infixParseFns:  map[token.TokenType]infixParseFn{},
	}
	// 初始化的parser的peekToken()为nil，故先消耗这个空token
	p.nextToken()

	p.registerPrefix(token.IDENT, p.ParseIdent)
	p.registerPrefix(token.INT, p.ParseInt)
	p.registerPrefix(token.TRUE, p.ParseBoolean)
	p.registerPrefix(token.FALSE, p.ParseBoolean)
	p.registerPrefix(token.STRING, p.ParseString)
	p.registerPrefix(token.MINUS, p.ParsePrefixExpression)
	p.registerPrefix(token.BANG, p.ParsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.ParseGroupedExp)
	p.registerPrefix(token.IF, p.ParseIfExp)
	p.registerPrefix(token.FUNCTION, p.ParseFnLiteral)
	p.registerPrefix(token.LBRACKET, p.ParseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.ParseHashLiteral)

	p.registerInfix(token.PLUS, p.ParseInfixExpression)
	p.registerInfix(token.MINUS, p.ParseInfixExpression)
	p.registerInfix(token.SLASH, p.ParseInfixExpression)
	p.registerInfix(token.ASTERISK, p.ParseInfixExpression)
	p.registerInfix(token.EQ, p.ParseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.ParseInfixExpression)
	p.registerInfix(token.LT, p.ParseInfixExpression)
	p.registerInfix(token.GT, p.ParseInfixExpression)
	p.registerInfix(token.AND, p.ParseInfixExpression)
	p.registerInfix(token.OR, p.ParseInfixExpression)
	p.registerInfix(token.BIT_AND, p.ParseInfixExpression)
	p.registerInfix(token.BIT_OR, p.ParseInfixExpression)
	p.registerInfix(token.LPAREN, p.ParseCallExp)
	p.registerInfix(token.LBRACKET, p.ParseIndexExpression)

	return p
}

func (p *Parser) peekToken() *token.Token {
	return p.peek
}

func (p *Parser) nextToken() *token.Token {
	cur := p.peekToken()
	p.peek = p.l.NextToken()
	return cur
}

// 获取目前token的优先级。
// 若token不在列表中，返回LOWEST
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken().Type]; ok {
		return p
	}
	return LOWEST
}

// 消耗所有到分号的token，包括这个分号
func (p *Parser) skipToSemicolonOrRBrace() {
	peek := p.peekToken()
	for peek.Type != token.SEMICOLON && peek.Type != token.RBRACE {
		p.nextToken()
		peek = p.peekToken()
	}
	p.nextToken()
}

// 会消耗形如"(...)"的词法单元，返回Ident列表
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier

	p.nextToken()
	if p.peekToken().Type == token.RPAREN {
		p.nextToken()
		return identifiers
	}

	ident := &ast.Identifier{
		Token: *p.peekToken(),
		Value: p.peekToken().Literal,
	}
	identifiers = append(identifiers, ident)
	p.nextToken()

	for p.peekToken().Type == token.COMMA {
		p.nextToken()
		ident := &ast.Identifier{
			Token: *p.peekToken(),
			Value: p.peekToken().Literal,
		}
		identifiers = append(identifiers, ident)
		p.nextToken()
	}

	if !p.expectPeekType(token.RPAREN) {
		return nil
	}
	p.nextToken()
	return identifiers
}

// 相较于前者，该方法返回表达式列表
func (p *Parser) parseCallArguments() []ast.Expression {
	var args []ast.Expression

	p.nextToken()
	if p.peekToken().Type == token.RPAREN {
		p.nextToken()
		return args
	}

	args = append(args, p.ParseExp(LOWEST))
	for p.peekToken().Type == token.COMMA {
		p.nextToken()
		args = append(args, p.ParseExp(LOWEST))
	}

	if !p.expectPeekType(token.RPAREN) {
		return nil
	}
	p.nextToken()
	return args
}

// 相较于前者，该方法解析形如"[...]"的列表，可以考虑将两个方法合并
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	var args []ast.Expression

	p.nextToken()
	if p.peekToken().Type == token.RPAREN {
		p.nextToken()
		return args
	}

	args = append(args, p.ParseExp(LOWEST))
	for p.peekToken().Type == token.COMMA {
		p.nextToken()
		args = append(args, p.ParseExp(LOWEST))
	}

	if !p.expectPeekType(end) {
		return nil
	}
	p.nextToken()
	return args
}

// 当期待token与peekToken不一致时，产生该error
func (p *Parser) expectTokenError(expect token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, found %s", expect, p.peekToken().Type)
	p.errors = append(p.errors, msg)
}

// 当当前token不是int时，产生该error
func (p *Parser) parseIntError(s string) {
	msg := fmt.Sprintf("could not parse %s as integer", s)
	p.errors = append(p.errors, msg)
}

// 当前token不存在prefix方法时，产生该error
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// 判断peekToken的TokenType是否为expect，返回bool值
// 不会消耗peekToken
// 当判断不符合时，自动添加错误
func (p *Parser) expectPeekType(expect token.TokenType) bool {
	if p.peekToken().Type == expect {
		return true
	} else {
		p.expectTokenError(expect)
		return false
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}
	for p.peekToken().Type != token.EOF {
		stmt := p.ParseStmt()
		program.Statements = append(program.Statements, stmt)
	}
	return program
}

func (p *Parser) ParseStmt() ast.Statement {
	// 路由过程不消耗token
	switch p.peekToken().Type {
	case token.LET:
		return p.ParseLetStmt()
	case token.RETURN:
		return p.ParseReturnStmt()
	default:
		return p.ParseExpStmt()
	}
}

func (p *Parser) ParseLetStmt() *ast.LetStatement {
	ls := &ast.LetStatement{Token: *p.nextToken()}

	if !p.expectPeekType(token.IDENT) {
		p.skipToSemicolonOrRBrace()
		return nil
	}

	idt := *p.nextToken()
	ls.Name = &ast.Identifier{
		Token: idt,
		Value: idt.Literal,
	}

	if !p.expectPeekType(token.ASSIGN) {
		p.skipToSemicolonOrRBrace()
		return nil
	}
	p.nextToken()
	// 省略表达式求值部分
	//p.skipToSemicolonOrRBrace()
	ls.Value = p.ParseExp(LOWEST)
	if !p.expectPeekType(token.SEMICOLON) {
		return nil
	}
	p.nextToken()
	return ls
}

func (p *Parser) ParseReturnStmt() *ast.ReturnStatement {
	rs := &ast.ReturnStatement{Token: *p.nextToken()}
	// 省略表达式求值部分
	//p.skipToSemicolonOrRBrace()
	rs.ReturnValue = p.ParseExp(LOWEST)
	if !p.expectPeekType(token.SEMICOLON) {
		return nil
	}
	p.nextToken()
	return rs
}

func (p *Parser) ParseBlockStmt() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token:      *p.peekToken(),
		Statements: []ast.Statement{},
	}
	p.nextToken()

	for p.peekToken().Type != token.RBRACE && p.peekToken().Type != token.EOF {
		stmt := p.ParseStmt()
		block.Statements = append(block.Statements, stmt)
	}
	if p.peekToken().Type == token.RBRACE {
		p.nextToken()
	}
	return block
}

func (p *Parser) ParseExpStmt() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: *p.peekToken()}
	stmt.Exp = p.ParseExp(LOWEST)
	// p.skipToSemicolonOrRBrace()
	// 跳过最后的token，如";"
	// 若为BlockStmt，可能没有";"
	if p.peekToken().Type == token.SEMICOLON {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) ParseExp(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.peekToken().Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.peekToken().Type)
		return nil
	}
	leftExp := prefix()
	for p.peekToken().Type != token.SEMICOLON && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken().Type]
		if infix == nil {
			return leftExp
		}
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) ParsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    *p.peekToken(),
		Operator: p.peekToken().Literal,
	}
	p.nextToken()
	expression.Right = p.ParseExp(PREFIX)
	return expression
}

func (p *Parser) ParseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    *p.peekToken(),
		Left:     left,
		Operator: p.peekToken().Literal,
	}
	precedence := p.peekPrecedence()
	p.nextToken()
	expression.Right = p.ParseExp(precedence)
	return expression
}

func (p *Parser) ParseIdent() ast.Expression {
	return &ast.Identifier{
		Token: *p.peekToken(),
		Value: p.nextToken().Literal,
	}
}

func (p *Parser) ParseInt() ast.Expression {
	lit := &ast.IntegerLiteral{Token: *p.peekToken()}
	value, err := strconv.ParseInt(p.peekToken().Literal, 0, 64)
	if err != nil {
		p.parseIntError(p.peekToken().Literal)
		return nil
	}
	lit.Value = value
	p.nextToken()
	return lit
}

func (p *Parser) ParseBoolean() ast.Expression {
	boolean := &ast.Boolean{
		Token: *p.peekToken(),
		Value: p.nextToken().Type == token.TRUE,
	}
	return boolean
}

func (p *Parser) ParseString() ast.Expression {
	str := &ast.StringLiteral{
		Token: *p.peekToken(),
		Value: p.nextToken().Literal,
	}
	return str
}

func (p *Parser) ParseGroupedExp() ast.Expression {
	p.nextToken()
	exp := p.ParseExp(LOWEST)

	if p.nextToken().Type != token.RPAREN {
		return nil
	}
	return exp
}

func (p *Parser) ParseIfExp() ast.Expression {
	exp := &ast.IfExpression{
		Token: *p.peekToken(),
	}
	p.nextToken()
	if p.peekToken().Type != token.LPAREN {
		return nil
	}
	p.nextToken()
	exp.Condition = p.ParseExp(LOWEST)
	if !p.expectPeekType(token.RPAREN) {
		return nil
	}
	p.nextToken()
	if !p.expectPeekType(token.LBRACE) {
		return nil
	}
	exp.Consequence = p.ParseBlockStmt()
	if p.peekToken().Type == token.ELSE {
		p.nextToken()
		exp.Alternative = p.ParseBlockStmt()
	}
	return exp
}

func (p *Parser) ParseFnLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: *p.peekToken()}
	p.nextToken()
	if !p.expectPeekType(token.LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
	if !p.expectPeekType(token.LBRACE) {
		return nil
	}
	lit.Body = p.ParseBlockStmt()
	return lit
}

// ParseCallExp 调用表达式，左值为函数变量名
func (p *Parser) ParseCallExp(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    *p.peekToken(), // "("
		Function: function,
	}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) ParseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: *p.peekToken()} // "["
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) ParseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{
		Token: *p.peekToken(), // "["
		Left:  left,
	}
	p.nextToken()
	exp.Index = p.ParseExp(LOWEST)
	if !p.expectPeekType(token.RBRACKET) {
		return nil
	}
	p.nextToken()
	return exp
}

func (p *Parser) ParseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: *p.peekToken()} // "{"
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for p.peekToken().Type != token.RBRACE {
		p.nextToken()
		key := p.ParseExp(LOWEST)
		if !p.expectPeekType(token.COLON) {
			return nil
		}
		p.nextToken()
		value := p.ParseExp(LOWEST)
		hash.Pairs[key] = value
		if p.peekToken().Type != token.RBRACE && !p.expectPeekType(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeekType(token.RBRACE) {
		return nil
	}
	p.nextToken()
	return hash
}
