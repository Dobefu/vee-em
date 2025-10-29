package vm

import (
	"errors"
)

func (v *VM) instructionPop() error {
	if v.pc+1 >= register(len(v.program)) {
		return errors.New("unexpected end of program")
	}

	rawDest := register(v.program[v.pc+1])
	dest := rawDest & NumRegistersMask

	if v.sp == 0 {
		return errors.New("stack underflow")
	}

	v.registers[dest] = v.stack[v.sp-1]
	v.sp--

	return nil
}
