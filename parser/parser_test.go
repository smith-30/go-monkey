package parser

import (
	"testing"

	"github.com/smith-30/go-monkey/ast"
	"github.com/smith-30/go-monkey/lexer"
)

func TestLetStatement(t *testing.T) {
	type fields struct {
		input string
	}
	type want struct {
		expectedIdentifier string
	}

	tests := []struct {
		name   string
		fields fields
		wants  []want
	}{
		{
			name: "let",
			fields: fields{
				input: `
let x = 5;
let y = 10;
let foobar = 838383;
`,
			},
			wants: []want{
				{"x"},
				{"y"},
				{"foobar"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.fields.input)
			p := New(l)

			program := p.ParseProgram()
			checkParseErrors(t, p)
			if program == nil {
				t.Fatalf("ParseProgram() returned nil")
			}

			if len(program.Statements) != 3 {
				t.Fatalf("Program.Statements does not contain 3 statements. got = %d", len(program.Statements))
			}

			for i, w := range tt.wants {
				stmt := program.Statements[i]
				if !testLetStatement(t, stmt, w.expectedIdentifier) {
					return
				}
			}
		})
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) (result bool) {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got = %q", s.TokenLiteral())
		return
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got = %T", s)
		return
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got = '%s'", name, letStmt.Name.Value)
		return
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got = '%s'", name, letStmt.Name.TokenLiteral())
		return
	}

	result = true
	return
}

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
