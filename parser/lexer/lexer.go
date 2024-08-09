package lexer

import (
	"github.com/stevecallear/mexl/ast/token"
)

type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      byte
}

var keywords = map[string]token.Type{
	"true":  token.True,
	"false": token.False,
	"and":   token.And,
	"or":    token.Or,
	"not":   token.Bang,
	"eq":    token.Equal,
	"ne":    token.NotEqual,
	"lt":    token.LessThan,
	"gt":    token.GreaterThan,
	"le":    token.LessThanOrEqual,
	"ge":    token.GreaterThanOrEqual,
	"sw":    token.StartsWith,
	"ew":    token.EndsWith,
	"in":    token.In,
	"null":  token.Null,
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() (t token.Token) {
	newToken := func(t token.Type, chs ...byte) token.Token {
		return token.Token{Type: t, Literal: string(chs)}
	}

	l.skipWhitespace()

	var skipRead bool // todo: it should be possible to remove this

	switch l.ch {
	case 0:
		t.Type = token.EOF
		t.Literal = "EOF"
	case '+':
		t = newToken(token.Plus, l.ch)

	case '-':
		t = newToken(token.Minus, l.ch)

	case '*':
		t = newToken(token.Asterisk, l.ch)

	case '/':
		t = newToken(token.Slash, l.ch)

	case '%':
		t = newToken(token.Percent, l.ch)

	case '.':
		pc := l.peekChar()
		switch {
		case isDigit(pc):
			t = l.readNumber()
			skipRead = true
		default:
			t = newToken(token.Stop, l.ch)
		}

	case ',':
		t = newToken(token.Comma, l.ch)

	case '(':
		t = newToken(token.LParen, l.ch)

	case ')':
		t = newToken(token.RParen, l.ch)

	case '[':
		t = newToken(token.LBracket, l.ch)

	case ']':
		t = newToken(token.RBracket, l.ch)

	case '"':
		t = l.readString()

	case '=':
		ch := l.ch
		l.readChar()

		switch l.ch {
		case '=':
			t = newToken(token.Equal, ch, l.ch)
		default:
			t = newToken(token.Illegal, ch, l.ch)
		}

	case '!':
		switch l.peekChar() {
		case '=':
			ch := l.ch
			l.readChar()
			t = newToken(token.NotEqual, ch, l.ch)
		default:
			t = newToken(token.Bang, l.ch)
		}

	case '<':
		switch l.peekChar() {
		case '=':
			ch := l.ch
			l.readChar()
			t = newToken(token.LessThanOrEqual, ch, l.ch)
		default:
			t = newToken(token.LessThan, l.ch)
		}

	case '>':
		switch l.peekChar() {
		case '=':
			ch := l.ch
			l.readChar()
			t = newToken(token.GreaterThanOrEqual, ch, l.ch)
		default:
			t = newToken(token.GreaterThan, l.ch)
		}

	case '&':
		ch := l.ch
		l.readChar()

		switch l.ch {
		case '&':
			t = newToken(token.And, ch, l.ch)
		default:
			t = newToken(token.Illegal, ch, l.ch)
		}

	case '|':
		ch := l.ch
		l.readChar()

		switch l.ch {
		case '|':
			t = newToken(token.Or, ch, l.ch)
		default:
			t = newToken(token.Illegal, ch, l.ch)
		}

	default:
		switch {
		case isLetter(l.ch):
			t.Literal = l.readWord()
			var ok bool
			t.Type, ok = keywords[t.Literal]
			if !ok {
				t.Type = token.Ident
			}
			skipRead = true
		case isDigit(l.ch):
			t = l.readNumber()
			skipRead = true
		default:
			t = newToken(token.Illegal, l.ch)
		}
	}

	if !skipRead {
		l.readChar()
	}

	return t
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}

	l.pos = l.readPos
	l.readPos += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPos]
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readWord() string {
	p := l.pos

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[p:l.pos]
}

func (l *Lexer) readNumber() token.Token {
	p := l.pos
	t := token.Token{Type: token.Int}

	for isDigit(l.ch) || isDot(l.ch) {
		if isDot(l.ch) {
			t.Type = token.Float
		}
		l.readChar()
	}

	t.Literal = l.input[p:l.pos]
	return t
}

func (l *Lexer) readString() token.Token {
	p := l.pos + 1
	t := token.Token{Type: token.String}

	for {
		l.readChar()
		if l.ch == 0 {
			t.Type = token.Illegal
			break
		}
		if l.ch == '"' {
			break
		}
	}

	t.Literal = l.input[p:l.pos]
	return t
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isDot(ch byte) bool {
	return ch == '.'
}
