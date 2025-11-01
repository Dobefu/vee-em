package vm

import (
	"errors"
)

func (v *VM) instructionCallRegister(
	instructionStart register,
	instructionEnd register,
) error {
	if instructionEnd > v.programLen {
		return errors.New("unexpected end of program")
	}

	if v.sp >= uint64(len(v.stack)) {
		return errors.New("stack overflow")
	}

	returnAddr := int64(v.pc) // #nosec: G115
	v.stack[v.sp] = returnAddr
	v.sp++

	src1 := register(v.program[instructionStart+1]) & NumRegistersMask
	addr := v.registers[src1]

	if addr < 0 || uint64(addr) >= v.programLen {
		return errors.New("memory address out of bounds")
	}

	v.pc = register(addr)

	return nil
}
