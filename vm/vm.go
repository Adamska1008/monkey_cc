package vm

import (
	"encoding/binary"
	"fmt"
	"monkey_cc/code"
	"monkey_cc/compiler"
	"monkey_cc/object"
)

const (
	StackSize  = 2048
	GlobalSize = 65536
)

var (
	True  = object.TRUE
	False = object.FALSE
	Null  = object.NULL
)

// 栈式虚拟机，包含三个核心部分：常量、指令、栈
type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack   []object.Object
	globals []object.Object
	sp      int
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack:   make([]object.Object, StackSize),
		globals: make([]object.Object, GlobalSize),
		sp:      -1,
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
		case code.OpBang:
			if err := vm.executeBangOperator(op); err != nil {
				return err
			}
		case code.OpMinus:
			if err := vm.executeMinusOperator(op); err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			if err := vm.executeBinaryOperator(op); err != nil {
				return err
			}
		case code.OpEqual, code.OpNotEqual, code.OpLess:
			if err := vm.executeComparison(op); err != nil {
				return err
			}
		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}
		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		case code.OpJumpNotTruthy:
			operand := int(binary.BigEndian.Uint16(vm.instructions[ip+1:]))
			ip += 2
			condition := vm.pop()
			if !isTruthy(condition) {
				ip = operand - 1
			}
		case code.OpJump:
			operand := int(binary.BigEndian.Uint16(vm.instructions[ip+1:]))
			ip = operand - 1
		case code.OpSetGlobal:
			globalIdx := int(binary.BigEndian.Uint16(vm.instructions[ip+1:]))
			ip += 2
			vm.globals[globalIdx] = vm.pop()
		case code.OpGetGlobal:
			globalIdx := int(binary.BigEndian.Uint16(vm.instructions[ip+1:]))
			ip += 2
			err := vm.push(vm.globals[globalIdx])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (vm *VM) executeBinaryOperator(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return vm.executeBinaryIntegerOperator(op, left, right)
	}
	return fmt.Errorf("unsupported types for binary operation: %s %s", left.Type(), right.Type())
}

func (vm *VM) executeBinaryIntegerOperator(op code.Opcode, left, right object.Object) error {
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

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	} else if left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ {
		return vm.executeBooleanComparison(op, left, right)
	}
	return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	var result bool
	switch op {
	case code.OpEqual:
		result = leftValue == rightValue
	case code.OpNotEqual:
		result = leftValue != rightValue
	case code.OpLess:
		result = leftValue < rightValue
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
	}
	vm.push(&object.Boolean{Value: result})
	return nil
}

func (vm *VM) executeBooleanComparison(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value
	var result bool
	switch op {
	case code.OpEqual:
		result = leftValue == rightValue
	case code.OpNotEqual:
		result = leftValue != rightValue
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
	}
	vm.push(&object.Boolean{Value: result})
	return nil
}

func (vm *VM) executeBangOperator(op code.Opcode) error {
	operand := vm.pop()
	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) executeMinusOperator(op code.Opcode) error {
	operand := vm.pop()
	switch operand := operand.(type) {
	case *object.Integer:
		return vm.push(&object.Integer{Value: -operand.Value})
	default:
		return fmt.Errorf("unknown operator: %d (%s)", op, operand.Type())
	}
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	default:
		return true
	}
}
