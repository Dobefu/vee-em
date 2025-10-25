// Package vm provides the VM struct.
package vm

type register = uint32

// VM defines the virtual machine.
type VM struct {
	pc register
}

// New creates a new VM instance.
func New() *VM {
	vm := &VM{
		pc: 0,
	}

	return vm
}
