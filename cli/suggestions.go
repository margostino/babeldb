package cli

import "github.com/c-bata/go-prompt"

func getSuggestions() []prompt.Suggest {
	return []prompt.Suggest{
		{
			Text:        "exit",
			Description: "Exit BabelDB CLI session",
		},
		{
			Text:        "help",
			Description: "List all commands available",
		},
	}
}
