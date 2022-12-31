package cli

import (
	"github.com/margostino/babeldb/common"
)

func normalize(input string) string {
	return common.NewString(input).
		ToLower().
		ReplaceAll("\n", " ").
		ReplaceAll("\r", " ").
		TrimSpace().
		Value()
}
