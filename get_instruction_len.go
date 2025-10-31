package vm

func (v *VM) getInstructionLen(opcode Opcode) register {
	// Note: The instruction length includes the opcode itself.
	switch opcode {
	case OpcodeNop:
		return 1

	case OpcodePush:
		return 2

	case OpcodePop:
		return 2

	case OpcodeLoadImmediate:
		return 10

	case OpcodeLoadRegister:
		return 3

	case OpcodeAdd:
		return 4

	case OpcodeSub:
		return 4

	case OpcodeMul:
		return 4

	case OpcodeDiv:
		return 4

	case OpcodeMod:
		return 4

	case OpcodeAND:
		return 4

	case OpcodeOR:
		return 4

	case OpcodeXOR:
		return 4

	case OpcodeJmpImmediate:
		return 9

	case OpcodeJmpImmediateIfZero:
		return 10

	case OpcodeJmpImmediateIfNotZero:
		return 10

	case OpcodeJmpRegister:
		return 2

	case OpcodeJmpRegisterIfZero:
		return 3

	case OpcodeJmpRegisterIfNotZero:
		return 3

	case OpcodeCMP:
		return 3

	default:
		return 0
	}
}
