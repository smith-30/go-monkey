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
	type wantVal struct {
		val interface{}
	}

	tests := []struct {
		name     string
		fields   fields
		wants    []want
		wantVals []wantVal
	}{
		{
			name: "let",
			fields: fields{
				input: `
let x = 5;
let y = true;
let foobar = y;
`,
			},
			wants: []want{
				{"x"},
				{"y"},
				{"foobar"},
			},
			wantVals: []wantVal{
				{5},
				{true},
				{"y"},
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

				val := stmt.(*ast.LetStatement).Value
				if !testLiteralExpression(t, val, tt.wantVals[i].val) {
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
	type wantVal struct {
		val interface{}
	}
	tests := []struct {
		name     string
		fields   fields
		want     int
		wantVals []wantVal
	}{
		{
			name: "return",
			fields: fields{
				input: `
return 5;
return true;
return y;
`,
			},
			want: 3,
			wantVals: []wantVal{
				{5},
				{true},
				{"y"},
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

			if len(program.Statements) != tt.want {
				t.Fatalf("Program.Statements does not contain %d statements. got = %d", tt.want, len(program.Statements))
			}

			for i, stmt := range program.Statements {
				returnStmt, ok := stmt.(*ast.ReturnStatement)
				if !ok {
					t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
					continue
				}
				if returnStmt.TokenLiteral() != "return" {
					t.Errorf("returnStmt.TokenLiteral is not 'return', got %q", returnStmt.TokenLiteral())
				}
				val := returnStmt.ReturnValue
				if !testLiteralExpression(t, val, tt.wantVals[i].val) {
					return
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

func TestFunctionLiteralParsing(t *testing.T) {
	type fields struct {
		input string
	}
	type exp struct {
		value      string
		paramCount int
	}
	tests := []struct {
		name   string
		fields fields
		exp    exp
	}{
		{
			name: "fn(x, y) { x + y; }",
			fields: fields{
				input: `fn(x, y) { x + y; }`,
			},
			exp: exp{
				value:      "fn(x, y) { x + y; }",
				paramCount: 2,
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

				ident, ok := expStmt.Expression.(*ast.FunctionLiteral)
				if !ok {
					t.Errorf("expStmt not *ast.IfExpression. got=%T", expStmt.Expression)
					continue
				}

				if len(ident.Parameters) != tt.exp.paramCount {
					t.Errorf("Parameters does not contain %d. got = %d", tt.exp.paramCount, len(ident.Parameters))
				}

				// Todo attach exp field
				testLiteralExpression(t, ident.Parameters[0], "x")
				testLiteralExpression(t, ident.Parameters[1], "y")

				if len(ident.Body.Statements) != 1 {
					t.Fatalf("ident.Body.Statements does not contain %d. got = %d", 1, len(ident.Body.Statements))
				}

				bodyStmt, ok := ident.Body.Statements[0].(*ast.ExpressionStatement)
				if !ok {
					t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", ident.Body.Statements[0])
				}

				testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
			}
		})
	}
}

func TestFunctionParameterParsing(t *testing.T) {
	type fields struct {
		input string
	}
	type exp struct {
		values []string
	}
	tests := []struct {
		name   string
		fields fields
		exp    exp
	}{
		{
			name: "fn() {};",
			fields: fields{
				input: `fn() {};`,
			},
			exp: exp{
				values: []string{},
			},
		},
		{
			name: "fn(x) {};",
			fields: fields{
				input: `fn(x) {};`,
			},
			exp: exp{
				values: []string{"x"},
			},
		},
		{
			name: "fn(x, y, z) {}",
			fields: fields{
				input: `fn(x, y, z) {}`,
			},
			exp: exp{
				values: []string{"x", "y", "z"},
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

			for _, stmt := range program.Statements {
				expStmt, ok := stmt.(*ast.ExpressionStatement)
				if !ok {
					t.Errorf("expStmt not *ast.ExpressionStatement. got=%T", stmt)
					continue
				}

				ident, ok := expStmt.Expression.(*ast.FunctionLiteral)
				if !ok {
					t.Errorf("expStmt not *ast.IfExpression. got=%T", expStmt.Expression)
					continue
				}

				if len(ident.Parameters) != len(tt.exp.values) {
					t.Errorf("Parameters does not contain %d. got = %d", len(tt.exp.values), len(ident.Parameters))
				}

				for i, idt := range tt.exp.values {
					testLiteralExpression(t, ident.Parameters[i], idt)
				}
			}
		})
	}
}

func TestCallExpressionParsing(t *testing.T) {
	type fields struct {
		input string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "add(1, 2 * 3, 4 + 5)",
			fields: fields{
				input: `add(1, 2 * 3, 4 + 5)`,
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

				ident, ok := expStmt.Expression.(*ast.CallExpression)
				if !ok {
					t.Errorf("expStmt not *ast.CallExpression. got=%T", expStmt.Expression)
					continue
				}

				if !testIdentifier(t, ident.Function, "add") {
					return
				}

				if len(ident.Arguments) != 3 {
					t.Errorf("Arguments does not contain %d. got = %d", 3, len(ident.Arguments))
				}

				testLiteralExpression(t, ident.Arguments[0], 1)
				testInfixExpression(t, ident.Arguments[1], 2, "*", 3)
				testInfixExpression(t, ident.Arguments[2], 4, "+", 5)
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

func TestStringLiteralExpression(t *testing.T) {
	type fields struct {
		input string
	}
	type exp struct {
		wantStmtCount int
		literal       string
		value         string
	}
	tests := []struct {
		name   string
		fields fields
		exp    exp
	}{
		{
			fields: fields{
				input: `"hello world";`,
			},
			exp: exp{
				wantStmtCount: 1,
				literal:       "hello world",
				value:         "hello world",
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

				literal, ok := expStmt.Expression.(*ast.StringLiteral)
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

func TestParsingArrayLiterals(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			input: "[1, 2 * 2, 3 + 3]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			program := p.ParseProgram()
			checkParseErrors(t, p)

			stmt := program.Statements[0].(*ast.ExpressionStatement)
			arr, ok := stmt.Expression.(*ast.ArrayLiteral)

			if !ok {
				t.Fatalf("expStmt not *ast.ArrayLiteral. got=%T", stmt.Expression)
			}

			if len(arr.Elements) != 3 {
				t.Fatalf("len(arr.Elements) not 3. got=%d", len(arr.Elements))
			}

			testIntegerLiteral(t, arr.Elements[0], 1)
			testInfixExpression(t, arr.Elements[1], 2, "*", 2)
			testInfixExpression(t, arr.Elements[2], 3, "+", 3)
		})
	}
}

func TestParsingIndexExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			input: "myArray[1 + 1]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			program := p.ParseProgram()
			checkParseErrors(t, p)

			stmt := program.Statements[0].(*ast.ExpressionStatement)
			indexExp, ok := stmt.Expression.(*ast.IndexExpression)

			if !ok {
				t.Fatalf("expStmt not *ast.IndexExpression. got=%T", stmt.Expression)
			}

			if !testIdentifier(t, indexExp.Left, "myArray") {
				return
			}

			if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
				return
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

func TestOperatorPrecedenceParsing(t *testing.T) {
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
		{
			name:   "a + add(b * c) + d",
			fields: fields{input: `a + add(b * c) + d`},
			exp:    exp{val: "((a + add((b * c))) + d)"},
		},
		{
			name:   "add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			fields: fields{input: `add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))`},
			exp:    exp{val: "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		},
		{
			name:   "add(a + b + c * d/ f + g)",
			fields: fields{input: `add(a + b + c * d/ f + g)`},
			exp:    exp{val: "add((((a + b) + ((c * d) / f)) + g))"},
		},
		{
			name:   "a * [1, 2, 3, 4][b * c] * d",
			fields: fields{input: `a * [1, 2, 3, 4][b * c] * d`},
			exp:    exp{val: "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		},
		{
			name:   "add(a * b[2], b[1], 2 * [1, 2][1])",
			fields: fields{input: `add(a * b[2], b[1], 2 * [1, 2][1])`},
			exp:    exp{val: "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
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

	t.Errorf("type of exp not handled. got=%T, expected=%#v", exp, expected)
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
