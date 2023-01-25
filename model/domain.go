package model

import (
	"github.com/margostino/babeldb/common"
	"golang.org/x/net/html"
)

type Operator int32
type Type int32
type QueryType int32

const (
	TypeField  string = "type"
	DataField  string = "data"
	TokenField string = "token"
	HrefField  string = "href"
)

var Fields = common.NewStringSlice(TypeField, DataField, TokenField, HrefField)
var AttributeFields = common.NewStringSlice(HrefField)

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

type Attribute struct {
	Key   string
	Value string
}

type Token struct {
	Type       html.TokenType
	Data       string
	Attributes []*Attribute
}

type Source struct {
	Name string
	Url  string
	Page *Page
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
	Fields     *common.StringSlice
	QueryType  QueryType
	Distinct   bool
	Expression *ExpressionTree
}

func (q *Query) InOrderPrint() {
	q.Expression.InOrderPrint(q.Expression.Root)
}

func (q *Query) Match(token *Token) bool {
	return q.Expression.Match(q.Expression.Root, token)
}

type Attributes struct {
	attributes []html.Attribute
}

func NewAttributes(attributes []html.Attribute) *Attributes {
	return &Attributes{
		attributes: attributes,
	}
}

func (s *Attributes) Get(key string) string {
	for _, attribute := range s.attributes {
		if attribute.Key == key {
			return attribute.Val
		}
	}
	return ""
}
