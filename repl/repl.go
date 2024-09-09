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

	for {
		fmt.Fprintf(out, ">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		run(scanner.Text(), out, env)
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
		io.WriteString(out, MONKEY_FACE)
		io.WriteString(out, "parser errors:\n")
		for _, msg := range p.Errors() {
			io.WriteString(out, msg)
			io.WriteString(out, "\n")
		}
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}
