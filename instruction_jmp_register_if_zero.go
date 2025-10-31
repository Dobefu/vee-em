package vm

import (
	"errors"
)

func (v *VM) instructionJmpRegisterIfZero(
	instructionStart register,
	instructionEnd register,
) error {
	if instructionEnd > v.programLen {
		return errors.New("unexpected end of program")
	}

	checkReg := register(v.program[instructionStart+1]) & NumRegistersMask
	addrReg := register(v.program[instructionStart+2]) & NumRegistersMask

	if v.registers[checkReg] == 0 {
		addr := v.registers[addrReg]

		if addr < 0 || uint64(addr) >= v.programLen {
			return errors.New("memory address out of bounds")
		}

		v.pc = register(addr)
	}

	return nil
}
