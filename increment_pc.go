package vm

import (
	"fmt"
)

func (v *VM) incrementPC(opcode Opcode) error {
	var instructionLen register

	// Note: The instruction length includes the opcode itself.
	switch opcode {
	case OpcodeNop:
		instructionLen = 1

	case OpcodePush:
		instructionLen = 2

	case OpcodePop:
		instructionLen = 2

	case OpcodeLoadImmediate:
		instructionLen = 10

	case OpcodeLoadRegister:
		instructionLen = 3

	case OpcodeAdd:
		instructionLen = 4

	case OpcodeSub:
		instructionLen = 4

	case OpcodeMul:
		instructionLen = 4

	case OpcodeDiv:
		instructionLen = 4

	case OpcodeMod:
		instructionLen = 4

	case OpcodeAND:
		instructionLen = 4

	case OpcodeOR:
		instructionLen = 4

	case OpcodeXOR:
		instructionLen = 4

	case OpcodeJmpImmediate:
		instructionLen = 9

	case OpcodeJmpRegister:
		instructionLen = 2

	default:
		return fmt.Errorf("unknown opcode: %08b", opcode)
	}

	v.pc += instructionLen

	return nil
}
