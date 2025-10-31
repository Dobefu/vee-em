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
		expected [NumRegisters]int64
	}{
		{
			name: "nop",
			program: []byte{
				byte(OpcodeNop),
			},
			expected: [NumRegisters]int64{},
		},
		{
			name: "load immediate",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
			},
			expected: [NumRegisters]int64{1},
		},
		{
			name: "load register",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadRegister), 1, 0,
			},
			expected: [NumRegisters]int64{1, 1},
		},
		{
			name: "push",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodePush), 0,
			},
			expected: [NumRegisters]int64{1},
		},
		{
			name: "pop",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodePush), 0,
				byte(OpcodePop), 1,
			},
			expected: [NumRegisters]int64{1, 1},
		},
		{
			name: "add",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeAdd), 2, 0, 1,
			},
			expected: [NumRegisters]int64{1, 2, 3},
		},
		{
			name: "sub",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeSub), 2, 0, 1,
			},
			expected: [NumRegisters]int64{1, 2, -1},
		},
		{
			name: "mul",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeMul), 2, 0, 1,
			},
			expected: [NumRegisters]int64{2, 2, 4},
		},
		{
			name: "div",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 4,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeDiv), 2, 0, 1,
			},
			expected: [NumRegisters]int64{4, 2, 2},
		},
		{
			name: "mod",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 5,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeMod), 2, 0, 1,
			},
			expected: [NumRegisters]int64{5, 2, 1},
		},
		{
			name: "and",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 0b10000011,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 0b11000001,
				byte(OpcodeAND), 2, 0, 1,
			},
			expected: [NumRegisters]int64{0b10000011, 0b11000001, 0b10000001},
		},
		{
			name: "or",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 0b10000011,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 0b11000001,
				byte(OpcodeOR), 2, 0, 1,
			},
			expected: [NumRegisters]int64{0b10000011, 0b11000001, 0b11000011},
		},
		{
			name: "xor",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 0b00011111,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 0b11111000,
				byte(OpcodeXOR), 2, 0, 1,
			},
			expected: [NumRegisters]int64{0b00011111, 0b11111000, 0b11100111},
		},
		{
			name: "jmp immediate",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeJmpImmediate), 0, 0, 0, 0, 0, 0, 0, 33,
				byte(OpcodeAdd), 0, 0, 0, // This should get skipped.
				byte(OpcodeAdd), 0, 0, 0,
			},
			expected: [NumRegisters]int64{2, 1},
		},
		{
			name: "jmp register",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 36,
				byte(OpcodeJmpRegister), 2,
				byte(OpcodeAdd), 0, 0, 0, // This should get skipped.
				byte(OpcodeAdd), 0, 0, 0,
			},
			expected: [NumRegisters]int64{2, 1, 36},
		},
		{
			name: "jmp if zero",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 0,
				byte(OpcodeJmpImmediateIfZero), 0, 0, 0, 0, 0, 0, 0, 0, 20,
				byte(OpcodeAdd), 0, 0, 0, // This should get skipped.
			},
			expected: [NumRegisters]int64{0},
		},
		{
			name: "jmp if zero (not zero)",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeJmpImmediateIfZero), 0, 0, 0, 0, 0, 0, 0, 0, 20,
				byte(OpcodeAdd), 0, 0, 0, // This should not get skipped.
			},
			expected: [NumRegisters]int64{2},
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
			name:     "missing magic header",
			program:  []byte{},
			expected: errors.New("invalid magic header"),
		},
		{
			name: "invalid magic header",
			program: []byte{
				0xFF,
				byte(OpcodeNop),
			},
			expected: errors.New("invalid magic header"),
		},
		{
			name: "unexpected end of program",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate),
			},
			expected: errors.New("unexpected end of program"),
		},
		{
			name: "opcode push too few arguments",
			program: []byte{
				0x00,
				byte(OpcodePush),
			},
			expected: errors.New("unexpected end of program"),
		},
		{
			name: "opcode pop too few arguments",
			program: []byte{
				0x00,
				byte(OpcodePop),
			},
			expected: errors.New("unexpected end of program"),
		},
		{
			name: "opcode load register too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeLoadRegister),
			},
			expected: errors.New("unexpected end of program"),
		},
		{
			name: "opcode add too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeAdd),
			},
			expected: errors.New("unexpected end of program"),
		},
		{
			name: "opcode sub too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeSub),
			},
			expected: errors.New("unexpected end of program"),
		},
		{
			name: "opcode mul too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeMul),
			},
			expected: errors.New("unexpected end of program"),
		},
		{
			name: "opcode div too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeDiv),
			},
			expected: errors.New("unexpected end of program"),
		},
		{
			name: "opcode mod too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeMod),
			},
			expected: errors.New("unexpected end of program"),
		},
		{
			name: "opcode AND too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeAND),
			},
			expected: errors.New("unexpected end of program"),
		},
		{
			name: "opcode OR too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeOR),
			},
			expected: errors.New("unexpected end of program"),
		},
		{
			name: "opcode XOR too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeXOR),
			},
			expected: errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp immediate too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpImmediate),
			},
			expected: errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp register too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpRegister),
			},
			expected: errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp if zero too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpImmediateIfZero),
			},
			expected: errors.New("unexpected end of program"),
		},
		{
			name: "division by zero",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 4,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 0,
				byte(OpcodeDiv), 2, 0, 1,
			},
			expected: errors.New("division by zero"),
		},
		{
			name: "modulo by zero",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 5,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 0,
				byte(OpcodeMod), 2, 0, 1,
			},
			expected: errors.New("modulo by zero"),
		},
		{
			name: "stack overflow",
			program: func() []byte {
				program := make([]byte, 0, StackSize*2)
				program = append(program, 0x00)

				for range StackSize + 1 {
					program = append(program, byte(OpcodePush), 0)
				}

				return program
			}(),
			expected: errors.New("stack overflow"),
		},
		{
			name: "stack underflow",
			program: []byte{
				0x00,
				byte(OpcodePop), 0,
			},
			expected: errors.New("stack underflow"),
		},
		{
			name: "unknown opcode",
			program: []byte{
				0x00,
				byte(255),
			},
			expected: errors.New("unknown opcode: 11111111"),
		},
		{
			name: "jmp immediate target out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeJmpImmediate), 0, 0, 0, 0, 0, 0, 39, 15,
			},
			expected: errors.New("memory address out of bounds"),
		},
		{
			name: "jmp register memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 39, 15,
				byte(OpcodeJmpRegister), 0,
			},
			expected: errors.New("memory address out of bounds"),
		},
		{
			name: "jmp if zero memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 0,
				byte(OpcodeJmpImmediateIfZero), 0, 0, 0, 0, 0, 0, 0, 0, 21,
			},
			expected: errors.New("memory address out of bounds"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			vm := New(test.program, WithMagicHeader([]byte{0x00}))
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
