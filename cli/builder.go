package cli

import (
	"github.com/c-bata/go-prompt"
)

func New() *Cli {
	return &Cli{
		Prompt:      "cli@babel",
		Suggestions: getSuggestions(),
		Executor:    newExecutor(),
	}
}

//func getMetadata(commandMap map[string]*Command) []prompt.Suggest {
//	var suggestions = make([]prompt.Suggest, 0)
//	for key, value := range commandMap {
//		var commandText string
//		if value.Args > 0 {
//			commandText = key + " x" + strconv.Itoa(value.Args)
//		} else {
//			commandText = key
//		}
//		suggestion := prompt.Suggest{
//			Text:        commandText,
//			Description: value.Description,
//		}
//		suggestions = append(suggestions, suggestion)
//	}
//	return suggestions
//}

func completer(suggestions []prompt.Suggest) func(d prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
	}
}

//func commandBind(commandsList []CommandConfiguration, cli *Shell) *CommandMap {
//	commands := make(map[string]*Command)
//
//	for _, value := range commandsList {
//		action := getAction(&value, cli)
//		if action != nil && (isValidAction(&value, action.Function) || isValidInputAction(&value, action.InputFunction)) {
//			commands[value.Id] = &Command{
//				Id:          value.Id,
//				Args:        value.Args,
//				Action:      getAction(&value, cli),
//				Pattern:     value.Pattern,
//				Description: value.Description,
//			}
//		} else {
//			log.Printf("Command %s with Args %d, Pattern %s and action %s is not valid\n", value.Id, value.Args, value.Pattern, value.Action)
//		}
//
//	}
//
//	return newCommandMap(commands)
//}
//
//func getAction(command *CommandConfiguration, cli *Shell) *Action {
//	var commandAction *Action = nil
//	if command.Args > 0 {
//		function := cli.ActionOneStringMap[command.Action]
//		commandAction = NewInputAction(function)
//	} else {
//		function := cli.ActionMap[command.Action]
//		commandAction = NewAction(function)
//	}
//	return commandAction
//}
//
//func newCommandMap(commands map[string]*Command) *CommandMap {
//	return &CommandMap{Commands: commands}
//}
//
//func isValidInputAction(command *CommandConfiguration, function func([]string)) bool {
//	return command.Args > 0 && command.Pattern != "" && function != nil
//}
//
//func isValidAction(command *CommandConfiguration, function func()) bool {
//	return command.Args == 0 && command.Pattern == "" && function != nil
//}
