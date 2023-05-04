package compiler

import (
	"fmt"
	"monkey_cc/ast"
	"monkey_cc/code"
	"monkey_cc/object"
)

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type Compiler struct {
	instructions        []byte
	constants           []object.Object
	lastInstruction     EmittedInstruction // 前一个表达式
	previousInstruction EmittedInstruction // 前两个表达式，仅在回退时使用
}

func New() *Compiler {
	return &Compiler{
		instructions:        []byte{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
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
	case *ast.BlockStatement:
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
		c.emitOp(code.OpPop)
	case *ast.InfixExpression:
		if node.Operator == ">" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}
			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			c.emitOp(code.OpLess)
			return nil
		}
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
		case "-":
			c.emitOp(code.OpSub)
		case "*":
			c.emitOp(code.OpMul)
		case "/":
			c.emitOp(code.OpDiv)
		case "<":
			c.emitOp(code.OpLess)
		case "==":
			c.emitOp(code.OpEqual)
		case "!=":
			c.emitOp(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return nil
		}
		switch node.Operator {
		case "-":
			c.emitOp(code.OpMinus)
		case "!":
			c.emitOp(code.OpBang)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}
		jumpNotTruthyPos := c.emitOp(code.OpJumpNotTruthy, 9999)
		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}
		if c.lastInstruction.Opcode == code.OpPop {
			// 如果If中表达式块的最后末尾有Pop，则移除这个Pop
			// 这是为了使得If中表达式块的数值留在栈中
			c.removeLastPop()
		}
		jumpPos := c.emitOp(code.OpJump, 9999)
		afterConsequencePos := len(c.instructions)
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)
		if node.Alternative == nil {
			c.emitOp(code.OpNull)
		} else {
			err = c.Compile(node.Alternative)
			if err != nil {
				return err
			}
			if c.lastInstruction.Opcode == code.OpPop {
				c.removeLastPop()
			}
		}
		afterAlternativePos := len(c.instructions)
		c.changeOperand(jumpPos, afterAlternativePos)
	case *ast.IntegerLiteral: // 对于整型常量值，转化为*object.Integer并保存在常量池中
		integer := &object.Integer{Value: node.Value}
		c.emitOp(code.OpConstant, c.pushConstant(integer))
	case *ast.Boolean:
		if node.Value {
			c.emitOp(code.OpTrue)
		} else {
			c.emitOp(code.OpFalse)
		}
	}
	return nil
}

// 压入常量池，返回在池中的索引
func (c *Compiler) pushConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// 添加操作，返回指令在指令流中的地址
func (c *Compiler) emitOp(op code.Opcode, oprands ...int) int {
	ins := code.Make(op, oprands...)
	pos := c.addIns(ins)
	c.setLastInstruction(op, pos)
	return pos
}

// 修改某一操作的操作数
func (c *Compiler) changeOperand(opPos, operand int) {
	op := code.Opcode(c.instructions[opPos])
	newIns := code.Make(op, operand)
	c.replaceIns(opPos, newIns)
}

// 记录两条历史指令
func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	c.previousInstruction = c.lastInstruction
	c.lastInstruction = EmittedInstruction{op, pos}
}

// 将指令码置入指令流中，返回指令码的起始地址
func (c *Compiler) addIns(ins code.Instructions) int {
	pos := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return pos
}

func (c *Compiler) replaceIns(pos int, newIns code.Instructions) {
	for i := 0; i < len(newIns); i++ {
		c.instructions[i+pos] = newIns[i]
	}
}

// 移除字节码最后的OpPop指令
func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
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
