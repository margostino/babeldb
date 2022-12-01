package cli

type Command struct {
	id          string
	pattern     string
	function    func(interface{})
	description string
}

func isValidCommand(command string) bool {
	// TODO
	return true
}

func newCommands() []*Command {
	return []*Command{
		{
			id:          "exit",
			pattern:     "exit",
			function:    exit,
			description: "Exit BabelDB CLI",
		},
		{
			id:          "help",
			pattern:     "help",
			function:    help,
			description: "List commands available",
		},
		{
			id:          "create-source",
			pattern:     "^create source [a-zA-Z0-9_ .=\"\\/*]+$",
			function:    createSource,
			description: "List commands available",
		},
	}
}
