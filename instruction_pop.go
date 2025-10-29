package vm

import "errors"

func (v *VM) instructionPop(rawDest register) error {
	dest := rawDest & NumRegistersMask

	if v.sp == 0 {
		return errors.New("stack underflow")
	}

	v.registers[dest] = v.stack[v.sp-1]
	v.sp--

	return nil
}
