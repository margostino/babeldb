package model

import (
	"fmt"
	"github.com/margostino/babeldb/common"
	"golang.org/x/net/html"
	"strings"
)

type Operator int32
type Type int32
type QueryType int32

const (
	TypeField  string = "type"
	DataField  string = "data"
	TokenField string = "token"
)

var Fields = common.NewStringSlice(TypeField, DataField, TokenField)

const (
	EqualOperator Operator = iota
	LikeOperator
	NotLikeOperator
	AndOperator
	OrOperator
	StringType Type = iota
	TokenType
	TextTokenType
	ErrorTokenType
	StartTagTokenType
	EndTagTokenType
	SelfClosingTagTokenType
	CommentTokenType
	DocTypeTokenType
	SelectType QueryType = iota
	CreateType
)

type QueryVar struct {
	Operator Operator
	Type     Type
	Field    string
	Value    interface{}
}

type ExpressionTree struct {
	Root *ExpressionNode
}

type ExpressionNodeKey struct {
	Value    string
	Type     Type
	Operator Operator
}

type ExpressionNode struct {
	Key   *ExpressionNodeKey
	Left  *ExpressionNode
	Right *ExpressionNode
}

type Query struct {
	Source     string
	Url        string
	Schedule   string
	QueryType  QueryType
	Distinct   bool
	Expression *ExpressionTree
}

type Attributes struct {
	Key   string
	Value string
}

type Token struct {
	Type       html.TokenType
	Data       string
	Attributes []*Attributes
}

type Source struct {
	Name   string
	Url    string
	Tokens []*Token
}

func (q *Query) InOrderPrint() {
	q.Expression.InOrderPrint(q.Expression.Root)
}

func (q *Query) Match(token *Token) bool {
	return q.Expression.Match(q.Expression.Root, token)
}

func (n *ExpressionNode) GetKey() string {
	return n.Key.Value
}

func (n *ExpressionNode) GetType() Type {
	return n.Key.Type
}

func (n *ExpressionNode) GetOperator() Operator {
	return n.Key.Operator
}

func (n *ExpressionNode) IsComparisonOperatorNode() bool {
	return isComparisonOperator(n.GetKey())
}

func (n *ExpressionNode) IsInode() bool {
	return isInode(n.GetKey())
}

func (n *ExpressionNode) IsValueFieldNode() bool {
	return !isFieldNode(n.GetKey()) && !n.IsInode()
}

func (n *ExpressionNode) IsLogicalOperatorNode() bool {
	return isLogicalOperator(n.GetKey())
}

func (n *ExpressionNode) isParamValue() bool {
	return len(n.GetKey()) == 2 && n.GetKey()[0:1] == ":"
}

func (t *ExpressionTree) GetParamNode(node *ExpressionNode) *ExpressionNode {

	if node.isParamValue() {
		return node
	}

	return t.GetParamNode(node.Right)
}

func (t *ExpressionTree) Match(node *ExpressionNode, token *Token) bool {
	var match bool
	if !isLeaf(node.GetKey()) {
		if node.IsComparisonOperatorNode() {
			var match bool
			field := node.Left.GetKey()
			value := node.Right.GetKey()

			if field == "type" {
				switch node.GetOperator() {
				case EqualOperator:
					match = token.Type == GetTokenType(value)
				case NotLikeOperator:
					match = token.Type.String() == value
				}
			} else if field == "data" {
				data := strings.ToLower(token.Data)
				value = strings.ReplaceAll(value, "%", "")
				switch node.GetOperator() {
				case EqualOperator:
					match = data == value
				case LikeOperator:
					match = strings.Contains(data, value)
				case NotLikeOperator:
					match = !strings.Contains(token.Data, value)
				}
			}
			return match
		}
		match = t.Match(node.Left, token)

		if node.GetOperator() == AndOperator {
			match = match && t.Match(node.Right, token)
		} else if node.GetOperator() == OrOperator {
			match = match || t.Match(node.Right, token)
		} else {
			// TODO
		}

	}

	return match
}

func (t *ExpressionTree) InOrderPrint(node *ExpressionNode) {
	if node != nil {
		t.InOrderPrint(node.Left)
		fmt.Printf("InOrder: %s\n", node.GetKey())
		t.InOrderPrint(node.Right)
	}
}

func (t *ExpressionTree) Insert(key string) {
	if t.Root == nil {
		t.Root = &ExpressionNode{Key: &ExpressionNodeKey{Value: key}}
	} else {
		t.Root.Insert(key)
	}
}

func (n *ExpressionNode) Insert(key string) {
	if isLeaf(key) {
		if n.Left == nil {
			n.Left = &ExpressionNode{Key: &ExpressionNodeKey{Value: key}}
		} else if n.Right == nil {
			n.Right = &ExpressionNode{Key: &ExpressionNodeKey{Value: key}}
			if isComparisonOperator(n.GetKey()) && n.Left.GetKey() == TypeField {
				n.Right.Key.Type = TokenType
			}
		} else {
			n.Right.Insert(key)
		}
	} else {
		operator := GetOperator(key)

		if isLogicalOperator(key) {
			leafNode := *n
			node := &ExpressionNode{
				Key: &ExpressionNodeKey{
					Value:    key,
					Operator: operator,
				},
				Left:  &leafNode,
				Right: nil,
			}
			n.Key = node.Key
			n.Left = node.Left
			n.Right = node.Right
		} else {
			leafNode := *n
			if isLogicalOperator(n.GetKey()) {
				n.Right = &ExpressionNode{
					Key: &ExpressionNodeKey{
						Value:    key,
						Operator: operator,
					},
					Left:  n.Right,
					Right: nil,
				}
			} else {
				n.Key = &ExpressionNodeKey{
					Value:    key,
					Operator: operator,
				}
				n.Left = &leafNode
				n.Right = nil
			}
		}
	}

}
