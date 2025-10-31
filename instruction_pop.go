package vm

import (
	"errors"
)

func (v *VM) instructionPop(
	instructionStart register,
	instructionEnd register,
) error {
	if instructionEnd > v.programLen {
		return errors.New("unexpected end of program")
	}

	rawDest := register(v.program[instructionStart+1])
	dest := rawDest & NumRegistersMask

	if v.sp == 0 {
		return errors.New("stack underflow")
	}

	v.registers[dest] = v.stack[v.sp-1]
	v.sp--

	return nil
}
