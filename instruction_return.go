package vm

import (
	"errors"
)

func (v *VM) instructionReturn(_ register, _ register) error {
	if v.sp == 0 {
		return errors.New("stack underflow")
	}

	returnAddr := register(v.stack[v.sp-1]) // #nosec: G115
	v.sp--

	if returnAddr >= v.programLen {
		return errors.New("memory address out of bounds")
	}

	v.pc = returnAddr

	return nil
}
