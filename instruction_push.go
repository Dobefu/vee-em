package vm

import (
	"errors"
)

func (v *VM) instructionPush(
	instructionStart register,
	instructionLen register,
) error {
	if instructionStart+instructionLen-1 >= register(len(v.program)) {
		return errors.New("unexpected end of program")
	}

	rawSrc1 := register(v.program[instructionStart+1])
	src1 := rawSrc1 & NumRegistersMask

	if v.sp >= uint64(len(v.stack)) {
		return errors.New("stack overflow")
	}

	v.stack[v.sp] = v.registers[src1]
	v.sp++

	return nil
}
