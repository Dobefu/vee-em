package vm

import "errors"

func (v *VM) instructionPop(dest register) error {
	if v.sp == 0 {
		return errors.New("stack underflow")
	}

	v.registers[dest] = v.stack[v.sp-1]
	v.sp--

	return nil
}
