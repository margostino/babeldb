package main

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/margostino/babeldb/shell"
	"strconv"
)

func main() {
	//cmd.Execute()
	shell.Welcome()
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("Invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Number",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)
}
