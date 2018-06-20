package lexer

import (
	"reflect"
	"testing"

	"github.com/smith-30/go-monkey/token"
)

func TestLexer_NextToken(t *testing.T) {
	type fields struct {
		input string
	}
	tests := []struct {
		name   string
		fields fields
		want   []token.Token
	}{
		{
			name: "success",
			fields: fields{
				input: "=+(){},;",
			},
			want: []token.Token{
				token.Token{Type: token.ASSIGN, Literal: "="},
				token.Token{Type: token.PLUS, Literal: "+"},
				token.Token{Type: token.LPAREN, Literal: "("},
				token.Token{Type: token.RPAREN, Literal: ")"},
				token.Token{Type: token.LBRACE, Literal: "{"},
				token.Token{Type: token.RBRACE, Literal: "}"},
				token.Token{Type: token.COMMA, Literal: ","},
				token.Token{Type: token.SEMICOLON, Literal: ";"},
				token.Token{Type: token.EOF, Literal: ""},
			},
		},
		{
			name: "success source code",
			fields: fields{
				input: `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x+y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;

`,
			},
			want: []token.Token{
				token.Token{Type: token.LET, Literal: "let"},
				token.Token{Type: token.IDENT, Literal: "five"},
				token.Token{Type: token.ASSIGN, Literal: "="},
				token.Token{Type: token.INT, Literal: "5"},
				token.Token{Type: token.SEMICOLON, Literal: ";"},
				token.Token{Type: token.LET, Literal: "let"},
				token.Token{Type: token.IDENT, Literal: "ten"},
				token.Token{Type: token.ASSIGN, Literal: "="},
				token.Token{Type: token.INT, Literal: "10"},
				token.Token{Type: token.SEMICOLON, Literal: ";"},
				token.Token{Type: token.LET, Literal: "let"},
				token.Token{Type: token.IDENT, Literal: "add"},
				token.Token{Type: token.ASSIGN, Literal: "="},
				token.Token{Type: token.FUNCTION, Literal: "fn"},
				token.Token{Type: token.LPAREN, Literal: "("},
				token.Token{Type: token.IDENT, Literal: "x"},
				token.Token{Type: token.COMMA, Literal: ","},
				token.Token{Type: token.IDENT, Literal: "y"},
				token.Token{Type: token.RPAREN, Literal: ")"},
				token.Token{Type: token.LBRACE, Literal: "{"},
				token.Token{Type: token.IDENT, Literal: "x"},
				token.Token{Type: token.PLUS, Literal: "+"},
				token.Token{Type: token.IDENT, Literal: "y"},
				token.Token{Type: token.SEMICOLON, Literal: ";"},
				token.Token{Type: token.RBRACE, Literal: "}"},
				token.Token{Type: token.SEMICOLON, Literal: ";"},
				token.Token{Type: token.LET, Literal: "let"},
				token.Token{Type: token.IDENT, Literal: "result"},
				token.Token{Type: token.ASSIGN, Literal: "="},
				token.Token{Type: token.IDENT, Literal: "add"},
				token.Token{Type: token.LPAREN, Literal: "("},
				token.Token{Type: token.IDENT, Literal: "five"},
				token.Token{Type: token.COMMA, Literal: ","},
				token.Token{Type: token.IDENT, Literal: "ten"},
				token.Token{Type: token.RPAREN, Literal: ")"},
				token.Token{Type: token.SEMICOLON, Literal: ";"},
				token.Token{Type: token.BANG, Literal: "!"},
				token.Token{Type: token.MINUS, Literal: "-"},
				token.Token{Type: token.SLASH, Literal: "/"},
				token.Token{Type: token.ASTERISK, Literal: "*"},
				token.Token{Type: token.INT, Literal: "5"},
				token.Token{Type: token.SEMICOLON, Literal: ";"},
				token.Token{Type: token.INT, Literal: "5"},
				token.Token{Type: token.LT, Literal: "<"},
				token.Token{Type: token.INT, Literal: "10"},
				token.Token{Type: token.GT, Literal: ">"},
				token.Token{Type: token.INT, Literal: "5"},
				token.Token{Type: token.SEMICOLON, Literal: ";"},
				token.Token{Type: token.IF, Literal: "if"},
				token.Token{Type: token.LPAREN, Literal: "("},
				token.Token{Type: token.INT, Literal: "5"},
				token.Token{Type: token.LT, Literal: "<"},
				token.Token{Type: token.INT, Literal: "10"},
				token.Token{Type: token.RPAREN, Literal: ")"},
				token.Token{Type: token.LBRACE, Literal: "{"},
				token.Token{Type: token.RETURN, Literal: "return"},
				token.Token{Type: token.TRUE, Literal: "true"},
				token.Token{Type: token.SEMICOLON, Literal: ";"},
				token.Token{Type: token.RBRACE, Literal: "}"},
				token.Token{Type: token.ELSE, Literal: "else"},
				token.Token{Type: token.LBRACE, Literal: "{"},
				token.Token{Type: token.RETURN, Literal: "return"},
				token.Token{Type: token.FALSE, Literal: "false"},
				token.Token{Type: token.SEMICOLON, Literal: ";"},
				token.Token{Type: token.RBRACE, Literal: "}"},
				token.Token{Type: token.INT, Literal: "10"},
				token.Token{Type: token.EQ, Literal: "=="},
				token.Token{Type: token.INT, Literal: "10"},
				token.Token{Type: token.SEMICOLON, Literal: ";"},
				token.Token{Type: token.INT, Literal: "10"},
				token.Token{Type: token.NOT_EQ, Literal: "!="},
				token.Token{Type: token.INT, Literal: "9"},
				token.Token{Type: token.SEMICOLON, Literal: ";"},
				token.Token{Type: token.EOF, Literal: ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.fields.input)
			for _, token := range tt.want {
				if got := l.NextToken(); !reflect.DeepEqual(got, token) {
					t.Errorf("\nact  %#v\nwant %#v", got, token)
					return
				}
			}
		})
	}
}
