package vm

import (
	"errors"
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		program  []byte
		expected [32]int64
	}{
		{
			name: "nop",
			program: []byte{
				byte(OpcodeNop), 0, 0, 0,
			},
			expected: [32]int64{},
		},
		{
			name: "load immediate",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 1,
			},
			expected: [32]int64{1},
		},
		{
			name: "push",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 1,
				byte(OpcodePush), 0, 0, 0,
			},
			expected: [32]int64{1},
		},
		{
			name: "pop",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 1,
				byte(OpcodePush), 0, 0, 0,
				byte(OpcodePop), 1, 0, 0,
			},
			expected: [32]int64{1, 1},
		},
		{
			name: "add",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 2,
				byte(OpcodeAdd), 2, 0, 1,
			},
			expected: [32]int64{1, 2, 3},
		},
		{
			name: "sub",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 2,
				byte(OpcodeSub), 2, 0, 1,
			},
			expected: [32]int64{1, 2, -1},
		},
		{
			name: "mul",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 2,
				byte(OpcodeLoadImmediate), 1, 0, 2,
				byte(OpcodeMul), 2, 0, 1,
			},
			expected: [32]int64{2, 2, 4},
		},
		{
			name: "div",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 4,
				byte(OpcodeLoadImmediate), 1, 0, 2,
				byte(OpcodeDiv), 2, 0, 1,
			},
			expected: [32]int64{4, 2, 2},
		},
		{
			name: "mod",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 5,
				byte(OpcodeLoadImmediate), 1, 0, 2,
				byte(OpcodeMod), 2, 0, 1,
			},
			expected: [32]int64{5, 2, 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			vm := New(test.program)
			err := vm.Run()

			if err != nil {
				t.Fatalf("expected no error, got %s", err.Error())
			}

			if !reflect.DeepEqual(vm.registers, test.expected) {
				t.Fatalf(
					"expected registers to be %v, got %v",
					test.expected,
					vm.registers,
				)
			}
		})
	}
}

func TestRunErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		program  []byte
		expected error
	}{
		{
			name: "division by zero",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 4,
				byte(OpcodeLoadImmediate), 1, 0, 0,
				byte(OpcodeDiv), 2, 0, 1,
			},
			expected: errors.New("division by zero"),
		},
		{
			name: "modulo by zero",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 5,
				byte(OpcodeLoadImmediate), 1, 0, 0,
				byte(OpcodeMod), 2, 0, 1,
			},
			expected: errors.New("modulo by zero"),
		},
		{
			name: "stack overflow",
			program: func() []byte {
				program := make([]byte, 0, StackSize*4)

				for range StackSize + 1 {
					program = append(program, byte(OpcodePush), 0, 0, 0)
				}

				return program
			}(),
			expected: errors.New("stack overflow"),
		},
		{
			name: "stack underflow",
			program: []byte{
				byte(OpcodePop), 0, 0, 0,
			},
			expected: errors.New("stack underflow"),
		},
		{
			name: "unknown opcode",
			program: []byte{
				byte(255), 0, 0, 0,
			},
			expected: errors.New("unknown opcode: 11111111"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			vm := New(test.program)
			err := vm.Run()

			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if err.Error() != test.expected.Error() {
				t.Fatalf(
					"expected error to be \"%s\", got \"%s\"",
					test.expected.Error(),
					err.Error(),
				)
			}
		})
	}
}
