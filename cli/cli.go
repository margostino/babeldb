package cli

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/margostino/babeldb/common"
	"github.com/margostino/babeldb/engine"
	"github.com/xwb1989/sqlparser"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Cli struct {
	prompt      string
	engine      *engine.Engine
	suggestions []prompt.Suggest
}

func (cli *Cli) Start() {
	welcome()

	var input string
	var parseErr error
	var statement *sqlparser.Statement
	var isNewLine = false
	var inputs = make([]string, 0)

	for {
		if isNewLine {
			input = cli.printNewLine()
			//inputs = make([]string, 0)
		} else {
			input = cli.printPromptAndGetInput()
		}

		normalizedInput := normalize(input)

		if normalizedInput == "" {
			continue
		} else if !isNewLine && normalizedInput == "help" {
			cli.help()
		} else if !isNewLine && normalizedInput == "exit" {
			cli.exit()
		} else if !isEndOfCommand(normalizedInput) {
			isNewLine = true
			inputs = append(inputs, normalizedInput)
		} else if cli.isValidCommand(normalizedInput) {
			isNewLine = false
			if len(inputs) > 0 {
				multilineInput := fmt.Sprintf("%s %s", strings.Join(inputs, " "), normalizedInput)
				statement, parseErr = cli.engine.Parse(multilineInput)
			} else {
				statement, parseErr = cli.engine.Parse(normalizedInput)
			}

			if !common.IsError(parseErr, "when parsing input") {
				cli.execute(statement)
			}

		} else {
			fmt.Printf("input %q is not valid\n", normalizedInput)
		}

	}
}

func (cli *Cli) isValidCommand(command string) bool {
	// TODO
	return true
}

func (cli *Cli) help() {
	for _, option := range cli.suggestions {
		fmt.Printf("[ %s ] - %s\n", option.Text, option.Description)
	}
}

func (cli *Cli) exit() {
	fmt.Println("bye!")
	os.Exit(0)
}

func (cli *Cli) printPromptAndGetInput() string {
	prefix := fmt.Sprintf("%s> ", cli.prompt)
	return prompt.Input(strings.ToLower(prefix), completer(cli.suggestions))
}

func (cli *Cli) printNewLine() string {
	prefix := "           | "
	return prompt.Input(strings.ToLower(prefix), completer(cli.suggestions))
}

func (cli *Cli) execute(statement *sqlparser.Statement) {
	//query.Solver(cli.engine, query.Params)
}

func isEndOfCommand(command string) bool {
	return command[len(command)-1:] == ";"
}

func welcome() {
	file, err := ioutil.ReadFile("banner.txt")
	if err != nil {
		log.Printf("Failure when reading banner file: %v ", err)
	}
	banner := string(file)
	fmt.Printf("\n")
	fmt.Printf("%s\n", banner)
	fmt.Printf("Welcome to BabelDB CLI! - Version 0.1\n\n")
	log.Printf("DB engine healthy\n\n") // TODO
}
