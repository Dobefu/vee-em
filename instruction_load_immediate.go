package vm

func (v *VM) instructionLoadImmediate(rawDest register) error {
	dest := rawDest & NumRegistersMask

	highByte := int64(v.program[v.pc+2]) << 8
	lowByte := int64(v.program[v.pc+3])
	v.registers[dest] = highByte | lowByte

	return nil
}
