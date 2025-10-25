package vm

import (
	"errors"
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

		opcode, dest, src1, src2, err := v.decodeInstruction()

		if err != nil {
			instructionErr = err
			opcode = OpcodeNop
		}

		switch opcode {
		case OpcodeNop:
			// noop

		case OpcodePush:
			if v.sp >= uint64(len(v.stack)) {
				instructionErr = errors.New("stack overflow")

				break
			}

			v.stack[v.sp] = v.registers[src1]
			v.sp++

		case OpcodePop:
			if v.sp == 0 {
				instructionErr = errors.New("stack underflow")

				break
			}

			v.registers[dest] = v.stack[v.sp-1]
			v.sp--

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
			if v.registers[src2] == 0 {
				instructionErr = errors.New("modulo by zero")

				break
			}

			v.registers[dest] = v.registers[src1] % v.registers[src2]

		default:
			instructionErr = fmt.Errorf("unknown opcode: %08b", opcode)
		}

		if instructionErr != nil {
			return instructionErr
		}

		v.incrementPC()
	}

	return nil
}
