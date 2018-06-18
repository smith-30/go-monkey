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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.fields.input)
			for _, token := range tt.want {
				if got := l.NextToken(); !reflect.DeepEqual(got, token) {
					t.Errorf("\nLexer.NextToken() = %#v, want %#v", got, token)
				}
			}
		})
	}
}
