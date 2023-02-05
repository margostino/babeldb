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

const UrlSeparator = "@@@"

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

func (e *Engine) createSource(source *model.Source, schedule string) {

	job := cron.New(cron.WithSeconds())
	collector := collector.New(source)

	e.jobs = append(e.jobs, job)
	e.storage.AddSource(source)

	collector.Collect()
	job.AddFunc(schedule, collector.Collect)
	job.Start()
}

func (e *Engine) Parse(input string) (*model.Query, error) {
	var queryInput string
	var query = model.NewQuery()
	var bindVars = make(map[string]*querypb.BindVariable)
	var params = make(map[string]*model.ExpressionNode, 0)

	if strings.ToLower(input) == "show sources" {
		query.Source = model.Wildcard
		query.QueryType = model.ShowSources
		return query, nil
	}

	if shouldCreateSource(input) {
		parts := common.NewString(input).
			ReplaceAll("'", "").
			TrimSpace().
			Split(" ").
			Values()

		name := parts[2]
		url := parts[4]
		schedule := fmt.Sprintf("%s %s %s %s %s", parts[6], parts[7], parts[8], parts[9], parts[10])

		queryInput = fmt.Sprintf("insert into %s (url, schedule) values ('%s', '%s')", name, url, schedule)
	} else {
		queryInput = input
	}

	statement, err := sqlparser.Parse(queryInput)
	sqlparser.Normalize(statement, bindVars, "")

	if !common.IsError(err, "when parsing SQL input") {
		switch stmt := statement.(type) {
		case *sqlparser.Select:
			var whereClauses string
			whereBuffer := sqlparser.NewTrackedBuffer(nil)
			sourceBuffer := sqlparser.NewTrackedBuffer(nil)
			selectBuffer := sqlparser.NewTrackedBuffer(nil)

			stmt.SelectExprs.Format(selectBuffer)
			stmt.From.Format(sourceBuffer)

			if stmt.Where != nil && stmt.Where.Expr != nil {
				stmt.Where.Expr.Format(whereBuffer)
				whereClauses = whereBuffer.String()
			}

			fields := common.NewString(selectBuffer.String()).
				ReplaceAll(" ", "").
				ReplaceAll("`", "").
				Split(",").
				Values()

			// TODO: improve and move validations
			for _, field := range fields {
				if !model.Fields.Contains(field) {
					return nil, errors.New("invalid fields")
				}
			}

			//conditions := strings.Split(whereBuffer.String(), " and ")
			var tokens []string
			if whereClauses != "" {
				tokens = common.NewString(whereClauses).
					ReplaceAll("not like", "not_like").
					ReplaceAll("`", "").
					Split(" ").
					Values()
			}

			var expression = &model.ExpressionTree{
				Root: nil,
			}

			for _, value := range tokens {
				token := strings.Split(value, " ")[0]
				expression.Insert(token)

				if len(token) == 2 && token[0:1] == ":" {
					params[token[1:]] = expression.GetParamNode(expression.Root)
				}
			}

			for key, bindVar := range bindVars {
				params[key].Key.Value = string(bindVar.Value)
				if bindVar.Type == querypb.Type_VARBINARY {
					params[key].Key.Type = model.StringType
				}
			}

			query.Source = sourceBuffer.String()
			query.QueryType = model.SelectType
			query.Fields = common.NewStringSlice(fields...)
			query.Distinct = strings.HasPrefix(queryInput, "select distinct")
			query.Expression = expression

		case *sqlparser.Insert:
			tableBuffer := sqlparser.NewTrackedBuffer(nil)
			stmt.Table.Format(tableBuffer)
			if len(bindVars) == 2 {
				query.Source = tableBuffer.String()
				query.QueryType = model.CreateType
				query.Url = string(bindVars["1"].Value)
				query.Schedule = string(bindVars["2"].Value)
			} else {
				err = errors.New("wrong query variables size")
			}
		}
	}

	//if query.QueryType == model.SelectType {
	//	query.InOrderPrint()
	//}

	return query, err
}

func (e *Engine) Solve(query *model.Query) *model.QueryResults {
	var results = &model.QueryResults{
		Sources: make([]*model.Source, 0),
		Page:    model.NewPage(),
	}

	sourceName := query.Source
	switch query.QueryType {
	case model.SelectType:
		meta, sections := e.storage.Select(sourceName, query)
		results.Page = &model.Page{
			Meta:     meta,
			Sections: sections,
		}
	case model.ShowSources:
		sources := e.storage.Show()
		results.Sources = sources
	case model.CreateType:
		url := query.Url
		schedule := query.Schedule
		source := &model.Source{
			Name: sourceName,
			Url:  url,
			Page: model.NewPage(),
		}
		results.Sources = append(results.Sources, source)
		go e.createSource(source, schedule)
	}
	return results
}
