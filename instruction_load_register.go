package vm

import (
	"errors"
)

func (v *VM) instructionLoadRegister(dest register, rawSrc1 register) error {
	if rawSrc1 >= NumRegisters || dest >= NumRegisters {
		return errors.New("register out of bounds")
	}

	v.registers[dest] = v.registers[rawSrc1]

	return nil
}
