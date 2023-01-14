package model

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

type Operator int32
type VarType int32
type QueryType int32

const (
	EqualOperator Operator = iota
	LikeOperator
	NotLikeOperator
	StringType VarType = iota
	TokenType
	SelectType QueryType = iota
	CreateType
)

type QueryVar struct {
	Operator Operator
	VarType  VarType
	Field    string
	Value    interface{}
}

type ExpressionTree struct {
	Root *ExpressionNode
}

type ExpressionNode struct {
	Key     interface{}
	VarType VarType
	Left    *ExpressionNode
	Right   *ExpressionNode
}

type Query struct {
	Source     string
	Url        string
	Schedule   string
	QueryType  QueryType
	Vars       []*QueryVar
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
	return q.Expression.Match(q.Expression.Root, "", token)
}

func (t *ExpressionTree) Match(node *ExpressionNode, key interface{}, token *Token) bool {
	var match bool
	//var key string
	if node != nil {
		match = t.Match(node.Left, node.Key, token)

		if isComparisonOperator(node.Key.(string)) {
			field := node.Left.Key
			value := node.Right.Key

			if field == "type" {
				switch node.Key {
				case "=":
					match = token.Type.String() == value
				case "notLike":
					match = token.Type.String() == value
				}
			} else if field == "data" {
				switch node.Key {
				case "=":
					match = token.Type.String() == value
				case "notLike":
					match = !strings.Contains(token.Data, value.(string))
				}
			}

		}
		fmt.Printf("InOrder: %s (%t)\n", node.Key, match)
		//if isLeaf(node.Key) && isLeafKey {
		//	field = node.Key
		//} else if isLeaf(node.Key) && !isLeafKey {
		//	value = node.Key
		//} else if isComparisonOperator(node.Key) {
		//	switch node.Key {
		//	case "=":
		//		operator = EqualOperator
		//	case "like":
		//		operator = LikeOperator
		//	case "not_like":
		//		operator = NotLikeOperator
		//	}
		//} else if isLogicalOperator(node.Key) {
		//	println("ds")
		//}
		match = t.Match(node.Right, node.Key, token)
	}

	return match
}

func (t *ExpressionTree) InOrderPrint(node *ExpressionNode) {
	if node != nil {
		t.InOrderPrint(node.Left)
		fmt.Printf("InOrder: %s\n", node.Key)
		t.InOrderPrint(node.Right)
	}
}

func (t *ExpressionTree) Insert(key string) {
	if t.Root == nil {
		t.Root = &ExpressionNode{Key: key}
	} else {
		t.Root.Insert(key)
	}
}

func (n *ExpressionNode) Insert(key string) {
	if isLeaf(key) {
		if n.Left == nil {
			n.Left = &ExpressionNode{Key: key}
		} else if n.Right == nil {
			n.Right = &ExpressionNode{Key: key}
			if isComparisonOperator(n.Key.(string)) && n.Left.Key == "type" {
				n.Right.VarType = TokenType
			}
		}
	} else if isInode(key) {
		leafNode := *n
		n.Key = key
		n.Left = &leafNode
		n.Right = nil
	} else {
		if n.Right == nil {
			n.Right = &ExpressionNode{Key: key}
		} else {
			n.Right.Insert(key)
		}
	}

	//if isInode(key) {
	//	leafNode := *n
	//	n.Key = key
	//	n.Left = &leafNode
	//	n.Right = nil
	//} else if isLeaf(key) {
	//	if n.Left == nil {
	//		n.Left = &ExpressionNode{Key: key}
	//	} else if n.Right == nil {
	//		n.Right = &ExpressionNode{Key: key}
	//	}
	//} else {
	//	if n.Right == nil {
	//		n.Right = &ExpressionNode{Key: key}
	//	} else {
	//		n.Right.Insert(key)
	//	}
	//}
}
