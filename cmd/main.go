package main

import (
	"bufio"
	"fmt"
	"go-interpreter/lexer"
	"go-interpreter/parser"
	"go-interpreter/utils"
	"os"
	"strings"
)

func runFromFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}

	input := string(data)

	run(filename, input)
}

func runFromRepl() {
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

		run("<stdin>", input)
	}
}

func run(fn string, src string) {
	lex := lexer.NewLexer(fn, src)
	tokens, err := lex.Tokenize()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	pars := parser.NewParser(tokens)
	ast := pars.Parse()

	if ast.Error != nil {
		fmt.Println(ast.Error.Error())
		os.Exit(1)
	}

	utils.LogAsJSON("AST", ast)
}

func main() {
	if len(os.Args) > 1 {
		runFromFile(os.Args[1])
	} else {
		runFromRepl()
	}
}
