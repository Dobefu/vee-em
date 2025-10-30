package vm

import (
	"testing"
)

func TestIncrementPC(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		opcode   Opcode
		expected register
	}{
		{
			name:     "jmp immediate",
			opcode:   OpcodeJmpImmediate,
			expected: 9,
		},
		{
			name:     "jmp register",
			opcode:   OpcodeJmpRegister,
			expected: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			vm := New([]byte{})
			vm.incrementPC(test.opcode)

			if vm.pc != test.expected {
				t.Fatalf("expected PC to be %d, got %d", test.expected, vm.pc)
			}
		})
	}
}
