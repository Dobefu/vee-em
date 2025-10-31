package vm

import (
	"encoding/binary"
	"errors"
)

func (v *VM) instructionJmpImmediateIfNotZero(
	instructionStart register,
	instructionEnd register,
) error {
	if instructionEnd > v.programLen {
		return errors.New("unexpected end of program")
	}

	src1 := register(v.program[instructionStart+1]) & NumRegistersMask

	addr := binary.BigEndian.Uint64(
		v.program[instructionStart+2 : instructionEnd],
	)

	if addr >= v.programLen {
		return errors.New("memory address out of bounds")
	}

	if v.registers[src1] != 0 {
		v.pc = addr
	}

	return nil
}
