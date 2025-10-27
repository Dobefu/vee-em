package vm

import (
	"errors"
)

func (v *VM) instructionJmpRegister(rawSrc1 register) error {
	if rawSrc1 >= NumRegisters {
		return errors.New("register out of bounds")
	}

	val := v.registers[rawSrc1]

	if val < 0 || val >= int64(len(v.program)) {
		return errors.New("memory address out of bounds")
	}

	v.pc = register(val)

	return nil
}
