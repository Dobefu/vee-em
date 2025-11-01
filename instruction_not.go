package vm

import (
	"errors"
)

func (v *VM) instructionNOT(
	instructionStart register,
	instructionEnd register,
) error {
	if instructionEnd > v.programLen {
		return errors.New("unexpected end of program")
	}

	dest := register(v.program[instructionStart+1]) & NumRegistersMask
	src := register(v.program[instructionStart+2]) & NumRegistersMask

	v.registers[dest] = ^v.registers[src]

	v.flags.isZero = v.registers[dest] == 0
	v.flags.isNegative = v.registers[dest] < 0

	return nil
}
