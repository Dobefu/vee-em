package vm

import (
	"errors"
)

func (v *VM) instructionStoreMemory(
	instructionStart register,
	instructionEnd register,
) error {
	if instructionEnd > v.programLen {
		return errors.New("unexpected end of program")
	}

	srcReg := register(v.program[instructionStart+1]) & NumRegistersMask
	addrReg := register(v.program[instructionStart+2]) & NumRegistersMask

	addr := v.registers[addrReg]

	if addr < 0 || uint64(addr) >= HeapSize {
		return errors.New("memory address out of bounds")
	}

	v.heap[addr] = v.registers[srcReg]

	return nil
}
