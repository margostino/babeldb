package cli

import (
	"fmt"
	"github.com/margostino/babeldb/common"
	"strings"
)

func createSource(input interface{}) {
	command := common.NewString(input.(string)).
		ReplaceAll("create source", "").
		TrimSpace().
		Value()
	parts := strings.Split(command, " ")

	if len(parts) > 0 && parts[0] != "" {
		fmt.Printf("source name: %s\n", parts[0])
	} else {
		fmt.Printf("source name invalid: %s\n", input)
	}

}
