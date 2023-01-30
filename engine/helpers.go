package engine

import (
	"github.com/margostino/babeldb/common"
	"github.com/margostino/babeldb/model"
)

func isComparisonOperator(value string) bool {
	return value == "=" || value == ">" || value == "<" || value == ">=" || value == "<=" || value == "<>" || value == "not_like"
}

func isLogicalOperator(value string) bool {
	return value == "and" || value == "or" // TODO
}

func isInode(value string) bool {
	return isLogicalOperator(value) || isComparisonOperator(value)
}

func isLeaf(value string) bool {
	return !isLogicalOperator(value) && !isComparisonOperator(value)
}

func shouldCreateSource(input string) bool {
	return common.NewString(input).
		ToLower().
		HasPrefix(model.CreateSource)
}
