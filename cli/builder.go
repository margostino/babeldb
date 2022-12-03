package cli

import (
	"github.com/c-bata/go-prompt"
)

func New() *Cli {
	return &Cli{
		prompt:      "cli@babel",
		suggestions: newSuggestions(),
	}
}

func completer(suggestions []prompt.Suggest) func(d prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
	}
}

func newSuggestions() []prompt.Suggest {
	return []prompt.Suggest{
		{
			Text:        "CREATE SOURCE {name} FROM {url}",
			Description: "Create new data source for a given URL",
		},
	}
}
