package vm

func (v *VM) decodeInstruction() (opcode Opcode) {
	return Opcode(v.program[v.pc])
}
