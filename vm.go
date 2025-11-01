// Package vm provides the VM struct.
package vm

// Option is a function that can be used to configure the VM.
type Option func(*VM)

type register = uint64

// NumRegisters is the numbers of registers in the CPU.
// Please note that this has to be a power of two.
const NumRegisters = 32

// NumRegistersMask is the bitmask to use when masking a raw register.
const NumRegistersMask = NumRegisters - 1

// StackSize is the size of the stack in bytes.
const StackSize = 1024

// HeapSize is the size of the heap in bytes.
const HeapSize = 65536

// flags defines the flags of the CPU.
type flags struct {
	// Whether the result of the last operation was zero.
	isZero bool
	// Whether the result of the last operation was negative.
	isNegative bool
}

// VM defines the virtual machine.
type VM struct {
	// The magic header of the program.
	magicHeader []byte
	// The program counter.
	pc register
	// The registers to use when storing or loading data.
	registers [NumRegisters]int64
	// The bytecode of the program to execute.
	program []byte
	// The length of the program.
	programLen register
	// The stack of the virtual machine.
	stack [StackSize]int64
	// The stack pointer of the virtual machine.
	sp register
	// The heap memory of the virtual machine.
	heap [HeapSize]int64
	// The flags register.
	flags flags
}

// New creates a new VM instance.
func New(program []byte, options ...Option) *VM {
	vm := &VM{
		magicHeader: []byte{},
		pc:          0,
		registers:   [NumRegisters]int64{},
		program:     program,
		programLen:  register(len(program)),
		stack:       [StackSize]int64{},
		sp:          0,
		heap:        [HeapSize]int64{},
		flags: flags{
			isZero:     false,
			isNegative: false,
		},
	}

	for _, option := range options {
		option(vm)
	}

	return vm
}

// WithMagicHeader adds a magic header to the program.
func WithMagicHeader(header []byte) Option {
	return func(v *VM) {
		v.magicHeader = header
		v.pc = register(len(header))
	}
}
