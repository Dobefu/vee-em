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
			highByte := int64(v.program[v.pc+2]) << 8
			lowByte := int64(v.program[v.pc+3])
			v.registers[dest] = highByte | lowByte

		case OpcodeAdd:
			v.registers[dest] = v.registers[src1] + v.registers[src2]

		case OpcodeSub:
			v.registers[dest] = v.registers[src1] - v.registers[src2]

		case OpcodeMul:
			v.registers[dest] = v.registers[src1] * v.registers[src2]

		case OpcodeDiv:
			instructionErr = v.instructionDiv(dest, src1, src2)

		case OpcodeMod:
			instructionErr = v.instructionMod(dest, src1, src2)

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
