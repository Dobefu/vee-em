package vm

import (
	"errors"
)

func (v *VM) instructionOR() error {
	if v.pc+3 >= register(len(v.program)) {
		return errors.New("unexpected end of program")
	}

	dest := register(v.program[v.pc+1]) & NumRegistersMask
	src1 := register(v.program[v.pc+2]) & NumRegistersMask
	src2 := register(v.program[v.pc+3]) & NumRegistersMask

	v.registers[dest] = v.registers[src1] | v.registers[src2]

	return nil
}
