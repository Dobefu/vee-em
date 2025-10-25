package vm

import (
	"errors"
	"fmt"
	"os"
)

// Run runs the VM.
func (v *VM) Run() error {
	for int(v.pc) < len(v.program) {
		opcode, dest, src1, src2 := v.decodeInstruction()

		switch opcode {
		case OpcodeNop:
			// noop

		case OpcodePush:
			if int(v.sp) >= len(v.stack) {
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

		case OpcodeAdd:
			v.registers[dest] = v.registers[src1] + v.registers[src2]

		case OpcodeSub:
			v.registers[dest] = v.registers[src1] - v.registers[src2]

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
