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

	// OpcodeLoadImmediate loads an immediate value into a register.
	OpcodeLoadImmediate

	// OpcodeAdd adds two values from the stack.
	OpcodeAdd
	// OpcodeSub subtracts two values from the stack.
	OpcodeSub
	// OpcodeMul multiplies two values from the stack.
	OpcodeMul
	// OpcodeDiv divides two values from the stack.
	OpcodeDiv
	// OpcodeMod takes the modulo of two values from the stack.
	OpcodeMod

	// OpcodeJmp jumps to an address.
	OpcodeJmp
)
