package main

import (
	"fmt"
	"monkey/repl"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: monkey [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		repl.Script(os.Args[1], os.Stdout)
	} else {
		repl.Start(os.Stdin, os.Stdout)
	}
}
