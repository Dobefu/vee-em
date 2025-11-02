package vm

// Note: The instruction length includes the opcode itself.
var instructionLengths = map[Opcode]register{
	OpcodeNop:                          1,
	OpcodePush:                         2,
	OpcodePop:                          2,
	OpcodeLoadImmediate:                10,
	OpcodeLoadRegister:                 3,
	OpcodeLoadMemory:                   3,
	OpcodeStoreMemory:                  3,
	OpcodeAdd:                          4,
	OpcodeSub:                          4,
	OpcodeMul:                          4,
	OpcodeDiv:                          4,
	OpcodeMod:                          4,
	OpcodeAND:                          4,
	OpcodeOR:                           4,
	OpcodeXOR:                          4,
	OpcodeNOT:                          3,
	OpcodeShiftLeft:                    4,
	OpcodeShiftRight:                   4,
	OpcodeShiftRightArithmetic:         4,
	OpcodeCMP:                          3,
	OpcodeJmpImmediate:                 9,
	OpcodeJmpImmediateIfZero:           10,
	OpcodeJmpImmediateIfNotZero:        10,
	OpcodeJmpImmediateIfEqual:          9,
	OpcodeJmpImmediateIfNotEqual:       9,
	OpcodeJmpImmediateIfGreater:        9,
	OpcodeJmpImmediateIfGreaterOrEqual: 9,
	OpcodeJmpImmediateIfLess:           9,
	OpcodeJmpImmediateIfLessOrEqual:    9,
	OpcodeJmpRegister:                  2,
	OpcodeJmpRegisterIfZero:            3,
	OpcodeJmpRegisterIfNotZero:         3,
	OpcodeJmpRegisterIfEqual:           2,
	OpcodeJmpRegisterIfNotEqual:        2,
	OpcodeJmpRegisterIfGreater:         2,
	OpcodeJmpRegisterIfGreaterOrEqual:  2,
	OpcodeJmpRegisterIfLess:            2,
	OpcodeJmpRegisterIfLessOrEqual:     2,
	OpcodeCallImmediate:                9,
	OpcodeCallRegister:                 2,
	OpcodeReturn:                       1,
	OpcodeHostCall:                     11,
	OpcodeHalt:                         1,
}

func (v *VM) getInstructionLen(opcode Opcode) register {
	length, ok := instructionLengths[opcode]

	if !ok {
		return 0
	}

	return length
}
