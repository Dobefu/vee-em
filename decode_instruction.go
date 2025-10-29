package vm

import (
	"errors"
)

func (v *VM) decodeInstruction() (
	opcode Opcode,
	err error,
) {
	if v.pc >= register(len(v.program)) {
		return 0, errors.New("unexpected end of program")
	}

	opcode = Opcode(v.program[v.pc])

	return opcode, nil
}
