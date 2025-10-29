package vm

import (
	"encoding/binary"
	"errors"
)

func (v *VM) instructionJmpImmediate() error {
	if v.pc+8 >= register(len(v.program)) {
		return errors.New("unexpected end of program")
	}

	val := int64(binary.BigEndian.Uint64(v.program[v.pc+1 : v.pc+9])) // #nosec: G115

	if val < 0 || val >= int64(len(v.program)) {
		return errors.New("memory address out of bounds")
	}

	v.pc = register(val)

	return nil
}
