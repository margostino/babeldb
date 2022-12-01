package cli

import (
	"github.com/margostino/babeldb/common"
)

func normalize(input string) string {
	command := common.NewString(input)
	return command.ToLower().
		TrimSpace().
		Value()
}
