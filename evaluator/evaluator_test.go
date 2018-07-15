package evaluator

import (
	"testing"

	"github.com/smith-30/go-monkey/lexer"
	"github.com/smith-30/go-monkey/object"
	"github.com/smith-30/go-monkey/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		name  string
		input string
		exp   int64
	}{
		{"1", "5", 5},
		{"2", "10", 10},
		{"3", "-5", -5},
		{"4", "-10", -10},
		{"5", "5 + 5 + 5 - 10", 5},
		{"6", "2 * 2 * 2 * 2 * 2", 32},
		{"7", "-50 + 100 + -50", 0},
		{"8", "5 * 2 + 10", 20},
		{"9", "5 + 2 * 10", 25},
		{"10", "20 + 2 * -10", 0},
		{"11", "50 / 2 * 2 + 10", 60},
		{"12", "2 * (5 + 10)", 30},
		{"13", "3 * 3 * 3 + 10", 37},
		{"14", "3 * (3 * 3) + 10 ", 37},
		{"15", "(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			testIntegerObject(t, evaluated, tt.exp)
		})
	}
}

func testIntegerObject(t *testing.T, obj object.Object, exp int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != exp {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, exp)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		name  string
		input string
		exp   bool
	}{
		{"1", "false", false},
		{"2", "true", true},
		{"3", "1 < 2", true},
		{"4", "1 > 2", false},
		{"5", "1 < 1", false},
		{"6", "1 > 1", false},
		{"7", "1 == 1", true},
		{"8", "1 != 1", false},
		{"9", "1 == 2", false},
		{"10", "1 != 2", true},
		{"11", "true == true", true},
		{"12", "false == false", true},
		{"13", "true == false", false},
		{"14", "true != false", true},
		{"15", "false != true", true},
		{"16", "(1 < 2) == true", true},
		{"17", "(1 < 2) == false", false},
		{"18", "(1 > 2) == true", false},
		{"19", "(1 > 2) == false", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			testBooleanObject(t, evaluated, tt.exp)
		})
	}
}

func testBooleanObject(t *testing.T, obj object.Object, exp bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != exp {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, exp)
		return false
	}
	return true
}

func TestEvalBangOperator(t *testing.T) {
	tests := []struct {
		name  string
		input string
		exp   bool
	}{
		{"1", "!false", true},
		{"2", "!true", false},
		{"3", "!5", false},
		{"4", "!!true", true},
		{"5", "!!false", false},
		{"6", "!!5", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			testBooleanObject(t, evaluated, tt.exp)
		})
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		name  string
		input string
		exp   interface{}
	}{
		{"1", "if (true) { 10 }", 10},
		{"2", "if (false) { 10 }", nil},
		{"3", "if (1) { 10 }", 10},
		{"4", "if (1 < 2) { 10 }", 10},
		{"5", "if (1 > 2) { 10 }", nil},
		{"6", "if (1 < 2) { 10 } else { 20 }", 10},
		{"7", "if (1 > 2) { 10 } else { 20 }", 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			integer, ok := tt.exp.(int)
			if ok {
				testIntegerObject(t, evaluated, int64(integer))
			} else {
				testNullObject(t, evaluated)
			}

		})
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		name  string
		input string
		exp   int64
	}{
		{"1", "return 10;", 10},
		{"2", "return 10; 9;", 10},
		{"3", "return 2 * 5; 9;", 10},
		{"4", "9; return 2 * 5; 9;", 10},
		{
			"5",
			`
if (10 > 1) {
	if (10 > 1) {
		return 10;
	}

	return 1;
}
			`,
			10,
		},
		{
			"6",
			`
let f = fn(x) {
  return x;
  x + 10;
};
f(10);`,
			10,
		},
		{
			"7",
			`
let f = fn(x) {
   let result = x + 10;
   return result;
   return 10;
};
f(10);`,
			20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			testIntegerObject(t, evaluated, tt.exp)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name  string
		input string
		exp   string
	}{
		{
			"1",
			`5 + true;`,
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"2",
			`5 + true; 5;`,
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"3",
			`-true`,
			"unknown operator: -BOOLEAN",
		},
		{
			"4",
			`true + false;`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5",
			`5; true + false; 5;`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"6",
			`if (10 > 1) { true + false; }`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"7",
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"8",
			`foobar;`,
			"identifier not found: foobar",
		},
		{
			"9",
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			"10",
			`{"name": "Mokey"}[fn(x) { x }];`,
			"unusable as hash key: FUNCTION",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			errObj, ok := evaluated.(*object.Error)

			if !ok {
				t.Fatalf("no error object returned. got=%T (%#v)", evaluated, evaluated)
			}

			if errObj.Message != tt.exp {
				t.Errorf("want %q, but %q", tt.exp, errObj.Message)
			}
		})
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		name  string
		input string
		exp   int64
	}{
		{"1", "let a = 5; a;", 5},
		{"2", "let a = 5 * 5; a;", 25},
		{"3", "let a = 5; let b = a; b;", 5},
		{"4", "let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			testIntegerObject(t, evaluated, tt.exp)
		})
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2 };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x, got=%q", fn.Parameters[0])
	}

	expBody := "(x + 2)"
	if fn.Body.String() != expBody {
		t.Fatalf("Body is not %q. got=%q", expBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		name  string
		input string
		exp   int64
	}{
		{
			input: "let identity = fn(x) { x; }; identity(5);",
			exp:   5,
		},
		{
			input: "let identity = fn(x) { return x; }; identity(5);",
			exp:   5,
		},
		{
			input: "let double = fn(x) { x * 2; }; double(5)",
			exp:   10,
		},
		{
			input: "let add = fn(x, y) { x + y; }; add(5, 5)",
			exp:   10,
		},
		{
			input: "let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5))",
			exp:   20,
		},
		{
			input: "fn(x) { x; }(5)",
			exp:   5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testIntegerObject(t, testEval(tt.input), tt.exp)
		})
	}
}

func TestClosures(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			input: `
			let newAdder = fn(x) {
				fn(y) {x + y};
			};

			let addTwo = newAdder(2);
			addTwo(2);
			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testIntegerObject(t, testEval(tt.input), 4)
		})
	}
}

func TestStringLiteral(t *testing.T) {
	tests := []struct {
		name  string
		input string
		exp   string
	}{
		{
			input: `"Hello World!;"`,
			exp:   "Hello World!;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			s, ok := evaluated.(*object.String)
			if !ok {
				t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
			}

			if s.Value != tt.exp {
				t.Errorf("string has wrong value. got=%q", s.Value)
			}
		})
	}
}

func TestStringConcat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		exp   string
	}{
		{
			input: `"Hello" + " " + "World!"`,
			exp:   "Hello World!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			str, ok := evaluated.(*object.String)
			if !ok {
				t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
			}

			if str.Value != tt.exp {
				t.Errorf("String has wrong value. got=%q", str.Value)
			}
		})
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		name  string
		input string
		exp   interface{}
	}{
		{
			input: `len("")`,
			exp:   0,
		},
		{
			input: `len("four")`,
			exp:   4,
		},
		{
			input: `len("hello world")`,
			exp:   11,
		},
		{
			input: `len(1)`,
			exp:   "argument to `len` not supported, got INTEGER",
		},
		{
			input: `len("one", "two")`,
			exp:   "wrong number of arguments. got=2, want=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch exp := tt.exp.(type) {
			case int:
				testIntegerObject(t, evaluated, int64(exp))
			case string:
				errObj, ok := evaluated.(*object.Error)
				if !ok {
					t.Fatalf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				}
				if errObj.Message != tt.exp {
					t.Errorf("exp=%q. got=%q", tt.exp, errObj.Message)
				}
			}
		})
	}
}

func Test(t *testing.T) {
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
			evaluated := testEval(tt.input)

			arr, ok := evaluated.(*object.Array)
			if !ok {
				t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
			}

			if len(arr.Elements) != 3 {
				t.Fatalf("exp %d, but %d", 3, len(arr.Elements))
			}

			testIntegerObject(t, arr.Elements[0], 1)
			testIntegerObject(t, arr.Elements[1], 4)
			testIntegerObject(t, arr.Elements[2], 6)
		})
	}
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
		exp   interface{}
	}{
		{
			input: "[1, 2, 3][0]",
			exp:   1,
		},
		{
			input: "[1, 2, 3][1]",
			exp:   2,
		},
		{
			input: "[1, 2, 3][2]",
			exp:   3,
		},
		{
			input: "let i = 0; [1][i]",
			exp:   1,
		},
		{
			input: "[1, 2, 3][1 + 1]",
			exp:   3,
		},
		{
			input: "let myArray = [1, 2, 3]; myArray[2];",
			exp:   3,
		},
		{
			input: "let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2]",
			exp:   6,
		},
		{
			input: "let myArray= [1, 2, 3]; let i = myArray[0]; myArray[i]",
			exp:   2,
		},
		{
			input: "[1, 2, 3][3]",
			exp:   nil,
		},
		{
			input: "[1, 2, 3][-1]",
			exp:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			integer, ok := tt.exp.(int)

			if ok {
				testIntegerObject(t, evaluated, int64(integer))
			} else {
				testNullObject(t, evaluated)
			}
		})
	}
}

func TestHashLiteralList(t *testing.T) {
	tests := []struct {
		name  string
		input string
		exp   map[object.HashKey]int64
	}{
		{
			input: `let two = "two";
			{
				"one": 10 - 9,
				two: 1 + 1,
				"thr" + "ee": 6 / 2,
				4: 4,
				true: 5,
				false: 6
			}`,
			exp: map[object.HashKey]int64{
				(&object.String{Value: "one"}).HashKey():   1,
				(&object.String{Value: "two"}).HashKey():   2,
				(&object.String{Value: "three"}).HashKey(): 3,
				(&object.Integer{Value: 4}).HashKey():      4,
				TRUE.HashKey():                             5,
				FALSE.HashKey():                            6,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			result, ok := evaluated.(*object.Hash)
			if !ok {
				t.Fatalf("object is not Hash. got=%T (%+v)", evaluated, evaluated)
			}

			if len(result.Pairs) != len(tt.exp) {
				t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
			}

			for expKey, expVal := range tt.exp {
				pair, ok := result.Pairs[expKey]
				if !ok {
					t.Errorf("no pair for given key in Pairs")
				}
				testIntegerObject(t, pair.Value, expVal)
			}

		})
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		name  string
		input string
		exp   interface{}
	}{
		{
			input: `{"foo": 5}["foo"]`,
			exp:   5,
		},
		{
			input: `{"foo": 5}["bar"]`,
			exp:   nil,
		},
		{
			input: `let key = "foo"; {"foo": 5}[key]`,
			exp:   5,
		},
		{
			input: `{}["foo"]`,
			exp:   nil,
		},
		{
			input: `{5: 5}[5]`,
			exp:   5,
		},
		{
			input: `{true: 5}[true]`,
			exp:   5,
		},
		{
			input: `{false: 5}[false]`,
			exp:   5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			integer, ok := tt.exp.(int)
			if ok {
				testIntegerObject(t, evaluated, int64(integer))
			} else {
				testNullObject(t, evaluated)
			}
		})
	}
}

func testEval(i string) object.Object {
	l := lexer.New(i)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}
