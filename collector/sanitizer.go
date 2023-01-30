package collector

import (
	"github.com/margostino/babeldb/common"
	"regexp"
)

var regExpLeadClose = regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
var regExpInsideClose = regexp.MustCompile(`[\s\p{Zs}]{2,}`)

func sanitize(value string) string {
	return common.NewString(value).
		ReplaceAll("\t", " ").
		ReplaceAll("\n", " ").
		TrimSpace().
		Value()
}
