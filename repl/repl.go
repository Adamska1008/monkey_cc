package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey_cc/evaluator"
	"monkey_cc/lexer"
	"monkey_cc/parser"
)

const PROMPT = ">> "

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
		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
		// io.WriteString(out, program.String())
		// io.WriteString(out, "\n")
		//for {
		//	t := l.NextToken()
		//	if t.Type == token.EOF {
		//		break
		//	}
		//	fmt.Fprintf(out, "%v\n", *t)
		//}
	}
}
