package cli

import (
	"errors"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/margostino/babeldb/common"
	"regexp"
	"strings"
)

type Executor struct {
	commands []*Command
}

func newExecutor() *Executor {
	return &Executor{
		commands: newCommands(),
	}
}

func (e *Executor) execute(input string) {
	command, err := e.lookup(input)
	if !common.IsError(err, "command lookup failed") {
		command.function(input)
	}
}

func (e *Executor) lookup(input string) (*Command, error) {
	for _, value := range e.commands {
		// TODO: validate not null/empty Pattern
		match, _ := regexp.MatchString(value.pattern, input)
		if match {
			return value, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("command not found for input [%s]", input))
}

func (e *Executor) newSuggestions() []prompt.Suggest {
	var suggestions = make([]prompt.Suggest, 0)
	for _, command := range e.commands {
		text := strings.ReplaceAll(command.pattern, "^", "")
		suggestion := prompt.Suggest{
			Text:        text,
			Description: command.description,
		}
		suggestions = append(suggestions, suggestion)
	}
	return suggestions
}
