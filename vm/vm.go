package vm

import (
	"monkey_cc/code"
	"monkey_cc/compiler"
	"monkey_cc/object"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    -1,
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == -1 {
		return nil
	}
	return vm.stack[vm.sp]
}

func (vm *VM) Run() error {
	return nil
}
