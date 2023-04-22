package compiler

import (
	"monkey_cc/ast"
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
	return nil
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
