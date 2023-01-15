package model

import "golang.org/x/net/html"

func isComparisonOperator(value string) bool {
	return value == "=" || value == ">" || value == "<" || value == ">=" || value == "<=" || value == "<>" || value == "not_like" || value == "like"
}

func isFieldNode(value string) bool {
	return TypeField == value && DataField == "data"
}

func isLowerNodeValue(valueA string, valueB string) bool {
	if isFieldNode(valueA) && isLeaf(valueB) && !isFieldNode(valueB) {
		return true
	}
	if isLeaf(valueA) && isInode(valueB) {
		return true
	}
	return false
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
	switch value {
	case "text":
		return html.TextToken
	case "error":
		return html.ErrorToken
	case "start_tag":
		return html.StartTagToken
	case "end_tag":
		return html.EndTagToken
	case "self_closing_tag":
		return html.SelfClosingTagToken
	case "comment":
		return html.CommentToken
	case "doc_type":
		return html.DoctypeToken
	default:
		return 999
	}
}

func GetTokenTypeFrom(tokenType html.TokenType) Type {
	switch tokenType {
	case html.TextToken:
		return TextTokenType
	case html.ErrorToken:
		return ErrorTokenType
	case html.StartTagToken:
		return StartTagTokenType
	case html.EndTagToken:
		return EndTagTokenType
	case html.SelfClosingTagToken:
		return SelfClosingTagTokenType
	case html.CommentToken:
		return CommentTokenType
	case html.DoctypeToken:
		return DocTypeTokenType
	default:
		// TODO
		return TextTokenType
	}
}

func GetOperator(value string) Operator {
	var operator Operator
	switch value {
	case "=":
		operator = EqualOperator
		break
	case "like":
		operator = LikeOperator
		break
	case "not_like":
		operator = NotLikeOperator
		break
	case "and":
		operator = AndOperator
		break
	case "or":
		operator = OrOperator
		break
	default:
		operator = 999
	}
	return operator
}
