package vm

import (
	"errors"
)

func (v *VM) instructionDiv(
	instructionStart register,
	instructionEnd register,
) error {
	if instructionEnd > v.programLen {
		return errors.New("unexpected end of program")
	}

	dest := register(v.program[instructionStart+1]) & NumRegistersMask
	src1 := register(v.program[instructionStart+2]) & NumRegistersMask
	src2 := register(v.program[instructionStart+3]) & NumRegistersMask

	if v.registers[src2] == 0 {
		return errors.New("division by zero")
	}

	v.registers[dest] = v.registers[src1] / v.registers[src2]

	v.flags.isZero = v.registers[dest] == 0
	v.flags.isNegative = v.registers[dest] < 0

	return nil
}
