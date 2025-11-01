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
		instructionEnd := instructionStart + instructionLen
		v.pc += instructionLen

		switch opcode {
		case OpcodeNop:
			// noop

		case OpcodePush:
			instructionErr = v.instructionPush(instructionStart, instructionEnd)

		case OpcodePop:
			instructionErr = v.instructionPop(instructionStart, instructionEnd)

		case OpcodeLoadImmediate:
			instructionErr = v.instructionLoadImmediate(instructionStart, instructionEnd)

		case OpcodeLoadRegister:
			instructionErr = v.instructionLoadRegister(instructionStart, instructionEnd)

		case OpcodeAdd:
			instructionErr = v.instructionAdd(instructionStart, instructionEnd)

		case OpcodeSub:
			instructionErr = v.instructionSub(instructionStart, instructionEnd)

		case OpcodeMul:
			instructionErr = v.instructionMul(instructionStart, instructionEnd)

		case OpcodeDiv:
			instructionErr = v.instructionDiv(instructionStart, instructionEnd)

		case OpcodeMod:
			instructionErr = v.instructionMod(instructionStart, instructionEnd)

		case OpcodeAND:
			instructionErr = v.instructionAND(instructionStart, instructionEnd)

		case OpcodeOR:
			instructionErr = v.instructionOR(instructionStart, instructionEnd)

		case OpcodeXOR:
			instructionErr = v.instructionXOR(instructionStart, instructionEnd)

		case OpcodeCMP:
			instructionErr = v.instructionCMP(instructionStart, instructionEnd)

		case OpcodeJmpImmediate:
			instructionErr = v.instructionJmpImmediate(instructionStart, instructionEnd)

		case OpcodeJmpImmediateIfZero:
			instructionErr = v.instructionJmpImmediateIfZero(instructionStart, instructionEnd)

		case OpcodeJmpImmediateIfNotZero:
			instructionErr = v.instructionJmpImmediateIfNotZero(instructionStart, instructionEnd)

		case OpcodeJmpImmediateIfEqual:
			instructionErr = v.instructionJmpImmediateIfEqual(instructionStart, instructionEnd)

		case OpcodeJmpImmediateIfNotEqual:
			instructionErr = v.instructionJmpImmediateIfNotEqual(instructionStart, instructionEnd)

		case OpcodeJmpImmediateIfGreater:
			instructionErr = v.instructionJmpImmediateIfGreater(instructionStart, instructionEnd)

		case OpcodeJmpImmediateIfGreaterOrEqual:
			instructionErr = v.instructionJmpImmediateIfGreaterOrEqual(instructionStart, instructionEnd)

		case OpcodeJmpImmediateIfLess:
			instructionErr = v.instructionJmpImmediateIfLess(instructionStart, instructionEnd)

		case OpcodeJmpImmediateIfLessOrEqual:
			instructionErr = v.instructionJmpImmediateIfLessOrEqual(instructionStart, instructionEnd)

		case OpcodeJmpRegister:
			instructionErr = v.instructionJmpRegister(instructionStart, instructionEnd)

		case OpcodeJmpRegisterIfZero:
			instructionErr = v.instructionJmpRegisterIfZero(instructionStart, instructionEnd)

		case OpcodeJmpRegisterIfNotZero:
			instructionErr = v.instructionJmpRegisterIfNotZero(instructionStart, instructionEnd)

		case OpcodeJmpRegisterIfEqual:
			instructionErr = v.instructionJmpRegisterIfEqual(instructionStart, instructionEnd)

		case OpcodeJmpRegisterIfNotEqual:
			instructionErr = v.instructionJmpRegisterIfNotEqual(instructionStart, instructionEnd)

		case OpcodeJmpRegisterIfGreater:
			instructionErr = v.instructionJmpRegisterIfGreater(instructionStart, instructionEnd)

		case OpcodeJmpRegisterIfGreaterOrEqual:
			instructionErr = v.instructionJmpRegisterIfGreaterOrEqual(instructionStart, instructionEnd)

		case OpcodeJmpRegisterIfLess:
			instructionErr = v.instructionJmpRegisterIfLess(instructionStart, instructionEnd)

		default:
			instructionErr = fmt.Errorf("unknown opcode: %08b", opcode)
		}

		if instructionErr != nil {
			return instructionErr
		}
	}

	return nil
}
