package cli

import (
	"fmt"
	"os"
)

func exit() {
	fmt.Println("bye!")
	os.Exit(0)
}

func help() {
	fmt.Println("help")
}
