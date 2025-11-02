package vm

import (
	"errors"
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		program           []byte
		hostCallHandler   HostCallHandler
		expectedRegisters [NumRegisters]int64
		expectedFlags     flags
	}{
		{
			name: "nop",
			program: []byte{
				byte(OpcodeNop),
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "load immediate",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{1},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "load register",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadRegister), 1, 0,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{1, 1},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "load memory",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 123,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 10,
				byte(OpcodeStoreMemory), 0, 1,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 0,
				byte(OpcodeLoadMemory), 0, 1,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{123, 10},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "store memory",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 42,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 0,
				byte(OpcodeStoreMemory), 0, 1,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{42, 0},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "push",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodePush), 0,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{1},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "pop",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodePush), 0,
				byte(OpcodePop), 1,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{1, 1},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "add",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeAdd), 2, 0, 1,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{1, 2, 3},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "sub",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeSub), 2, 0, 1,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{1, 2, -1},
			expectedFlags: flags{
				isZero:     false,
				isNegative: true,
			},
		},
		{
			name: "mul",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeMul), 2, 0, 1,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 2, 4},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "div",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 4,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeDiv), 2, 0, 1,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{4, 2, 2},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "mod",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 5,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeMod), 2, 0, 1,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{5, 2, 1},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "and",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 0b10000011,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 0b11000001,
				byte(OpcodeAND), 2, 0, 1,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{0b10000011, 0b11000001, 0b10000001},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "or",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 0b10000011,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 0b11000001,
				byte(OpcodeOR), 2, 0, 1,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{0b10000011, 0b11000001, 0b11000011},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "xor",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 0b00011111,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 0b11111000,
				byte(OpcodeXOR), 2, 0, 1,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{0b00011111, 0b11111000, 0b11100111},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "not",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 0b00001111,
				byte(OpcodeNOT), 1, 0,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{0b00001111, ^int64(0b00001111)},
			expectedFlags: flags{
				isZero:     false,
				isNegative: true,
			},
		},
		{
			name: "shift left",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 3,
				byte(OpcodeShiftLeft), 2, 0, 1,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{1, 3, 8},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "shift right",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 16,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeShiftRight), 2, 0, 1,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{16, 2, 4},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "shift right arithmetic positive",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 16,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeShiftRightArithmetic), 2, 0, 1,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{16, 2, 4},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "shift right arithmetic negative",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 255, 255, 255, 255, 255, 255, 255, 240,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeShiftRightArithmetic), 2, 0, 1,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{-16, 2, -4},
			expectedFlags: flags{
				isZero:     false,
				isNegative: true,
			},
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
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 1},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
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
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 1, 36},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "jmp immediate if zero",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 0,
				byte(OpcodeJmpImmediateIfZero), 0, 0, 0, 0, 0, 0, 0, 0, 20,
				byte(OpcodeAdd), 0, 0, 0, // This should get skipped.
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{0},
			expectedFlags: flags{
				isZero:     true,
				isNegative: false,
			},
		},
		{
			name: "jmp immediate if zero (not zero)",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeJmpImmediateIfZero), 0, 0, 0, 0, 0, 0, 0, 0, 20,
				byte(OpcodeAdd), 0, 0, 0, // This should not get skipped.
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "jmp immediate if not zero",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeJmpImmediateIfNotZero), 0, 0, 0, 0, 0, 0, 0, 0, 23,
				byte(OpcodeAdd), 0, 0, 0, // This should get skipped.
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{1},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "jmp immediate if not zero (zero)",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 0,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeJmpImmediateIfNotZero), 0, 0, 0, 0, 0, 0, 0, 0, 33,
				byte(OpcodeAdd), 0, 0, 1, // This should not get skipped.
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{1, 1},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "jmp immediate if equal",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpImmediateIfEqual), 0, 0, 0, 0, 0, 0, 0, 36,
				byte(OpcodeAdd), 0, 0, 0, // This should get skipped.
				byte(OpcodeAdd), 0, 0, 0,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 1},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "jmp immediate if equal (not equal)",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpImmediateIfEqual), 0, 0, 0, 0, 0, 0, 0, 36,
				byte(OpcodeAdd), 0, 0, 0, // This should not get skipped.
				byte(OpcodeAdd), 0, 0, 0,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{4, 2},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "jmp immediate if not equal",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpImmediateIfNotEqual), 0, 0, 0, 0, 0, 0, 0, 36,
				byte(OpcodeAdd), 0, 0, 0, // This should get skipped.
				byte(OpcodeAdd), 0, 0, 0,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 2},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "jmp immediate if not equal (equal)",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpImmediateIfNotEqual), 0, 0, 0, 0, 0, 0, 0, 36,
				byte(OpcodeAdd), 0, 0, 0, // This should not get skipped.
				byte(OpcodeAdd), 0, 0, 0,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{4, 1},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "jmp immediate if greater",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpImmediateIfGreater), 0, 0, 0, 0, 0, 0, 0, 42,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 2, // This should get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 1},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "jmp immediate if greater (not greater)",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpImmediateIfGreater), 0, 0, 0, 0, 0, 0, 0, 42,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 2, // This should not get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 1, 2},
			expectedFlags: flags{
				isZero:     true,
				isNegative: false,
			},
		},
		{
			name: "jmp immediate if greater or equal",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpImmediateIfGreaterOrEqual), 0, 0, 0, 0, 0, 0, 0, 42,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 2, // This should get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 1},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "jmp immediate if greater or equal (not greater or equal)",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpImmediateIfGreaterOrEqual), 0, 0, 0, 0, 0, 0, 0, 42,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 2, // This should not get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 2, 2},
			expectedFlags: flags{
				isZero:     false,
				isNegative: true,
			},
		},
		{
			name: "jmp immediate if less",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpImmediateIfLess), 0, 0, 0, 0, 0, 0, 0, 42,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 2, // This should get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 2},
			expectedFlags: flags{
				isZero:     false,
				isNegative: true,
			},
		},
		{
			name: "jmp immediate if less (not less)",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpImmediateIfLess), 0, 0, 0, 0, 0, 0, 0, 42,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 2, // This should not get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 1, 2},
			expectedFlags: flags{
				isZero:     true,
				isNegative: false,
			},
		},
		{
			name: "jmp immediate if less or equal",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpImmediateIfLessOrEqual), 0, 0, 0, 0, 0, 0, 0, 42,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 2, // This should get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 2},
			expectedFlags: flags{
				isZero:     false,
				isNegative: true,
			},
		},
		{
			name: "jmp immediate if less or equal (not less or equal)",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpImmediateIfLessOrEqual), 0, 0, 0, 0, 0, 0, 0, 42,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 2, // This should not get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 1, 2},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "jmp register if zero",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 0,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 23,
				byte(OpcodeJmpRegisterIfZero), 0, 1,
				byte(OpcodeAdd), 0, 0, 0, // This should get skipped.
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{0, 23},
			expectedFlags: flags{
				isZero:     true,
				isNegative: false,
			},
		},
		{
			name: "jmp register if zero (not zero)",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 23,
				byte(OpcodeJmpRegisterIfZero), 0, 1,
				byte(OpcodeAdd), 0, 0, 0, // This should not get skipped.
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 23},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "jmp register if not zero",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 26,
				byte(OpcodeJmpRegisterIfNotZero), 0, 1,
				byte(OpcodeAdd), 0, 0, 0, // This should get skipped.
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{1, 26},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "jmp register if not zero (zero)",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 0,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 36,
				byte(OpcodeJmpRegisterIfNotZero), 0, 0,
				byte(OpcodeAdd), 0, 0, 1, // This should not get skipped.
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{1, 1, 36},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "jmp register if equal",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 45,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpRegisterIfEqual), 2,
				byte(OpcodeLoadImmediate), 3, 0, 0, 0, 0, 0, 0, 0, 2, // This should get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 1, 45},
			expectedFlags: flags{
				isZero:     true,
				isNegative: false,
			},
		},
		{
			name: "jmp register if equal (not equal)",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 45,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpRegisterIfEqual), 2,
				byte(OpcodeLoadImmediate), 3, 0, 0, 0, 0, 0, 0, 0, 2, // This should not get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 2, 45, 2},
			expectedFlags: flags{
				isZero:     false,
				isNegative: true,
			},
		},
		{
			name: "jmp register if not equal",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 45,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpRegisterIfNotEqual), 2,
				byte(OpcodeLoadImmediate), 3, 0, 0, 0, 0, 0, 0, 0, 2, // This should get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 2, 45},
			expectedFlags: flags{
				isZero:     false,
				isNegative: true,
			},
		},
		{
			name: "jmp register if not equal (equal)",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 45,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpRegisterIfNotEqual), 2,
				byte(OpcodeLoadImmediate), 3, 0, 0, 0, 0, 0, 0, 0, 2, // This should not get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 1, 45, 2},
			expectedFlags: flags{
				isZero:     true,
				isNegative: false,
			},
		},
		{
			name: "jmp register if greater",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 45,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpRegisterIfGreater), 2,
				byte(OpcodeLoadImmediate), 3, 0, 0, 0, 0, 0, 0, 0, 2, // This should get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 1, 45},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "jmp register if greater (not greater)",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 45,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpRegisterIfGreater), 2,
				byte(OpcodeLoadImmediate), 3, 0, 0, 0, 0, 0, 0, 0, 2, // This should not get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 1, 45, 2},
			expectedFlags: flags{
				isZero:     true,
				isNegative: false,
			},
		},
		{
			name: "jmp register if greater or equal",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 45,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpRegisterIfGreaterOrEqual), 2,
				byte(OpcodeLoadImmediate), 3, 0, 0, 0, 0, 0, 0, 0, 2, // This should get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 1, 45},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "jmp register if greater or equal (not greater or equal)",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 45,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpRegisterIfGreaterOrEqual), 2,
				byte(OpcodeLoadImmediate), 3, 0, 0, 0, 0, 0, 0, 0, 2, // This should not get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 2, 45, 2},
			expectedFlags: flags{
				isZero:     false,
				isNegative: true,
			},
		},
		{
			name: "jmp register if less",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 45,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpRegisterIfLess), 2,
				byte(OpcodeLoadImmediate), 3, 0, 0, 0, 0, 0, 0, 0, 2, // This should get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 2, 45},
			expectedFlags: flags{
				isZero:     false,
				isNegative: true,
			},
		},
		{
			name: "jmp register if less (not less)",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 45,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpRegisterIfLess), 2,
				byte(OpcodeLoadImmediate), 3, 0, 0, 0, 0, 0, 0, 0, 2, // This should not get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 1, 45, 2},
			expectedFlags: flags{
				isZero:     true,
				isNegative: false,
			},
		},
		{
			name: "jmp register if less or equal",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 45,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpRegisterIfLessOrEqual), 2,
				byte(OpcodeLoadImmediate), 3, 0, 0, 0, 0, 0, 0, 0, 2, // This should get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 2, 45},
			expectedFlags: flags{
				isZero:     false,
				isNegative: true,
			},
		},
		{
			name: "jmp register if less or equal (not less or equal)",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 45,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeJmpRegisterIfLessOrEqual), 2,
				byte(OpcodeLoadImmediate), 3, 0, 0, 0, 0, 0, 0, 0, 2, // This should not get skipped.
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{2, 1, 45, 2},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "cmp",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeCMP), 0, 1,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{1, 2},
			expectedFlags: flags{
				isZero:     false,
				isNegative: true,
			},
		},
		{
			name: "call immediate",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 42,
				byte(OpcodeCallImmediate), 0, 0, 0, 0, 0, 0, 0, 22,
				byte(OpcodeAdd), 0, 0, 0,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 100,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{42, 100},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "call register",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 42,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 25,
				byte(OpcodeCallRegister), 2,
				byte(OpcodeAdd), 0, 0, 0,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 100,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{42, 100, 25},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "return",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 42,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 33,
				byte(OpcodePush), 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 100,
				byte(OpcodeReturn),
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 200,
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{42, 200},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "host call",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 42,
				byte(OpcodeHostCall), 0, 0, 0, 0, 0, 0, 0, 0, 1, 0,
			},
			hostCallHandler: func(
				_ int64,
				_ register,
				_ register,
				_ [NumRegisters]int64,
			) (int64, error) {
				return 1, nil
			},
			expectedRegisters: [NumRegisters]int64{42, 1},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
		{
			name: "halt",
			program: []byte{
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 42,
				byte(OpcodeHalt),
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 100, // This should not get executed.
			},
			hostCallHandler:   nil,
			expectedRegisters: [NumRegisters]int64{42},
			expectedFlags: flags{
				isZero:     false,
				isNegative: false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			vm := New(test.program, WithHostCallHandler(test.hostCallHandler))
			err := vm.Run()

			if err != nil {
				t.Fatalf("expected no error, got %s", err.Error())
			}

			if !reflect.DeepEqual(vm.registers, test.expectedRegisters) {
				t.Fatalf(
					"expected registers to be %v, got %v",
					test.expectedRegisters,
					vm.registers,
				)
			}

			if vm.flags != test.expectedFlags {
				t.Fatalf(
					"expected flags to be %v, got %v",
					test.expectedFlags,
					vm.flags,
				)
			}
		})
	}
}

func TestRunErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		program         []byte
		hostCallHandler HostCallHandler
		expected        error
	}{
		{
			name:            "missing magic header",
			program:         []byte{},
			hostCallHandler: nil,
			expected:        errors.New("invalid magic header"),
		},
		{
			name: "invalid magic header",
			program: []byte{
				0xFF,
				byte(OpcodeNop),
			},
			hostCallHandler: nil,
			expected:        errors.New("invalid magic header"),
		},
		{
			name: "unexpected end of program",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode push too few arguments",
			program: []byte{
				0x00,
				byte(OpcodePush),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
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
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode load memory too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeLoadMemory),
			},
			expected: errors.New("unexpected end of program"),
		},
		{
			name: "opcode store memory too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeStoreMemory),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode add too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeAdd),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
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
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode div too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeDiv),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode mod too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeMod),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode AND too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeAND),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode OR too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeOR),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode XOR too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeXOR),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode NOT too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeNOT),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode shift left too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeShiftLeft),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode shift right too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeShiftRight),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode shift right arithmetic too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeShiftRightArithmetic),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode CMP too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeCMP),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
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
			name: "opcode jmp immediate if zero too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpImmediateIfZero),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp immediate if not zero too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpImmediateIfNotZero),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp immediate if equal too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpImmediateIfEqual),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp immediate if not equal too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpImmediateIfNotEqual),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp immediate if greater too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpImmediateIfGreater),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp immediate if greater or equal too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpImmediateIfGreaterOrEqual),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp immediate if less too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpImmediateIfLess),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp immediate if less or equal too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpImmediateIfLessOrEqual),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp register too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpRegister),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp register if zero too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpRegisterIfZero),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp register if not zero too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpRegisterIfNotZero),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp register if equal too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpRegisterIfEqual),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp register if not equal too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpRegisterIfNotEqual),
			},
			expected: errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp register if greater too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpRegisterIfGreater),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp register if greater or equal too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpRegisterIfGreaterOrEqual),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp register if less too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpRegisterIfLess),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode jmp register if less or equal too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeJmpRegisterIfLessOrEqual),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode call immediate too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeCallImmediate),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode call register too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeCallRegister),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "opcode host call too few arguments",
			program: []byte{
				0x00,
				byte(OpcodeHostCall),
			},
			hostCallHandler: nil,
			expected:        errors.New("unexpected end of program"),
		},
		{
			name: "division by zero",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 4,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 0,
				byte(OpcodeDiv), 2, 0, 1,
			},
			hostCallHandler: nil,
			expected:        errors.New("division by zero"),
		},
		{
			name: "modulo by zero",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 5,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 0,
				byte(OpcodeMod), 2, 0, 1,
			},
			hostCallHandler: nil,
			expected:        errors.New("modulo by zero"),
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
			hostCallHandler: nil,
			expected:        errors.New("stack overflow"),
		},
		{
			name: "stack underflow",
			program: []byte{
				0x00,
				byte(OpcodePop), 0,
			},
			hostCallHandler: nil,
			expected:        errors.New("stack underflow"),
		},
		{
			name: "unknown opcode",
			program: []byte{
				0x00,
				byte(255),
			},
			hostCallHandler: nil,
			expected:        errors.New("unknown opcode: 11111111"),
		},
		{
			name: "jmp immediate target out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeJmpImmediate), 0, 0, 0, 0, 0, 0, 39, 15,
			},
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "jmp immediate if zero memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 0,
				byte(OpcodeJmpImmediateIfZero), 0, 0, 0, 0, 0, 0, 0, 0, 21,
			},
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "jmp immediate if not zero memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeJmpImmediateIfNotZero), 0, 0, 0, 0, 0, 0, 0, 0, 21,
			},
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "jmp immediate if equal memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeJmpImmediateIfEqual), 0, 0, 0, 0, 0, 0, 0, 34,
			},
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "jmp immediate if not equal memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeJmpImmediateIfNotEqual), 0, 0, 0, 0, 0, 0, 0, 34,
			},
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "jmp immediate if greater memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeJmpImmediateIfGreater), 0, 0, 0, 0, 0, 0, 0, 34,
			},
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "jmp immediate if greater or equal memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeJmpImmediateIfGreaterOrEqual), 0, 0, 0, 0, 0, 0, 0, 34,
			},
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "jmp immediate if less memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeJmpImmediateIfLess), 0, 0, 0, 0, 0, 0, 0, 34,
			},
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "jmp immediate if less or equal memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeJmpImmediateIfLessOrEqual), 0, 0, 0, 0, 0, 0, 0, 34,
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
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "jmp register if zero memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 0,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 24,
				byte(OpcodeJmpRegisterIfZero), 0, 1,
			},
			expected: errors.New("memory address out of bounds"),
		},
		{
			name: "jmp register if not zero memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 24,
				byte(OpcodeJmpRegisterIfNotZero), 0, 1,
			},
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "jmp register if equal memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 36,
				byte(OpcodeJmpRegisterIfEqual), 2,
			},
			expected: errors.New("memory address out of bounds"),
		},
		{
			name: "jmp register if not equal memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 36,
				byte(OpcodeJmpRegisterIfNotEqual), 2,
			},
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "jmp register if greater memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 36,
				byte(OpcodeJmpRegisterIfGreater), 2,
			},
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "jmp register if greater or equal memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 36,
				byte(OpcodeJmpRegisterIfGreaterOrEqual), 2,
			},
			expected: errors.New("memory address out of bounds"),
		},
		{
			name: "jmp register if less memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 36,
				byte(OpcodeJmpRegisterIfLess), 2,
			},
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "jmp register if less or equal memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 1,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 0, 0, 2,
				byte(OpcodeCMP), 0, 1,
				byte(OpcodeLoadImmediate), 2, 0, 0, 0, 0, 0, 0, 0, 36,
				byte(OpcodeJmpRegisterIfLessOrEqual), 2,
			},
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "load memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 1, 0, 0,
				byte(OpcodeLoadMemory), 1, 0,
			},
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "store memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 42,
				byte(OpcodeLoadImmediate), 1, 0, 0, 0, 0, 0, 1, 0, 0,
				byte(OpcodeStoreMemory), 0, 1,
			},
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "call immediate stack overflow",
			program: func() []byte {
				program := make([]byte, 0, StackSize*2)
				program = append(program, 0x00)

				for range StackSize + 1 {
					program = append(program, byte(OpcodeCallImmediate), 0, 0, 0, 0, 0, 0, 0, 10)
				}

				return program
			}(),
			hostCallHandler: nil,
			expected:        errors.New("stack overflow"),
		},
		{
			name: "call immediate memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeCallImmediate), 0, 0, 0, 0, 0, 0, 0, 50,
			},
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "call register stack overflow",
			program: func() []byte {
				program := make([]byte, 0, StackSize*2)
				program = append(program, 0x00)
				targetAddr := uint64(12)
				program = append(program, byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, byte(targetAddr))
				program = append(program, byte(OpcodeNop))

				for range StackSize + 1 {
					program = append(program, byte(OpcodeCallRegister), 0)
				}

				return program
			}(),
			hostCallHandler: nil,
			expected:        errors.New("stack overflow"),
		},
		{
			name: "call register memory address out of bounds",
			program: []byte{
				0x00,
				byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 50,
				byte(OpcodeCallRegister), 0,
			},
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "return stack underflow",
			program: []byte{
				0x00,
				byte(OpcodeReturn),
			},
			hostCallHandler: nil,
			expected:        errors.New("stack underflow"),
		},
		{
			name: "return memory address out of bounds",
			program: func() []byte {
				program := make([]byte, 0, 20)
				program = append(program, 0x00)
				program = append(program, byte(OpcodeLoadImmediate), 0, 0, 0, 0, 0, 0, 0, 0, 50)
				program = append(program, byte(OpcodePush), 0)
				program = append(program, byte(OpcodeReturn))

				return program
			}(),
			hostCallHandler: nil,
			expected:        errors.New("memory address out of bounds"),
		},
		{
			name: "host call handler error",
			program: []byte{
				0x00,
				byte(OpcodeHostCall), 0, 0, 0, 0, 0, 0, 0, 0, 1, 0,
			},
			hostCallHandler: func(
				_ int64,
				_ register,
				_ register,
				_ [NumRegisters]int64,
			) (int64, error) {
				return 0, errors.New("host call handler error")
			},
			expected: errors.New("host call handler error"),
		},
		{
			name: "host call handler not set",
			program: []byte{
				0x00,
				byte(OpcodeHostCall), 0, 0, 0, 0, 0, 0, 0, 0, 1, 0,
			},
			hostCallHandler: nil,
			expected:        errors.New("host call handler not set"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			vm := New(
				test.program,
				WithMagicHeader([]byte{0x00}),
				WithHostCallHandler(test.hostCallHandler),
			)
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
