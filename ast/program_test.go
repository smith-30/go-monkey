package ast

import (
	"testing"

	"github.com/smith-30/go-monkey/token"
)

func TestProgram_String(t *testing.T) {
	type fields struct {
		Statements []Statement
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "let ok",
			fields: fields{
				Statements: []Statement{
					&LetStatement{
						Token: token.Token{Type: token.LET, Literal: "let"},
						Name: &Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "myVar"},
							Value: "myVar",
						},
						Value: &Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
							Value: "anotherVar",
						},
					},
				},
			},
			want: "let myVar = anotherVar;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Program{
				Statements: tt.fields.Statements,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("Program.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
