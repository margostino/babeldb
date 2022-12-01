package cli

type Command struct {
	id          string
	pattern     string
	function    func()
	description string
}

func newCommands() []*Command {
	exitCommand := &Command{
		id:          "exit",
		pattern:     "exit",
		function:    exit,
		description: "Exit BabelDB CLI",
	}
	return []*Command{
		exitCommand,
	}
}

func isValidCommand(command string) bool {
	// TODO
	return true && command != ""
}
