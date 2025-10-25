package vm

import (
	"bytes"
	"errors"
)

var errInvalidMagicHeader = errors.New("invalid magic header")

func (v *VM) validateMagicHeader() error {
	if len(v.magicHeader) > len(v.program) {
		return errInvalidMagicHeader
	}

	if !bytes.Equal(v.magicHeader, v.program[:len(v.magicHeader)]) {
		return errInvalidMagicHeader
	}

	return nil
}
