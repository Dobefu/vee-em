package vm

import (
	"errors"
)

func (v *VM) instructionJmpImmediate(rawSrc1 register) error {
	if rawSrc1+3 >= register(len(v.program)) {
		return errors.New("memory address out of bounds")
	}

	v.pc = rawSrc1

	return nil
}
