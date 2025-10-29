package vm

import (
	"errors"
)

func (v *VM) instructionDiv(
	rawDest register,
	rawSrc1 register,
	rawSrc2 register,
) error {
	dest := rawDest & NumRegistersMask
	src1 := rawSrc1 & NumRegistersMask
	src2 := rawSrc2 & NumRegistersMask

	if v.registers[src2] == 0 {
		return errors.New("division by zero")
	}

	v.registers[dest] = v.registers[src1] / v.registers[src2]

	return nil
}
