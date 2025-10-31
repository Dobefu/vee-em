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

	// OpcodeAdd adds two values.
	OpcodeAdd
	// OpcodeSub subtracts two values.
	OpcodeSub
	// OpcodeMul multiplies two values.
	OpcodeMul
	// OpcodeDiv divides two values.
	OpcodeDiv
	// OpcodeMod takes the modulo of two values.
	OpcodeMod

	// OpcodeAND performs an AND on two values.
	OpcodeAND
	// OpcodeOR performs an OR on two values.
	OpcodeOR
	// OpcodeXOR performs an exclusive OR on two values.
	OpcodeXOR

	// OpcodeJmpImmediate jumps to an address.
	OpcodeJmpImmediate
	// OpcodeJmpImmediateIfZero jumps to an address if a value is zero.
	OpcodeJmpImmediateIfZero
	// OpcodeJmpImmediateIfNotZero jumps to an address if a value is not zero.
	OpcodeJmpImmediateIfNotZero
	// OpcodeJmpRegister jumps to an address in a register.
	OpcodeJmpRegister
	// OpcodeJmpRegisterIfZero jumps to an address in a register if a value is zero.
	OpcodeJmpRegisterIfZero
	// OpcodeJmpRegisterIfNotZero jumps to an address in a register if a value is not zero.
	OpcodeJmpRegisterIfNotZero
)
