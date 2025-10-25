package vm

// Opcode defines operation codes.
type Opcode byte

const (
	// OpcodeNop does not do anything.
	OpcodeNop Opcode = iota

	// OpcodePush pushes a value onto the stack.
	OpcodePush
	// OpcodePop pops a value off the stack.
	OpcodePop

	// OpcodeAdd adds two values from the stack.
	OpcodeAdd
	// OpcodeSub subtracts two values from the stack.
	OpcodeSub
)
