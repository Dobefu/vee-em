package vm

import (
	"errors"
)

func (v *VM) instructionMod(dest register, src1 register, src2 register) error {
	if v.registers[src2] == 0 {
		return errors.New("modulo by zero")
	}

	v.registers[dest] = v.registers[src1] % v.registers[src2]

	return nil
}
