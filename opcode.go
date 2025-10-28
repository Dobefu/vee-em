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
	// OpcodeLoadRegister loads a value from a register into another register.
	OpcodeLoadRegister

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

	// OpcodeAND performs an AND on two values from the stack.
	OpcodeAND
	// OpcodeXOR performs an exclusive OR on two values from the stack.
	OpcodeXOR

	// OpcodeJmpImmediate jumps to an address.
	OpcodeJmpImmediate
	// OpcodeJmpRegister jumps to an address in a register.
	OpcodeJmpRegister
)
