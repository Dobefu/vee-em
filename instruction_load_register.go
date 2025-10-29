package vm

func (v *VM) instructionLoadRegister(rawDest register, rawSrc1 register) error {
	dest := rawDest & NumRegistersMask
	src1 := rawSrc1 & NumRegistersMask

	v.registers[dest] = v.registers[src1]

	return nil
}
