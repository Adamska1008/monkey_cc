package compiler

import (
	"fmt"
	"monkey_cc/ast"
	"monkey_cc/code"
	"monkey_cc/object"
)

type Compiler struct {
	instructions []byte
	constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: []byte{},
		constants:    []object.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Exp)
		if err != nil {
			return err
		}
	case *ast.InfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "+":
			c.emitOp(code.OpAdd)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IntegerLiteral: // 对于整型常量值，转化为*object.Integer并保存在常量池中
		integer := &object.Integer{Value: node.Value}
		c.emitOp(code.OpConstant, c.pushConstant(integer))
	}
	return nil
}

// 压入常量池，返回在池中的索引
func (c *Compiler) pushConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// 添加操作
func (c *Compiler) emitOp(op code.Opcode, oprands ...int) int {
	ins := code.Make(op, oprands...)
	pos := c.addIns(ins)
	return pos
}

// 将指令码置入指令流中
func (c *Compiler) addIns(ins code.Instructions) int {
	pos := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return pos
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

// 包含编译器生成的Instrucions和编译器求值的Constants
// 用于传输给虚拟机，并在测试中做断言
type Bytecode struct {
	Instructions []byte
	Constants    []object.Object
}
