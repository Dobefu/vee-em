package vm

import (
	"fmt"
)

// Run runs the VM.
func (v *VM) Run() error {
	err := v.validateMagicHeader()

	if err != nil {
		return err
	}

	for v.pc < v.programLen {
		var instructionErr error

		opcode := v.decodeInstruction()

		instructionLen := v.getInstructionLen(opcode)
		instructionStart := v.pc
		v.pc += instructionLen

		switch opcode {
		case OpcodeNop:
			// noop

		case OpcodePush:
			instructionErr = v.instructionPush(instructionStart, instructionLen)

		case OpcodePop:
			instructionErr = v.instructionPop(instructionStart, instructionLen)

		case OpcodeLoadImmediate:
			instructionErr = v.instructionLoadImmediate(instructionStart, instructionLen)

		case OpcodeLoadRegister:
			instructionErr = v.instructionLoadRegister(instructionStart, instructionLen)

		case OpcodeAdd:
			instructionErr = v.instructionAdd(instructionStart, instructionLen)

		case OpcodeSub:
			instructionErr = v.instructionSub(instructionStart, instructionLen)

		case OpcodeMul:
			instructionErr = v.instructionMul(instructionStart, instructionLen)

		case OpcodeDiv:
			instructionErr = v.instructionDiv(instructionStart, instructionLen)

		case OpcodeMod:
			instructionErr = v.instructionMod(instructionStart, instructionLen)

		case OpcodeAND:
			instructionErr = v.instructionAND(instructionStart, instructionLen)

		case OpcodeOR:
			instructionErr = v.instructionOR(instructionStart, instructionLen)

		case OpcodeXOR:
			instructionErr = v.instructionXOR(instructionStart, instructionLen)

		case OpcodeJmpImmediate:
			instructionErr = v.instructionJmpImmediate(instructionStart, instructionLen)

		case OpcodeJmpRegister:
			instructionErr = v.instructionJmpRegister(instructionStart, instructionLen)

		default:
			instructionErr = fmt.Errorf("unknown opcode: %08b", opcode)
		}

		if instructionErr != nil {
			return instructionErr
		}
	}

	return nil
}
