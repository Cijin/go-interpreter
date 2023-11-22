package lexer

import (
	"errors"

	"github.com/cijin/go-interpreter/token"
)

type Lexer struct {
	input        string
	ch           byte
	position     int
	readPosition int
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for null
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peakChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readString() (string, error) {
	// current position on '"'
	position := l.position + 1

	// keep reading till '"'
	for {
		l.readChar()

		if l.ch == 0 || l.ch == '\n' {
			return "", errors.New("string literal not terminated")
		}

		if l.ch == '"' {
			break
		}

		// escape '"', not sure how to do this yet
		// also look into \n, \r, \t
		// if l.ch == '\'
	}

	return l.input[position:l.position], nil
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peakChar() == '=' {
			ch := l.ch
			l.readChar()

			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}

	case '!':
		if l.peakChar() == '=' {
			ch := l.ch
			l.readChar()

			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}

	case '"':
		tok.Type = token.STRING
		tok.Literal, tok.Error = l.readString()

	case '<':
		tok = newToken(token.LT, l.ch)

	case '>':
		tok = newToken(token.GT, l.ch)

	case ',':
		tok = newToken(token.COMMA, l.ch)

	case ';':
		tok = newToken(token.SEMICOLON, l.ch)

	case '+':
		tok = newToken(token.PLUS, l.ch)

	case '-':
		tok = newToken(token.MINUS, l.ch)

	case '*':
		tok = newToken(token.ASTERISK, l.ch)

	case '/':
		tok = newToken(token.SLASH, l.ch)

	case '(':
		tok = newToken(token.LPAREN, l.ch)

	case ')':
		tok = newToken(token.RPAREN, l.ch)

	case '{':
		tok = newToken(token.LSQUIRLY, l.ch)

	case '}':
		tok = newToken(token.RSQUIRLY, l.ch)

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)

			// there is no need to readChar as readIdentifier
			// reads past the identifier
			return tok
		}

		if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT

			return tok
		}

		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
