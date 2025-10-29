package vm

func (v *VM) instructionXOR(
	rawDest register,
	rawSrc1 register,
	rawSrc2 register,
) error {
	dest := rawDest & NumRegistersMask
	src1 := rawSrc1 & NumRegistersMask
	src2 := rawSrc2 & NumRegistersMask

	v.registers[dest] = v.registers[src1] ^ v.registers[src2]

	return nil
}
