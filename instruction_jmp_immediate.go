package vm

import (
	"encoding/binary"
	"errors"
)

func (v *VM) instructionJmpImmediate(
	instructionStart register,
	instructionLen register,
) error {
	if instructionStart+instructionLen-1 >= v.programLen {
		return errors.New("unexpected end of program")
	}

	addr := binary.BigEndian.Uint64(
		v.program[instructionStart+1 : instructionStart+instructionLen],
	)

	if addr >= v.programLen {
		return errors.New("memory address out of bounds")
	}

	v.pc = addr

	return nil
}
