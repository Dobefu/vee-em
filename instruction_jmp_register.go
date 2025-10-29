package vm

import (
	"errors"
)

func (v *VM) instructionJmpRegister() error {
	if v.pc+1 >= register(len(v.program)) {
		return errors.New("unexpected end of program")
	}

	rawSrc1 := register(v.program[v.pc+1])
	src1 := rawSrc1 & NumRegistersMask

	addr := v.registers[src1]

	if addr < 0 || addr >= int64(len(v.program)) {
		return errors.New("memory address out of bounds")
	}

	v.pc = register(addr)

	return nil
}
