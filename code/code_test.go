package code

import (
	"testing"
)

const (
	NOT_EXPECTED = "expected %v to be %v, found %v"
	WRONG_LENGTH = "%v has wrong length: expected=%v, found=%v"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
	}
	for _, tt := range tests {
		instructions := Make(tt.op, tt.operands...)
		if len(instructions) != len(tt.expected) {
			t.Errorf(WRONG_LENGTH, "instructions", len(tt.expected), len(instructions))
		}
		for i := range tt.expected {
			if instructions[i] != tt.expected[i] {
				t.Errorf(NOT_EXPECTED, "instructions[i]", tt.expected[i], instructions[i])
			}
		}
	}
}

func TestInsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpAdd),
	}
	expected := `0000 OpConstant 1
0003 OpConstant 2
0006 OpConstant 65535
0009 OpAdd
`
	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}
	if concatted.String() != expected {
		t.Errorf(NOT_EXPECTED, "concatted.String()", expected, concatted.String())
	}
}
