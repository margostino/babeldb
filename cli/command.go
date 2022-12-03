package cli

type Command struct {
	function    func(interface{})
	description string
}

func isEndOfCommand(command string) bool {
	return command[len(command)-1:] == ";"
}
