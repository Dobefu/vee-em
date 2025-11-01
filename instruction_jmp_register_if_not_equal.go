package vm

import (
	"errors"
)

func (v *VM) instructionJmpRegisterIfNotEqual(
	instructionStart register,
	instructionEnd register,
) error {
	if instructionEnd > v.programLen {
		return errors.New("unexpected end of program")
	}

	addrReg := register(v.program[instructionStart+1]) & NumRegistersMask
	addr := v.registers[addrReg]

	if addr < 0 || uint64(addr) >= v.programLen {
		return errors.New("memory address out of bounds")
	}

	if !v.flags.isZero {
		v.pc = register(addr)
	}

	return nil
}
