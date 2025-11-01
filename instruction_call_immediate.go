package vm

import (
	"encoding/binary"
	"errors"
)

func (v *VM) instructionCallImmediate(
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

	addr := binary.BigEndian.Uint64(
		v.program[instructionStart+1 : instructionEnd],
	)

	if addr >= v.programLen {
		return errors.New("memory address out of bounds")
	}

	v.pc = addr

	return nil
}
