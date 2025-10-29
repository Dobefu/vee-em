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

	for v.pc < uint64(len(v.program)) {
		var instructionErr error
		shouldIncrementPC := true

		opcode, rawDest, rawSrc1, rawSrc2, err := v.decodeInstruction()

		if err != nil {
			instructionErr = err
			opcode = OpcodeNop
		}

		switch opcode {
		case OpcodeNop:
			// noop

		case OpcodePush:
			instructionErr = v.instructionPush(rawSrc1)

		case OpcodePop:
			instructionErr = v.instructionPop(rawDest)

		case OpcodeLoadImmediate:
			instructionErr = v.instructionLoadImmediate(rawDest)

		case OpcodeLoadRegister:
			instructionErr = v.instructionLoadRegister(rawDest, rawSrc1)

		case OpcodeAdd:
			instructionErr = v.instructionAdd(rawDest, rawSrc1, rawSrc2)

		case OpcodeSub:
			instructionErr = v.instructionSub(rawDest, rawSrc1, rawSrc2)

		case OpcodeMul:
			instructionErr = v.instructionMul(rawDest, rawSrc1, rawSrc2)

		case OpcodeDiv:
			instructionErr = v.instructionDiv(rawDest, rawSrc1, rawSrc2)

		case OpcodeMod:
			instructionErr = v.instructionMod(rawDest, rawSrc1, rawSrc2)

		case OpcodeAND:
			instructionErr = v.instructionAND(rawDest, rawSrc1, rawSrc2)

		case OpcodeOR:
			instructionErr = v.instructionOR(rawDest, rawSrc1, rawSrc2)

		case OpcodeXOR:
			instructionErr = v.instructionXOR(rawDest, rawSrc1, rawSrc2)

		case OpcodeJmpImmediate:
			instructionErr = v.instructionJmpImmediate()
			shouldIncrementPC = false

		case OpcodeJmpRegister:
			instructionErr = v.instructionJmpRegister(rawSrc1)
			shouldIncrementPC = false

		default:
			instructionErr = fmt.Errorf("unknown opcode: %08b", opcode)
		}

		if instructionErr != nil {
			return instructionErr
		}

		if shouldIncrementPC {
			v.incrementPC()
		}
	}

	return nil
}
