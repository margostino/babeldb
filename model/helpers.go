package model

import "golang.org/x/net/html"

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

func GetTokenType(value string) html.TokenType {
	var tokenType html.TokenType
	switch value {
	case "text":
		tokenType = html.TextToken
		break
	case "error":
		tokenType = html.TextToken
		break
	case "start_tag":
		tokenType = html.StartTagToken
		break
	case "end_tag":
		tokenType = html.EndTagToken
		break
	case "self_closing_tag":
		tokenType = html.SelfClosingTagToken
		break
	case "comment":
		tokenType = html.CommentToken
		break
	case "doc_type":
		tokenType = html.DoctypeToken
		break
	default:
		tokenType = 999
	}
	return tokenType
}
