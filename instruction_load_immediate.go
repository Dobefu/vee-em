package vm

import (
	"encoding/binary"
	"errors"
)

func (v *VM) instructionLoadImmediate() error {
	if v.pc+9 >= register(len(v.program)) {
		return errors.New("unexpected end of program")
	}

	dest := register(v.program[v.pc+1]) & NumRegistersMask
	val := int64(binary.BigEndian.Uint64(v.program[v.pc+2 : v.pc+10])) // #nosec: G115

	v.registers[dest] = val

	return nil
}
