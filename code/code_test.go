package code

import (
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{
			op:       OpConstant,
			operands: []int{65534},
			expected: []byte{byte(OpConstant), 255, 254},
		},
		{
			op:       OpConstant,
			operands: []int{32767},
			expected: []byte{byte(OpConstant), 127, 255},
		},
	}
	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length. want=%d, got=%d",
				len(tt.expected), len(instruction))
		}
		for i, b := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Errorf("wrong byte at pos %d. want=%d, got=%d", i, b, instruction[i])
			}
		}
	}
}
