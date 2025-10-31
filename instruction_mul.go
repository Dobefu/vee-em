package vm

import (
	"errors"
)

func (v *VM) instructionMul(
	instructionStart register,
	instructionEnd register,
) error {
	if instructionEnd > v.programLen {
		return errors.New("unexpected end of program")
	}

	dest := register(v.program[instructionStart+1]) & NumRegistersMask
	src1 := register(v.program[instructionStart+2]) & NumRegistersMask
	src2 := register(v.program[instructionStart+3]) & NumRegistersMask

	v.registers[dest] = v.registers[src1] * v.registers[src2]

	return nil
}
