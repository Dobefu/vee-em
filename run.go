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

		opcode, err := v.decodeInstruction()

		if err != nil {
			instructionErr = err
			opcode = OpcodeNop
		}

		switch opcode {
		case OpcodeNop:
			// noop

		case OpcodePush:
			instructionErr = v.instructionPush()

		case OpcodePop:
			instructionErr = v.instructionPop()

		case OpcodeLoadImmediate:
			instructionErr = v.instructionLoadImmediate()

		case OpcodeLoadRegister:
			instructionErr = v.instructionLoadRegister()

		case OpcodeAdd:
			instructionErr = v.instructionAdd()

		case OpcodeSub:
			instructionErr = v.instructionSub()

		case OpcodeMul:
			instructionErr = v.instructionMul()

		case OpcodeDiv:
			instructionErr = v.instructionDiv()

		case OpcodeMod:
			instructionErr = v.instructionMod()

		case OpcodeAND:
			instructionErr = v.instructionAND()

		case OpcodeOR:
			instructionErr = v.instructionOR()

		case OpcodeXOR:
			instructionErr = v.instructionXOR()

		case OpcodeJmpImmediate:
			instructionErr = v.instructionJmpImmediate()
			shouldIncrementPC = false

		case OpcodeJmpRegister:
			instructionErr = v.instructionJmpRegister()
			shouldIncrementPC = false

		default:
			instructionErr = fmt.Errorf("unknown opcode: %08b", opcode)
		}

		if instructionErr != nil {
			return instructionErr
		}

		if shouldIncrementPC {
			instructionErr = v.incrementPC(opcode)

			if instructionErr != nil {
				return instructionErr
			}
		}
	}

	return nil
}
