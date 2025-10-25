// Package vm provides the VM struct.
package vm

type register = uint64

// StackSize is the size of the stack in bytes.
const StackSize = 1024

// VM defines the virtual machine.
type VM struct {
	// The program counter.
	pc register
	// The registers to use when storing or loading data.
	registers [32]int64
	// The bytecode of the program to execute.
	program []byte
	// The stack of the virtual machine.
	stack [StackSize]int64
	// The stack pointer of the virtual machine.
	sp register
}

// New creates a new VM instance.
func New(program []byte) *VM {
	vm := &VM{
		pc:        0,
		registers: [32]int64{},
		program:   program,
		stack:     [StackSize]int64{},
		sp:        0,
	}

	return vm
}
