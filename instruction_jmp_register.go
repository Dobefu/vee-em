package vm

import (
	"errors"
)

func (v *VM) instructionJmpRegister(rawSrc1 register) error {
	src1 := rawSrc1 & NumRegistersMask

	val := v.registers[src1]

	if val < 0 || val+3 >= int64(len(v.program)) {
		return errors.New("memory address out of bounds")
	}

	v.pc = register(val)

	return nil
}
