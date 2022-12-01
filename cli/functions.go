package cli

import (
	"fmt"
	"os"
)

func exit(params interface{}) {
	fmt.Println("bye!")
	os.Exit(0)
}

func help(params interface{}) {
	fmt.Println("help")
}

func createSource(params interface{}) {
	source := params.(string)
	fmt.Printf("source: %s\n", source)
}
