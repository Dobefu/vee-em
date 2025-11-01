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

	// OpcodeCMP compares two registers and sets flags.
	OpcodeCMP

	// OpcodeJmpImmediate jumps to an address.
	OpcodeJmpImmediate
	// OpcodeJmpImmediateIfZero jumps to an address if a value is zero.
	OpcodeJmpImmediateIfZero
	// OpcodeJmpImmediateIfNotZero jumps to an address if a value is not zero.
	OpcodeJmpImmediateIfNotZero
	// OpcodeJmpImmediateIfEqual jumps to an address if flags indicate equality.
	OpcodeJmpImmediateIfEqual
	// OpcodeJmpImmediateIfNotEqual jumps to an address if flags indicate inequality.
	OpcodeJmpImmediateIfNotEqual
	// OpcodeJmpImmediateIfGreater jumps to an address if flags indicate greater than.
	OpcodeJmpImmediateIfGreater
	// OpcodeJmpImmediateIfGreaterOrEqual jumps to an address if flags indicate greater than or equal.
	OpcodeJmpImmediateIfGreaterOrEqual
	// OpcodeJmpImmediateIfLess jumps to an address if flags indicate less than.
	OpcodeJmpImmediateIfLess
	// OpcodeJmpImmediateIfLessOrEqual jumps to an address if flags indicate less than or equal.
	OpcodeJmpImmediateIfLessOrEqual
	// OpcodeJmpRegister jumps to an address in a register.
	OpcodeJmpRegister
	// OpcodeJmpRegisterIfZero jumps to an address in a register if a value is zero.
	OpcodeJmpRegisterIfZero
	// OpcodeJmpRegisterIfNotZero jumps to an address in a register if a value is not zero.
	OpcodeJmpRegisterIfNotZero
	// OpcodeJmpRegisterIfEqual jumps to an address in a register if flags indicate equality.
	OpcodeJmpRegisterIfEqual
	// OpcodeJmpRegisterIfNotEqual jumps to an address in a register if flags indicate inequality.
	OpcodeJmpRegisterIfNotEqual
	// OpcodeJmpRegisterIfGreater jumps to an address in a register if flags indicate greater than.
	OpcodeJmpRegisterIfGreater
	// OpcodeJmpRegisterIfGreaterOrEqual jumps to an address in a register if flags indicate greater than or equal.
	OpcodeJmpRegisterIfGreaterOrEqual
	// OpcodeJmpRegisterIfLess jumps to an address in a register if flags indicate less than.
	OpcodeJmpRegisterIfLess
	// OpcodeJmpRegisterIfLessOrEqual jumps to an address in a register if flags indicate less than or equal.
	OpcodeJmpRegisterIfLessOrEqual
)
