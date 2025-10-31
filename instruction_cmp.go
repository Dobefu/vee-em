package vm

import (
	"errors"
)

func (v *VM) instructionCMP(
	instructionStart register,
	instructionEnd register,
) error {
	if instructionEnd > v.programLen {
		return errors.New("unexpected end of program")
	}

	src1 := register(v.program[instructionStart+1]) & NumRegistersMask
	src2 := register(v.program[instructionStart+2]) & NumRegistersMask

	result := v.registers[src1] - v.registers[src2]

	v.flags.isZero = result == 0
	v.flags.isNegative = result < 0

	return nil
}
