package parser

import (
	"fmt"
	"strings"
)

var depth = 0

func printTrace(msg string) {
	fmt.Printf("%s%s\n", strings.Repeat("\t", depth), msg)
}

func trace(msg string) string {
	printTrace("BEGIN " + msg)
	depth = depth + 1
	return msg
}

func untrace(msg string) {
	depth = depth - 1
	printTrace("END " + msg)
}
