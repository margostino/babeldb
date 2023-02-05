package model

import (
	"github.com/margostino/babeldb/common"
)

type Query struct {
	Source     string
	Url        string
	Schedule   string
	Fields     *common.StringSlice
	QueryType  QueryType
	Distinct   bool
	Expression *ExpressionTree
}

type QueryResults struct {
	Sources []*Source
	Page    *Page
}

func NewQuery() *Query {
	return &Query{
		Fields: common.NewStringSlice(),
	}
}

func (q *Query) InOrderPrint() {
	q.Expression.InOrderPrint(q.Expression.Root)
}
