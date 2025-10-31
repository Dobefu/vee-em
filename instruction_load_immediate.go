package vm

import (
	"encoding/binary"
	"errors"
)

func (v *VM) instructionLoadImmediate(
	instructionStart register,
	instructionLen register,
) error {
	if instructionStart+instructionLen-1 >= v.programLen {
		return errors.New("unexpected end of program")
	}

	dest := register(v.program[instructionStart+1]) & NumRegistersMask

	val := int64(binary.BigEndian.Uint64(
		v.program[instructionStart+2 : instructionStart+instructionLen],
	)) // #nosec: G115

	v.registers[dest] = val

	return nil
}
