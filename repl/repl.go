package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"os"
)

const MONKEY_FACE = `
           __,__
  .--.  .-"     "-.  .--.
 / .. \/  .-. .-.  \/ .. \
| |  '|  /   Y   \  |'  | |
| \   \  \ 0 | 0 /  /   / |
 \ '- ,\.-"""""""-./, -' /
  ''-' /_   ^ ^   _\ '-''
      |  \._   _./  |
      \   \ '~' /   /
       '._ '-=-' _.'
          '-----'
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "\n")
	io.WriteString(out, "Welcome to Monkey.\n")
	io.WriteString(out, "Type \".help\" for more information.\n")
	for {
		fmt.Fprintf(out, "> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		bytes := scanner.Bytes()
		length := len(bytes)
		if length > 0 {
			line := string(bytes)
			if bytes[0] == '.' {
				switch line {
				case ".help":
					io.WriteString(out, ".help    Print this help message\n")
					io.WriteString(out, ".exit    Exit the REPL\n")
				case ".exit":
					os.Exit(0)
				}
			} else {
				run(line, out, env)
			}
		}
	}
}

func Script(in string, out io.Writer) {
	env := object.NewEnvironment()

	bytes, err := os.ReadFile(in)
	if err != nil {
		panic(err)
	}

	run(string(bytes), out, env)
}

func run(input string, out io.Writer, env *object.Environment) {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			io.WriteString(out, "parser errors:")
			io.WriteString(out, msg)
			io.WriteString(out, "\n")
		}
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}
