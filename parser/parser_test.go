package parser

import (
	"fmt"
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

			if len(program.Statements) != len(tt.wants) {
				t.Fatalf("Program.Statements does not contain %d statements. got = %d", len(tt.wants), len(program.Statements))
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

func TestReturnStatements(t *testing.T) {
	type fields struct {
		input string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "return",
			fields: fields{
				input: `
return 5;
return 10;
return 838383;
`,
			},
			want: 3,
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

			if len(program.Statements) != tt.want {
				t.Fatalf("Program.Statements does not contain %d statements. got = %d", tt.want, len(program.Statements))
			}

			for _, stmt := range program.Statements {
				returnStmt, ok := stmt.(*ast.ReturnStatement)
				if !ok {
					t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
					continue
				}
				if returnStmt.TokenLiteral() != "return" {
					t.Errorf("returnStmt.TokenLiteral is not 'return', got %q", returnStmt.TokenLiteral())
				}
			}
		})
	}
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

func TestIdentifierExpression(t *testing.T) {
	type fields struct {
		input string
	}
	type exp struct {
		want  int
		value string
	}
	tests := []struct {
		name   string
		fields fields
		exp    exp
	}{
		{
			name: "simple identifier",
			fields: fields{
				input: `foobar;`,
			},
			exp: exp{
				want:  1,
				value: "foobar",
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

			if len(program.Statements) != tt.exp.want {
				t.Fatalf("Program.Statements does not contain %d statements. got = %d", tt.exp.want, len(program.Statements))
			}

			for _, stmt := range program.Statements {
				expStmt, ok := stmt.(*ast.ExpressionStatement)
				if !ok {
					t.Errorf("expStmt not *ast.ExpressionStatement. got=%T", stmt)
					continue
				}

				ident, ok := expStmt.Expression.(*ast.Identifier)
				if !ok {
					t.Errorf("expStmt not *ast.Identifier. got=%T", expStmt.Expression)
					continue
				}
				if ident.Value != tt.exp.value {
					t.Errorf("ident.Value is not %q, got %q", tt.exp.value, ident.Value)
				}

				if ident.TokenLiteral() != tt.exp.value {
					t.Errorf("ident.TokenLiteral is not %q, got %q", tt.exp.value, ident.TokenLiteral())
				}
			}
		})
	}

}

func TestIntegerLiteralExpression(t *testing.T) {
	type fields struct {
		input string
	}
	type exp struct {
		wantStmtCount int
		literal       string
		value         int64
	}
	tests := []struct {
		name   string
		fields fields
		exp    exp
	}{
		{
			name: "integer",
			fields: fields{
				input: `5;`,
			},
			exp: exp{
				wantStmtCount: 1,
				literal:       "5",
				value:         5,
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

			if len(program.Statements) != tt.exp.wantStmtCount {
				t.Fatalf("Program.Statements does not contain %d statements. got = %d", tt.exp.wantStmtCount, len(program.Statements))
			}

			for _, stmt := range program.Statements {
				expStmt, ok := stmt.(*ast.ExpressionStatement)
				if !ok {
					t.Errorf("expStmt not *ast.IntegerLiteral. got=%T", stmt)
					continue
				}

				literal, ok := expStmt.Expression.(*ast.IntegerLiteral)
				if !ok {
					t.Errorf("expStmt not *ast.Identifier. got=%T", expStmt.Expression)
					continue
				}
				if literal.Value != tt.exp.value {
					t.Errorf("literal.Value is not %q, got %q", tt.exp.value, literal.Value)
				}

				if literal.TokenLiteral() != tt.exp.literal {
					t.Errorf("literal.TokenLiteral is not %q, got %q", tt.exp.literal, literal.TokenLiteral())
				}
			}
		})
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	type fields struct {
		input string
	}
	type exp struct {
		wantStmtCount int
		operator      string
		value         int64
	}
	tests := []struct {
		name   string
		fields fields
		exp    exp
	}{
		{
			name: "`!`",
			fields: fields{
				input: `!5;`,
			},
			exp: exp{
				wantStmtCount: 1,
				operator:      "!",
				value:         5,
			},
		},
		{
			name: "`-`",
			fields: fields{
				input: `-15;`,
			},
			exp: exp{
				wantStmtCount: 1,
				operator:      "-",
				value:         15,
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

			if len(program.Statements) != tt.exp.wantStmtCount {
				t.Fatalf("Program.Statements does not contain %d statements. got = %d", tt.exp.wantStmtCount, len(program.Statements))
			}

			for _, stmt := range program.Statements {
				expStmt, ok := stmt.(*ast.ExpressionStatement)
				if !ok {
					t.Errorf("expStmt not *ast.IntegerLiteral. got=%T", stmt)
					continue
				}

				prefixExp, ok := expStmt.Expression.(*ast.PrefixExpression)
				if !ok {
					t.Errorf("expStmt not *ast.PrefixExpression. got=%T", expStmt.Expression)
					continue
				}
				if prefixExp.Operator != tt.exp.operator {
					t.Errorf("prefixExp.Value is not %q, got %q", tt.exp.operator, prefixExp.Operator)
				}

				if !testIntegerLiteral(t, prefixExp.Right, tt.exp.value) {
					return
				}
			}
		})
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value is not %d. got=%d", value, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.Value is not %d. got=%s", value, integer.TokenLiteral())
		return false
	}
	return true
}
