package vm

import "errors"

func (v *VM) instructionPush(rawSrc1 register) error {
	src1 := rawSrc1 & NumRegistersMask

	if v.sp >= uint64(len(v.stack)) {
		return errors.New("stack overflow")
	}

	v.stack[v.sp] = v.registers[src1]
	v.sp++

	return nil
}
