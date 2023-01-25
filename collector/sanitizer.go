package collector

import "github.com/margostino/babeldb/common"

func sanitize(value string) string {
	return common.NewString(value).
		ReplaceAll("\t", " ").
		ReplaceAll("\n", " ").
		TrimSpace().
		Value()
}
