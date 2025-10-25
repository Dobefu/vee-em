package vm

import (
	"errors"
	"fmt"
	"os"
)

// Run runs the VM.
func (v *VM) Run() error {
	for v.pc < uint64(len(v.program)) {
		opcode, dest, src1, src2 := v.decodeInstruction()

		switch opcode {
		case OpcodeNop:
			// noop

		case OpcodePush:
			if v.sp >= uint64(len(v.stack)) {
				return errors.New("stack overflow")
			}

			v.stack[v.sp] = v.registers[src1]
			v.sp++

		case OpcodePop:
			if v.sp == 0 {
				return errors.New("stack underflow")
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
			if v.registers[src2] == 0 {
				return errors.New("division by zero")
			}

			v.registers[dest] = v.registers[src1] / v.registers[src2]

		case OpcodeMod:
			if v.registers[src2] == 0 {
				return errors.New("modulo by zero")
			}

			v.registers[dest] = v.registers[src1] % v.registers[src2]

		default:
			return fmt.Errorf("unknown opcode: %08b", opcode)
		}

		_, _ = fmt.Fprintf(os.Stdout, "Registers before: %v\n", v.registers)

		_, _ = fmt.Fprintf(os.Stdout, "Opcode: %08b (%d)\n", opcode, opcode)
		_, _ = fmt.Fprintf(os.Stdout, "Dest:   %08b (%d)\n", dest, dest)
		_, _ = fmt.Fprintf(os.Stdout, "Src1:   %08b (%d)\n", src1, src1)
		_, _ = fmt.Fprintf(os.Stdout, "Src2:   %08b (%d)\n", src2, src2)

		_, _ = fmt.Fprintf(os.Stdout, "Registers after:  %v\n", v.registers)

		v.incrementPC()
	}

	return nil
}
