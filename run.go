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

		opcode, dest, rawSrc1, rawSrc2, err := v.decodeInstruction()

		src1 := rawSrc1 & 0x1F
		src2 := rawSrc2 & 0x1F

		if err != nil {
			instructionErr = err
			opcode = OpcodeNop
		}

		switch opcode {
		case OpcodeNop:
			// noop

		case OpcodePush:
			instructionErr = v.instructionPush(src1)

		case OpcodePop:
			instructionErr = v.instructionPop(dest)

		case OpcodeLoadImmediate:
			instructionErr = v.instructionLoadImmediate(dest)

		case OpcodeLoadRegister:
			instructionErr = v.instructionLoadRegister(dest, rawSrc1)

		case OpcodeAdd:
			instructionErr = v.instructionAdd(dest, src1, src2)

		case OpcodeSub:
			instructionErr = v.instructionSub(dest, src1, src2)

		case OpcodeMul:
			instructionErr = v.instructionMul(dest, src1, src2)

		case OpcodeDiv:
			instructionErr = v.instructionDiv(dest, src1, src2)

		case OpcodeMod:
			instructionErr = v.instructionMod(dest, src1, src2)

		case OpcodeXOR:
			instructionErr = v.instructionXOR(dest, src1, src2)

		case OpcodeJmpImmediate:
			instructionErr = v.instructionJmpImmediate(rawSrc1)
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
