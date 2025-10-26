package vm

import (
	"errors"
)

func (v *VM) instructionLoadRegister(dest register, rawSrc1 register) error {
	if rawSrc1 >= register(len(v.registers)) || dest >= register(len(v.registers)) {
		return errors.New("register out of bounds")
	}

	v.registers[dest] = v.registers[rawSrc1]

	return nil
}
