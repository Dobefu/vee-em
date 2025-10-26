package vm

// Mode defines the mode of an instruction.
type Mode byte

const (
	// ModeImmediate is the immediate mode.
	ModeImmediate Mode = iota
	// ModeRegister is the register mode.
	ModeRegister
)
