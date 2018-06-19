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
