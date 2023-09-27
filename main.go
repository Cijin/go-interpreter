package main

import (
	"fmt"
	"os"

	"github.com/cijin/go-interpreter/repl"
)

func main() {
	fmt.Print("Welcome to monkey v0.0.0\nPress ctrl-d to exit.\n")

	repl.Start(os.Stdin, os.Stdout)
}
