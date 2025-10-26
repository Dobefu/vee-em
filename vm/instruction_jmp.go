package vm

import (
	"errors"
)

func (v *VM) instructionJmp(rawSrc1 register, rawSrc2 register) error {
	mode := Mode(rawSrc2 & 0x3)

	switch mode {
	case ModeImmediate:
		highBits := (rawSrc2 & 0xFC) >> 2
		target := rawSrc1 | (highBits << 8)

		if target > register(len(v.program)) {
			return errors.New("target address out of bounds")
		}

		v.pc = target

	case ModeRegister:
		if rawSrc1 >= register(len(v.registers)) {
			return errors.New("register out of bounds")
		}

		val := v.registers[rawSrc1]

		if val < 0 {
			return errors.New("target address cannot be negative")
		}

		v.pc = register(val)
	}

	return nil
}
