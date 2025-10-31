package vm

import (
	"encoding/binary"
	"errors"
)

func (v *VM) instructionJmpImmediate(
	instructionStart register,
	instructionEnd register,
) error {
	if instructionEnd > v.programLen {
		return errors.New("unexpected end of program")
	}

	addr := binary.BigEndian.Uint64(
		v.program[instructionStart+1 : instructionEnd],
	)

	if addr >= v.programLen {
		return errors.New("memory address out of bounds")
	}

	v.pc = addr

	return nil
}
