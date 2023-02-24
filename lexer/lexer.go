package lexer

import "monkey_cc/token"

type Lexer struct {
	input string
	// pos 指向当前读取的字符 pos == -1 表示尚未开始读取
	pos int
}

// 读取下一个字符，但是指针不向前移动
// 当下一个字符不存在时，返回0
func (l *Lexer) peekChar() byte {
	if l.pos == len(l.input)-1 {
		return 0
	} else {
		return l.input[l.pos+1]
	}
}

// 读取下一个字符，同时指针前移
// 当下一个字符不存在时，返回0，且指针不会移动
func (l *Lexer) nextChar() byte {
	next := l.peekChar()
	if next == 0 {
		return 0
	} else {
		l.pos++
		return next
	}
}

// 消耗输入中的空白字符，包括空格，制表符，换行符
// 使得执行完后，l.nextChar()是非空白字符
func (l *Lexer) consumeSpaces() {
	for {
		ch := l.peekChar()
		if ch != '\n' && ch != '\r' && ch != '\t' && ch != ' ' {
			break
		}
		l.nextChar()
	}
}

// 读取标识符
func (l *Lexer) readIdentifier() string {
	begin := l.pos
	for {
		ch := l.peekChar()
		if !isLetter(ch) {
			break
		}
		l.nextChar()
	}
	return l.input[begin : l.pos+1]
}

// 读取数字
func (l *Lexer) readInteger() string {
	begin := l.pos
	for {
		ch := l.peekChar()
		if !isDigit(ch) {
			break
		}
		l.nextChar()
	}
	return l.input[begin : l.pos+1]
}

func (l *Lexer) readString() string {
	// 跳过‘”’，需要将pos+1
	begin := l.pos + 1
	for l.nextChar() != '"' {
	}
	return l.input[begin:l.pos]
}

// NextToken 读取下一个词法单元，同时指针前移
func (l *Lexer) NextToken() *token.Token {
	l.consumeSpaces()
	ch := l.nextChar()
	switch ch {
	case '=':
		if l.peekChar() == '=' {
			l.nextChar()
			return token.New(token.EQ, "==")
		} else {
			return token.New(token.ASSIGN, "=")
		}
	case '+':
		return token.New(token.PLUS, "+")
	case '-':
		return token.New(token.MINUS, "-")
	case '!':
		if l.peekChar() == '=' {
			l.nextChar()
			return token.New(token.NOT_EQ, "!=")
		} else {
			return token.New(token.BANG, "!")
		}
	case '*':
		return token.New(token.ASTERISK, "*")
	case '/':
		return token.New(token.SLASH, "/")
	case '<':
		return token.New(token.LT, "<")
	case '>':
		return token.New(token.GT, ">")
	case ';':
		return token.New(token.SEMICOLON, ";")
	case ',':
		return token.New(token.COMMA, ",")
	case '(':
		return token.New(token.LPAREN, "(")
	case ')':
		return token.New(token.RPAREN, ")")
	case '{':
		return token.New(token.LBRACE, "{")
	case '}':
		return token.New(token.RBRACE, "}")
	case '&':
		if l.peekChar() == '&' {
			l.nextChar()
			return token.New(token.AND, "&&")
		} else {
			return token.New(token.BIT_AND, "&")
		}
	case '|':
		if l.peekChar() == '|' {
			l.nextChar()
			return token.New(token.OR, "||")
		} else {
			return token.New(token.BIT_OR, "|")
		}
	case '"':
		literal := l.readString()
		return token.New(token.STRING, literal)
	case 0:
		return token.New(token.EOF, "")
	default:
		if isLetter(ch) {
			literal := l.readIdentifier()
			return token.New(token.LookupIdent(literal), literal)
		} else if isDigit(ch) {
			literal := l.readInteger()
			return token.New(token.INT, literal)
		} else {
			return token.New(token.ILLEGAL, "")
		}
	}
}

func New(s string) *Lexer {
	return &Lexer{
		input: s,
		pos:   -1,
	}
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
