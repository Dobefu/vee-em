package vm

import "errors"

func (v *VM) instructionPush(src register) error {
	if v.sp >= uint64(len(v.stack)) {
		return errors.New("stack overflow")
	}

	v.stack[v.sp] = v.registers[src]
	v.sp++

	return nil
}
