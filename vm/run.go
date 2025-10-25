package vm

import (
	"fmt"
	"os"
)

// Run runs the VM.
func (v *VM) Run() {
	for int(v.pc) < len(v.program) {
		opcode, dest, src1, src2 := v.decodeInstruction()

		_, _ = fmt.Fprintf(os.Stdout, "Opcode: %08b\n", opcode)
		_, _ = fmt.Fprintf(os.Stdout, "Opcode: %08b\n", dest)
		_, _ = fmt.Fprintf(os.Stdout, "Opcode: %08b\n", src1)
		_, _ = fmt.Fprintf(os.Stdout, "Opcode: %08b\n", src2)

		v.incrementPC()
	}
}
