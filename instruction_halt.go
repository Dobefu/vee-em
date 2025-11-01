package vm

func (v *VM) instructionHalt(_ register, _ register) error {
	v.pc = v.programLen

	return nil
}
