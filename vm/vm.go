// Package vm provides the VM struct.
package vm

type register = uint32

// VM defines the virtual machine.
type VM struct {
	// The program counter.
	pc register
	// The registers to use when storing or loading data.
	registers [32]int64
	// The bytecode of the program to execute.
	program []byte
	stack   [1024]int64
	sp      register
}

// New creates a new VM instance.
func New(program []byte) *VM {
	vm := &VM{
		pc:        0,
		registers: [32]int64{},
		program:   program,
		stack:     [1024]int64{},
		sp:        0,
	}

	return vm
}
