package lexer

import "github.com/smith-30/go-monkey/token"

type Lexer struct {
	input        string
	position     int  // now position about input.(indicate now character)
	readPosition int  // position to read from now.(next to the current character)
	ch           byte // character currently being inspected
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()
	return l
}

// Todo: parsable whole Unicode. This func have not supported all Unicode yet.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// 0 corresponds to ASCII's NUL
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		t = token.NewToken(token.ASSIGN, l.ch)
	case ';':
		t = token.NewToken(token.SEMICOLON, l.ch)
	case ')':
		t = token.NewToken(token.RPAREN, l.ch)
	case '(':
		t = token.NewToken(token.LPAREN, l.ch)
	case '{':
		t = token.NewToken(token.LBRACE, l.ch)
	case '}':
		t = token.NewToken(token.RBRACE, l.ch)
	case ',':
		t = token.NewToken(token.COMMA, l.ch)
	case '+':
		t = token.NewToken(token.PLUS, l.ch)
	case 0:
		t.Literal = ""
		t.Type = token.EOF
	default:
		if isLetter(l.ch) {
			t.Literal = l.readIdentifier()
			t.Type = token.LookUpIdent(t.Literal)
			return t
		} else if isDigit(l.ch) {
			t.Type = token.INT
			t.Literal = l.readNumber()
			return t
		} else {
			t = token.NewToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return t
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// This func can parse only int.
// Todo parse float? hexadecimal?
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
