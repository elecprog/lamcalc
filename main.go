package main

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
	"github.com/elecprog/cLC/lambda"
)

var globals = map[string]lambda.Term{}

func main() {
	// A warm welcome
	showInfo()

	// Load files
	if len(os.Args) > 1 {
		loadFiles(os.Args[1:])
		fmt.Print("Switching to interactive mode...\n\n")
	}

	commandline, _ := readline.New("(cLC) ")
	defer commandline.Close()

	for {
		command, _ := commandline.Readline()
		stmnt, err := parseStatement(command)

		if err != nil {
			printError(err)
		} else {
			executeStatement(stmnt)
		}
	}
}