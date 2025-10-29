package vm

import (
	"errors"
)

func (v *VM) instructionJmpImmediate() error {
	highByte := int64(v.program[v.pc+2]) << 8
	lowByte := int64(v.program[v.pc+3])
	addr := highByte | lowByte

	if addr < 0 || addr+3 >= int64(len(v.program)) {
		return errors.New("memory address out of bounds")
	}

	v.pc = register(addr)

	return nil
}
