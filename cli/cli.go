package cli

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"io/ioutil"
	"log"
	"strings"
)

type Cli struct {
	Prompt      string
	Suggestions []prompt.Suggest
	Executor    *Executor
}

func (cli *Cli) Start() {
	welcome()
	for {
		input := cli.prompt()
		normalizedInput := normalize(input)
		if isValidCommand(normalizedInput) {
			cli.Executor.execute(normalizedInput)
		} else {
			fmt.Printf("input %q is not valid\n", normalizedInput)
		}
	}
}

func (cli *Cli) SetPrompt(name string) *Cli {
	cli.Prompt = name
	return cli
}

func (cli *Cli) Help() {
	for _, option := range cli.Suggestions {
		fmt.Printf("[ %s ] - %s\n", option.Text, option.Description)
	}
}

func (cli *Cli) prompt() string {
	prefix := fmt.Sprintf("%s> ", cli.Prompt)
	return prompt.Input(strings.ToLower(prefix), completer(cli.Suggestions))
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
