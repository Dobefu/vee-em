package vm

import (
	"encoding/binary"
	"errors"
)

func (v *VM) instructionJmpImmediate(
	instructionStart register,
	instructionLen register,
) error {
	if instructionStart+instructionLen-1 >= register(len(v.program)) {
		return errors.New("unexpected end of program")
	}

	val := int64(binary.BigEndian.Uint64(
		v.program[instructionStart+1 : instructionStart+instructionLen],
	)) // #nosec: G115

	if val < 0 || val >= int64(len(v.program)) {
		return errors.New("memory address out of bounds")
	}

	v.pc = register(val)

	return nil
}
