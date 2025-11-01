package vm

import (
	"errors"
)

func (v *VM) instructionShiftRight(
	instructionStart register,
	instructionEnd register,
) error {
	if instructionEnd > v.programLen {
		return errors.New("unexpected end of program")
	}

	dest := register(v.program[instructionStart+1]) & NumRegistersMask
	src := register(v.program[instructionStart+2]) & NumRegistersMask
	shiftAmount := register(v.program[instructionStart+3]) & NumRegistersMask

	v.registers[dest] = int64(uint64(v.registers[src]) >> v.registers[shiftAmount]) // #nosec: G115

	v.flags.isZero = v.registers[dest] == 0
	v.flags.isNegative = v.registers[dest] < 0

	return nil
}
