package engine

import (
	"errors"
	"fmt"
	"github.com/margostino/babeldb/collector"
	"github.com/margostino/babeldb/common"
	"github.com/margostino/babeldb/model"
	"github.com/margostino/babeldb/storage"
	"github.com/robfig/cron/v3"
	"github.com/xwb1989/sqlparser"
	"github.com/xwb1989/sqlparser/dependency/querypb"
	"strings"
)

type Engine struct {
	storage *storage.Storage
	jobs    []*cron.Cron
}

func New() *Engine {
	return &Engine{
		storage: storage.New(),
		jobs:    make([]*cron.Cron, 0),
	}

}

func (e *Engine) Solve(query *model.Query) {
	source := query.Source
	switch query.QueryType {
	case model.SelectType:
		results := e.selectTokens(source, query)
		show(results)
	case model.CreateType:
		url := query.Url
		schedule := query.Schedule
		e.createSource(source, url, schedule)
	}
}

func show(results []*model.Token) {
	if len(results) == 0 {
		fmt.Println("no results!")
	} else {
		fmt.Println("Type  ||  Data")
		for _, token := range results {
			fmt.Printf("%s  ||  %s\n", token.Type, token.Data)
		}
	}
}

func (e *Engine) createSource(name string, url string, schedule string) {

	source := &model.Source{
		Name: name,
		Url:  url,
	}

	job := cron.New(cron.WithSeconds())
	collector := collector.New(source)

	e.jobs = append(e.jobs, job)
	e.storage.AddSource(source)

	collector.Collect()
	job.AddFunc(schedule, collector.Collect)
	job.Start()
}

func (e *Engine) selectTokens(sourceName string, query *model.Query) []*model.Token {
	return e.storage.SelectTokens(sourceName, query)
}

func (e *Engine) Parse(input string) (*model.Query, error) {
	var queryInput string
	var query *model.Query
	var bindVars = make(map[string]*querypb.BindVariable)
	var params = make(map[string]*model.ExpressionNode, 0)

	if strings.HasPrefix(input, "create source") {
		parts := common.NewString(input).
			ReplaceAll("create source", "").
			ReplaceAll(" from ", "&").
			ReplaceAll(" when ", "&").
			ReplaceAll(";", "").
			ReplaceAll("'", "").
			TrimSpace().
			Split("&").
			Values()

		name := parts[0]
		url := parts[1]
		schedule := parts[2]

		queryInput = fmt.Sprintf("insert into %s (url, schedule) values ('%s', '%s')", name, url, schedule)
	} else {
		queryInput = input
	}

	statement, err := sqlparser.Parse(queryInput)
	sqlparser.Normalize(statement, bindVars, "")

	if !common.IsError(err, "when parsing SQL input") {
		switch stmt := statement.(type) {
		case *sqlparser.Select:
			//var preField = ""
			//var operator Operator
			//var varType VarType
			//queryVars[QueryType] = SelectType

			whereBuffer := sqlparser.NewTrackedBuffer(nil)
			sourceBuffer := sqlparser.NewTrackedBuffer(nil)

			stmt.Where.Expr.Format(whereBuffer)
			stmt.From.Format(sourceBuffer)

			//conditions := strings.Split(whereBuffer.String(), " and ")
			tokens := common.NewString(whereBuffer.String()).
				ReplaceAll("not like", "not_like").
				Split(" ").
				Values()

			var expression = &model.ExpressionTree{
				Root: nil,
			}

			for _, value := range tokens {
				token := strings.Split(value, " ")[0]
				expression.Insert(token)

				if len(token) == 2 && token[0:1] == ":" {
					params[token[1:]] = expression.Root.Right
				}
			}

			for key, bindVar := range bindVars {
				if params[key].VarType == model.TokenType {
					tokenType := model.GetTokenType(string(bindVar.Value))
					params[key].Key = tokenType
				} else {
					params[key].Key = string(bindVar.Value)
				}

				if bindVar.Type == querypb.Type_VARBINARY {
					params[key].VarType = model.StringType
				}
			}

			query = &model.Query{
				Source:     sourceBuffer.String(),
				QueryType:  model.SelectType,
				Expression: expression,
			}

		case *sqlparser.Insert:
			tableBuffer := sqlparser.NewTrackedBuffer(nil)
			stmt.Table.Format(tableBuffer)
			if len(bindVars) == 2 {
				query = &model.Query{
					Source:    tableBuffer.String(),
					QueryType: model.CreateType,
					Url:       string(bindVars["1"].Value),
					Schedule:  string(bindVars["2"].Value),
				}
			} else {
				err = errors.New("wrong query variables size")
			}
		}
	}

	if query.QueryType == model.SelectType {
		//query.InOrderPrint()
	}

	return query, err
}
