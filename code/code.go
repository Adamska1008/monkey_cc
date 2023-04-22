package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			return out.String()
		}
		operands, width := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmt(def, operands))
		i += 1 + width
	}
	return out.String()
}

// 将单个指令转化为字符串
func (ins Instructions) fmt(def *Definition, operands []int) string {
	switch len(def.OperandWidths) {
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}
	return fmt.Sprintf("ERROR: unmatched operand numbers for %s to be %d", def.Name, len(def.OperandWidths))
}

type Opcode byte

const (
	OpConstant Opcode = iota
)

type Definition struct {
	Name          string
	OperandWidths []int
}

// 将操作码和操作数转换为字节流
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}
	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)
	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 4:
			binary.BigEndian.PutUint32(instruction[offset:], uint32(o))
		case 8:
			binary.BigEndian.PutUint64(instruction[offset:], uint64(o))
		}
		offset += width
	}
	return instruction
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}}, // 常量操作码，操作数表示常量在常量池中的位置
}

// 查找对应操作码的定义
// 定义包括操作码名称和操作数数量和位宽
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

// 读取操作数
// 返回值为(操作数数组，总字节长度)
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0
	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(binary.BigEndian.Uint16(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}
