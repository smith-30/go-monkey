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

func TestIfExpression(t *testing.T) {
	type fields struct {
		input string
	}
	type exp struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		exp    exp
	}{
		{
			name: "if (x < y) { x }",
			fields: fields{
				input: `if (x < y) { x };`,
			},
			exp: exp{
				value: "if (x < y) { x }",
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

			if len(program.Statements) != 1 {
				t.Fatalf("Program.Statements does not contain %d statements. got = %d", 1, len(program.Statements))
			}

			for _, stmt := range program.Statements {
				expStmt, ok := stmt.(*ast.ExpressionStatement)
				if !ok {
					t.Errorf("expStmt not *ast.ExpressionStatement. got=%T", stmt)
					continue
				}

				ident, ok := expStmt.Expression.(*ast.IfExpression)
				if !ok {
					t.Errorf("expStmt not *ast.IfExpression. got=%T", expStmt.Expression)
					continue
				}

				if !testInfixExpression(t, ident.Condition, "x", "<", "y") {
					return
				}

				if len(ident.Consequence.Statements) != 1 {
					t.Errorf("Consequence does not contain %d statements. got = %d", 1, len(ident.Consequence.Statements))
				}

				consequence, ok := ident.Consequence.Statements[0].(*ast.ExpressionStatement)
				if !ok {
					t.Errorf("Statements[0] is not *ast.ExpressionStatement. got=%T", expStmt.Expression)
				}

				if !testIdentifier(t, consequence.Expression, "x") {
					return
				}

				if ident.Alternative != nil {
					t.Errorf("ident.Alternative.Statemants is not nil. got=%#v", ident.Alternative)
				}
			}
		})
	}
}

func TestIfElseExpression(t *testing.T) {
	type fields struct {
		input string
	}
	type exp struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		exp    exp
	}{
		{
			name: "if (x < y) { x } else { y }",
			fields: fields{
				input: `if (x < y) { x } else { y };`,
			},
			exp: exp{
				value: "if (x < y) { x } else { y }",
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

			if len(program.Statements) != 1 {
				t.Fatalf("Program.Statements does not contain %d statements. got = %d", 1, len(program.Statements))
			}

			for _, stmt := range program.Statements {
				expStmt, ok := stmt.(*ast.ExpressionStatement)
				if !ok {
					t.Errorf("expStmt not *ast.ExpressionStatement. got=%T", stmt)
					continue
				}

				ident, ok := expStmt.Expression.(*ast.IfExpression)
				if !ok {
					t.Errorf("expStmt not *ast.IfExpression. got=%T", expStmt.Expression)
					continue
				}

				if !testInfixExpression(t, ident.Condition, "x", "<", "y") {
					return
				}

				if len(ident.Consequence.Statements) != 1 {
					t.Errorf("Consequence does not contain %d statements. got = %d", 1, len(ident.Consequence.Statements))
				}

				consequence, ok := ident.Consequence.Statements[0].(*ast.ExpressionStatement)
				if !ok {
					t.Errorf("Statements[0] is not *ast.ExpressionStatement. got=%T", ident.Consequence.Statements[0])
				}

				if !testIdentifier(t, consequence.Expression, "x") {
					return
				}

				elseExp, ok := ident.Alternative.Statements[0].(*ast.ExpressionStatement)
				if !ok {
					t.Errorf("Statements[0] is not *ast.ExpressionStatement. got=%T", ident.Alternative.Statements[0])
				}

				if !testIdentifier(t, elseExp.Expression, "y") {
					return
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
					t.Errorf("expStmt not *ast.ExpressionStatement. got=%T", stmt)
					continue
				}

				literal, ok := expStmt.Expression.(*ast.IntegerLiteral)
				if !ok {
					t.Errorf("expStmt not *ast.IntegerLiteral. got=%T", expStmt.Expression)
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

func TestBooleanExpression(t *testing.T) {
	type fields struct {
		input string
	}
	type exp struct {
		wantStmtCount int
		literal       string
		value         bool
	}
	tests := []struct {
		name   string
		fields fields
		exp    exp
	}{
		{
			name: "false",
			fields: fields{
				input: `false;`,
			},
			exp: exp{
				wantStmtCount: 1,
				literal:       "false",
				value:         false,
			},
		},
		{
			name: "true",
			fields: fields{
				input: `true;`,
			},
			exp: exp{
				wantStmtCount: 1,
				literal:       "true",
				value:         true,
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
					t.Errorf("expStmt not *ast.ExpressionStatement. got=%T", stmt)
					continue
				}

				literal, ok := expStmt.Expression.(*ast.Boolean)
				if !ok {
					t.Errorf("expStmt not *ast.IntegerLiteral. got=%T", expStmt.Expression)
					continue
				}
				if literal.Value != tt.exp.value {
					t.Errorf("literal.Value is not %v, got %v", tt.exp.value, literal.Value)
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
		value         interface{}
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
		{
			name: "`!true`",
			fields: fields{
				input: `!true;`,
			},
			exp: exp{
				wantStmtCount: 1,
				operator:      "!",
				value:         true,
			},
		},
		{
			name: "`!false`",
			fields: fields{
				input: `!false;`,
			},
			exp: exp{
				wantStmtCount: 1,
				operator:      "!",
				value:         false,
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

				if !testLiteralExpression(t, prefixExp.Right, tt.exp.value) {
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

func TestParsingInfixExpressions(t *testing.T) {
	type fields struct {
		input string
	}
	type exp struct {
		wantStmtCount int
		leftValue     interface{}
		operator      string
		rightValue    interface{}
	}
	tests := []struct {
		name   string
		fields fields
		exp    exp
	}{
		{
			name:   "`+`",
			fields: fields{input: `5 + 5;`},
			exp: exp{
				wantStmtCount: 1,
				leftValue:     5,
				operator:      "+",
				rightValue:    5,
			},
		},
		{
			name:   "`-`",
			fields: fields{input: `5 - 5;`},
			exp: exp{
				wantStmtCount: 1,
				leftValue:     5,
				operator:      "-",
				rightValue:    5,
			},
		},
		{
			name:   "`*`",
			fields: fields{input: `5 * 5;`},
			exp: exp{
				wantStmtCount: 1,
				leftValue:     5,
				operator:      "*",
				rightValue:    5,
			},
		},
		{
			name:   "`/`",
			fields: fields{input: `5 / 5;`},
			exp: exp{
				wantStmtCount: 1,
				leftValue:     5,
				operator:      "/",
				rightValue:    5,
			},
		},
		{
			name:   "`>`",
			fields: fields{input: `5 > 5;`},
			exp: exp{
				wantStmtCount: 1,
				leftValue:     5,
				operator:      ">",
				rightValue:    5,
			},
		},
		{
			name:   "`<`",
			fields: fields{input: `5 < 5;`},
			exp: exp{
				wantStmtCount: 1,
				leftValue:     5,
				operator:      "<",
				rightValue:    5,
			},
		},
		{
			name:   "`==`",
			fields: fields{input: `5 == 5;`},
			exp: exp{
				wantStmtCount: 1,
				leftValue:     5,
				operator:      "==",
				rightValue:    5,
			},
		},
		{
			name:   "`!=`",
			fields: fields{input: `5 != 5;`},
			exp: exp{
				wantStmtCount: 1,
				leftValue:     5,
				operator:      "!=",
				rightValue:    5,
			},
		},
		{
			name:   "`true == true`",
			fields: fields{input: `true == true`},
			exp: exp{
				wantStmtCount: 1,
				leftValue:     true,
				operator:      "==",
				rightValue:    true,
			},
		},
		{
			name:   "`true != false`",
			fields: fields{input: `true != false`},
			exp: exp{
				wantStmtCount: 1,
				leftValue:     true,
				operator:      "!=",
				rightValue:    false,
			},
		},
		{
			name:   "`false == false`",
			fields: fields{input: `false == false`},
			exp: exp{
				wantStmtCount: 1,
				leftValue:     false,
				operator:      "==",
				rightValue:    false,
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

				prefixExp, ok := expStmt.Expression.(*ast.InfixExpression)
				if !ok {
					t.Errorf("expStmt not *ast.PrefixExpression. got=%T", expStmt.Expression)
					continue
				}

				if !testInfixExpression(t, prefixExp, tt.exp.leftValue, tt.exp.operator, tt.exp.rightValue) {
					return
				}
			}
		})
	}
}

func TestFunctionLiteralParsing(t *testing.T) {

}

func TestOperatorPresedenceParsing(t *testing.T) {
	type fields struct {
		input string
	}
	type exp struct {
		val string
	}
	tests := []struct {
		name   string
		fields fields
		exp    exp
	}{
		{
			name:   "`-a * b`",
			fields: fields{input: `-a * b`},
			exp:    exp{val: "((-a) * b)"},
		},
		{
			name:   "`!-a`",
			fields: fields{input: `!-a`},
			exp:    exp{val: "(!(-a))"},
		},
		{
			name:   "`a + b + c`",
			fields: fields{input: `a + b + c`},
			exp:    exp{val: "((a + b) + c)"},
		},
		{
			name:   "`a + b - c`",
			fields: fields{input: `a + b - c`},
			exp:    exp{val: "((a + b) - c)"},
		},
		{
			name:   "`a * b * c`",
			fields: fields{input: `a * b * c`},
			exp:    exp{val: "((a * b) * c)"},
		},
		{
			name:   "`a * b / c`",
			fields: fields{input: `a * b / c`},
			exp:    exp{val: "((a * b) / c)"},
		},
		{
			name:   "`a + b * c + d / e - f `",
			fields: fields{input: `a + b * c + d / e - f`},
			exp:    exp{val: "(((a + (b * c)) + (d / e)) - f)"},
		},
		{
			name:   "`3 + 4; -5 * 5`",
			fields: fields{input: `3 + 4; -5 * 5`},
			exp:    exp{val: "(3 + 4)((-5) * 5)"},
		},
		{
			name:   "`5 > 4 == 3 < 4`",
			fields: fields{input: `5 > 4 == 3 < 4`},
			exp:    exp{val: "((5 > 4) == (3 < 4))"},
		},
		{
			name:   "`5 < 4 != 3 > 4`",
			fields: fields{input: `5 < 4 != 3 > 4`},
			exp:    exp{val: "((5 < 4) != (3 > 4))"},
		},
		{
			name:   "`3 + 4 * 5 == 3 * 1 + 4 * 5`",
			fields: fields{input: `3 + 4 * 5 == 3 * 1 + 4 * 5`},
			exp:    exp{val: "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		},
		{
			name:   "true",
			fields: fields{input: `true`},
			exp:    exp{val: "true"},
		},
		{
			name:   "true",
			fields: fields{input: `false`},
			exp:    exp{val: "false"},
		},
		{
			name:   "3 > 5 == false",
			fields: fields{input: `3 > 5 == false`},
			exp:    exp{val: "((3 > 5) == false)"},
		},
		{
			name:   "3 < 5 == true",
			fields: fields{input: `3 < 5 == true`},
			exp:    exp{val: "((3 < 5) == true)"},
		},
		{
			name:   "1 + (2 + 3) + 4",
			fields: fields{input: `1 + (2 + 3) + 4`},
			exp:    exp{val: "((1 + (2 + 3)) + 4)"},
		},
		{
			name:   "(5 + 5) * 2",
			fields: fields{input: `(5 + 5) * 2`},
			exp:    exp{val: "((5 + 5) * 2)"},
		},
		{
			name:   "2 / (5 + 5)",
			fields: fields{input: `2 / (5 + 5)`},
			exp:    exp{val: "(2 / (5 + 5))"},
		},
		{
			name:   "-(5 + 5)",
			fields: fields{input: `-(5 + 5)`},
			exp:    exp{val: "(-(5 + 5))"},
		},
		{
			name:   "!(true == true)",
			fields: fields{input: `!(true == true)`},
			exp:    exp{val: "(!(true == true))"},
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

			act := program.String()
			if act != tt.exp.val {
				t.Errorf("exp=%q, got=%q", tt.exp.val, act)
			}
		})
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, val string) (result bool) {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return result
	}

	if ident.Value != val {
		t.Errorf("ident.Value not %s, got=%s", val, ident.Value)
		return result
	}

	if ident.TokenLiteral() != val {
		t.Errorf("ident.TokenLiteral() not %s, got=%s", val, ident.TokenLiteral())
		return result
	}

	result = true
	return result
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	b, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if b.Value != value {
		t.Errorf("b.Value %t, value %t", b.Value, value)
		return false
	}

	if b.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("b.TokenLiteral %s, value %t", b.TokenLiteral(), value)
	}

	return true
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) (result bool) {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return result
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return result
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return result
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return result
	}

	result = true
	return result
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
