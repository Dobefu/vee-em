package vm

import (
	"errors"
)

func (v *VM) decodeInstruction() (
	opcode Opcode,
	dest register,
	rawSrc1 register,
	rawSrc2 register,
	err error,
) {
	if v.pc+3 > register(len(v.program)) {
		return 0, 0, 0, 0, errors.New("unexpected end of program")
	}

	opcode = Opcode(v.program[v.pc])
	dest = register(v.program[v.pc+1])
	rawSrc1 = register(v.program[v.pc+2])
	rawSrc2 = register(v.program[v.pc+3])

	return opcode, dest, rawSrc1, rawSrc2, nil
}
