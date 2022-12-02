package cli

type Command struct {
	id          string
	pattern     string
	function    func(interface{})
	description string
}

func isEndOfCommand(command string) bool {
	return command[len(command)-1:] == ";"
}

func newCommands() []*Command {
	return []*Command{
		{
			id:          "create-source",
			pattern:     "^create source [a-zA-Z0-9_ .=\"\\/*:]+;$",
			function:    createSource,
			description: "List commands available",
		},
	}
}
