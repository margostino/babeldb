package cli

import (
	"github.com/c-bata/go-prompt"
)

func New() *Cli {
	executor := newExecutor()
	suggestions := executor.newSuggestions()
	return &Cli{
		Prompt:      "cli@babel",
		Suggestions: suggestions,
		Executor:    executor,
	}
}

func completer(suggestions []prompt.Suggest) func(d prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
	}
}
