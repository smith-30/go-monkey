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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			testIntegerObject(t, evaluated, tt.exp)
		})
	}
}

func testEval(i string) object.Object {
	l := lexer.New(i)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}
