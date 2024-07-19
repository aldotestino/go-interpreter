package main

import (
	"bufio"
	"fmt"
	"go-interpreter/lexer"
	"go-interpreter/parser"
	"go-interpreter/runtime"
	"os"
	"strings"
)

func runFromFile(filename string, env *runtime.Environment) {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}

	input := string(data)

	run(input, env)
}

func runFromRepl(env *runtime.Environment) {
	fmt.Println("\nRepl v0.1")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		input, err := reader.ReadString('\n')

		if err != nil {
			panic(fmt.Sprintf("Something went wrong while reading input: %s", err.Error()))
		}

		input = strings.Trim(input, "\n")

		if input == "" || strings.Contains(input, "exit") {
			os.Exit(0)
		}

		run(input, env)
	}
}

func run(src string, env *runtime.Environment) {
	lex := lexer.NewLexer(src)
	tokens, err := lex.Tokenize()

	if err != nil {
		fmt.Println(err.Error())
	} else {
		pars := parser.NewParser(tokens)
		ast, err := pars.Parse()

		if err != nil {
			fmt.Println(err.Error())
		} else {
			intr := runtime.NewInterpreter()
			res, err := intr.Visit(ast, env)

			if err != nil {
				fmt.Println(err.Error())
			} else if res != nil {
				fmt.Println(res.GetValue())
			}
		}
	}

}

func main() {

	env := runtime.NewEnvironment(nil)

	if len(os.Args) > 1 {
		runFromFile(os.Args[1], env)
	} else {
		runFromRepl(env)
	}
}
