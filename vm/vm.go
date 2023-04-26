package vm

import (
	"encoding/binary"
	"fmt"
	"monkey_cc/code"
	"monkey_cc/compiler"
	"monkey_cc/object"
)

const StackSize = 2048

// 栈式虚拟机，包含三个核心部分：常量、指令、栈
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

func (vm *VM) LastPopped() object.Object {
	return vm.stack[vm.sp+1]
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize-1 {
		return fmt.Errorf("stack overflow")
	}
	vm.sp += 1
	vm.stack[vm.sp] = o
	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp]
	vm.sp--
	return o
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		switch op {
		case code.OpConstant:
			constIdx := binary.BigEndian.Uint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIdx])
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			if err := vm.executeBinaryOperator(op); err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		}
	}
	return nil
}

func (vm *VM) executeBinaryOperator(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return vm.excuteBinaryIntegerOperator(op, left, right)
	}
	return fmt.Errorf("unsupported types for binary operation: %s %s", left.Type(), right.Type())
}

func (vm *VM) excuteBinaryIntegerOperator(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	var result int64

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	}

	return vm.push(&object.Integer{Value: result})
}
