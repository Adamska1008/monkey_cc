package main

import (
	"fmt"
	"monkey_cc/repl"
	"os"
)

func main() {
	fmt.Println("Start monkey lan interpreter.")
	repl.Start(os.Stdin, os.Stdout)
}
