# Vee-em

[![Go Version](https://img.shields.io/github/go-mod/go-version/Dobefu/vee-em)](https://go.dev/)
[![License](https://img.shields.io/github/license/Dobefu/vee-em)](https://go.dev/)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Dobefu_vee-em&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=Dobefu_vee-em)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=Dobefu_vee-em&metric=coverage)](https://sonarcloud.io/summary/overall?id=Dobefu_vee-em)
[![Go Report Card](https://goreportcard.com/badge/github.com/Dobefu/vee-em)](https://goreportcard.com/report/github.com/Dobefu/vee-em)

## Overview

The VM is register-based, which means arithmetic operations work with registers
directly rather than a stack.
This makes it faster and more similar to how real CPUs work.

### Architecture

- 32 general-purpose registers
- 512KB heap for memory operations
- 8KB stack for function calls
- Flags register (zero and negative flags)

### Instruction Set

Instructions use registers as operands,
though the exact format varies per instruction.
Arithmetic and bitwise operations work directly on registers,
since this is more efficient than stack-based approaches.

#### Arithmetic Operations

- `ADD`, `SUB`, `MUL`, `DIV`, `MOD` - Standard arithmetic.
  - All take three registers: destination and two source operands.
  - Example: `ADD r0, r1, r2` computes `r0 = r1 + r2`

#### Bitwise Operations

- `AND`, `OR`, `XOR` - Binary logic operations
  - Takes three registers: destination and two source operands
- `NOT` - One's complement (bitwise negation)
  - Takes two registers: destination and source
- `ShiftLeft`, `ShiftRight` - Logical shifts (zeros fill empty bits)
  - Takes three registers: destination, source, and shift amount (from register)
- `ShiftRightArithmetic` - Arithmetic shift right (preserves sign bit)
  - Takes three registers: destination, source, and shift amount (from register)

#### Memory Operations

- `LoadImmediate` - Load a constant value directly into a register (8-byte immediate)
- `LoadRegister` - Copy value from one register to another
- `LoadMemory` - Load from heap address (address in register) into register
- `StoreMemory` - Store register value to heap address

#### Control Flow

- `JmpImmediate`, `JmpRegister` - Unconditional jumps to address
- Conditional jumps based on register value (check if a register value is zero):
  - `JmpImmediateIfZero`, `JmpImmediateIfNotZero`
    - Takes register to check + immediate address
  - `JmpRegisterIfZero`, `JmpRegisterIfNotZero`
    - Takes register to check + register with address
- Conditional jumps based on flags (set by `CMP` instruction):
  - `JmpImmediateIfEqual`, `JmpImmediateIfNotEqual`, `JmpImmediateIfGreater`, `JmpImmediateIfGreaterOrEqual`, `JmpImmediateIfLess`, `JmpImmediateIfLessOrEqual`
    - Takes immediate address (no register argument, just checks flags)
  - `JmpRegisterIfEqual`, `JmpRegisterIfNotEqual`, `JmpRegisterIfGreater`, `JmpRegisterIfGreaterOrEqual`, `JmpRegisterIfLess`, `JmpRegisterIfLessOrEqual`
    - Takes register with address (checks flags)
- `CallImmediate`, `CallRegister`
  - Function calls that push return address to stack
- `Return`
  - Pop return address and jump back

#### Stack Operations

- `Push` - Push register value onto stack (for function call conventions)
- `Pop` - Pop stack value into register

#### Other

- `CMP` - Compares two registers and sets flags for conditional jumps
- `HALT` - Stop VM execution gracefully
- `NOP` - No operation

## Usage

```go
import "github.com/Dobefu/vee-em"

program := []byte{
    // Your bytecode instructions here
}

v := vm.New(
  program,
  vm.WithMagicHeader([]byte("VEE-EM")),
)

err := v.Run()

if err != nil {
  log.Fatalf("Error running VM: %s", err.Error())
}
```

Check out the tests in `run_test.go` for examples of how to construct programs.
