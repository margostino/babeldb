package cli

type Command struct {
	id          string
	pattern     string
	function    func()
	description string
}

func isValidCommand(command string) bool {
	// TODO
	return true && command != ""
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
	}
}
