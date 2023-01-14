package engine

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
