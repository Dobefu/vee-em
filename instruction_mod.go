package vm

import (
	"errors"
)

func (v *VM) instructionMod(
	instructionStart register,
	instructionLen register,
) error {
	if instructionStart+instructionLen-1 >= register(len(v.program)) {
		return errors.New("unexpected end of program")
	}

	dest := register(v.program[instructionStart+1]) & NumRegistersMask
	src1 := register(v.program[instructionStart+2]) & NumRegistersMask
	src2 := register(v.program[instructionStart+3]) & NumRegistersMask

	if v.registers[src2] == 0 {
		return errors.New("modulo by zero")
	}

	v.registers[dest] = v.registers[src1] % v.registers[src2]

	return nil
}
