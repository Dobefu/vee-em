package vm

import (
	"errors"
)

func (v *VM) instructionLoadMemory(
	instructionStart register,
	instructionEnd register,
) error {
	if instructionEnd > v.programLen {
		return errors.New("unexpected end of program")
	}

	dest := register(v.program[instructionStart+1]) & NumRegistersMask
	addrReg := register(v.program[instructionStart+2]) & NumRegistersMask

	addr := v.registers[addrReg]

	if addr < 0 || uint64(addr) >= HeapSize {
		return errors.New("memory address out of bounds")
	}

	v.registers[dest] = v.heap[addr]

	v.flags.isZero = v.registers[dest] == 0
	v.flags.isNegative = v.registers[dest] < 0

	return nil
}
