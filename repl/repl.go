package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey_cc/compiler"
	"monkey_cc/lexer"
	"monkey_cc/parser"
	"monkey_cc/vm"
)

const PROMPT = ">> "

// Start function based on Compiler and VM
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			for _, msg := range p.Errors() {
				io.WriteString(out, "\t"+msg+"\n")
			}
		}
		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}
		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
			continue
		}
		io.WriteString(out, machine.LastPopped().Inspect())
		io.WriteString(out, "\n")
	}
}

// Start function based on evaluater and environment
// func Start(in io.Reader, out io.Writer) {
// 	scanner := bufio.NewScanner(in)
// 	env := object.NewEnvironment()

// 	for {
// 		fmt.Fprint(out, PROMPT)
// 		scanned := scanner.Scan()
// 		if !scanned {
// 			return
// 		}
// 		line := scanner.Text()
// 		l := lexer.New(line)
// 		p := parser.New(l)
// 		program := p.ParseProgram()
// 		if len(p.Errors()) != 0 {
// 			for _, msg := range p.Errors() {
// 				io.WriteString(out, "\t"+msg+"\n")
// 			}
// 		}
// 		evaluated := evaluator.Eval(program, env)
// 		if evaluated != nil {
// 			io.WriteString(out, evaluated.Inspect())
// 			io.WriteString(out, "\n")
// 		}
// 	}
// }
