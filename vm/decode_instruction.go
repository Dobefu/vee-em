package vm

func (v *VM) decodeInstruction() (
	opcode Opcode,
	dest register,
	src1 register,
	src2 register,
) {
	opcode = Opcode(v.program[v.pc])
	dest = register(v.program[v.pc+1] & 0x1F)
	src1 = register(v.program[v.pc+2] & 0x1F)
	src2 = register(v.program[v.pc+3] & 0x1F)

	return opcode, dest, src1, src2
}
