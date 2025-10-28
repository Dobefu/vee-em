package vm

func (v *VM) instructionAdd(dest register, src1 register, src2 register) error {
	v.registers[dest] = v.registers[src1] + v.registers[src2]

	return nil
}
