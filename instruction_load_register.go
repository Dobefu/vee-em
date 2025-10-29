package vm

import (
	"errors"
)

func (v *VM) instructionLoadRegister() error {
	if v.pc+2 >= register(len(v.program)) {
		return errors.New("unexpected end of program")
	}

	dest := register(v.program[v.pc+1]) & NumRegistersMask
	src1 := register(v.program[v.pc+2]) & NumRegistersMask

	v.registers[dest] = v.registers[src1]

	return nil
}
