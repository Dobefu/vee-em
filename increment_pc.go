package vm

func (v *VM) incrementPC() {
	v.pc += 4
}
