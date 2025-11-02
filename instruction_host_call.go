package vm

import (
	"encoding/binary"
	"errors"
)

func (v *VM) instructionHostCall(
	instructionStart register,
	instructionEnd register,
) error {
	if instructionEnd > v.programLen {
		return errors.New("unexpected end of program")
	}

	funcIndex := int64(binary.BigEndian.Uint64(
		v.program[instructionStart+1 : instructionStart+9],
	)) // #nosec: G115

	arg1Reg := register(v.program[instructionStart+9]) & NumRegistersMask
	numArgs := register(v.program[instructionStart+10]) & NumRegistersMask

	if v.hostCallHandler == nil {
		return errors.New("host call handler not set")
	}

	result, err := v.hostCallHandler(funcIndex, arg1Reg, numArgs, v.registers)

	if err != nil {
		return err
	}

	v.registers[arg1Reg] = result

	return nil
}
