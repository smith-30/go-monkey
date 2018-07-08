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

func testEval(i string) object.Object {
	l := lexer.New(i)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}
