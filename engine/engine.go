package engine

import (
	"errors"
	"fmt"
	"github.com/margostino/babeldb/collector"
	"github.com/margostino/babeldb/common"
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

func (e *Engine) Solve(query map[interface{}]interface{}) {
	source := query[Source].(string)
	switch query[QueryType] {
	case SelectType:
		conditions := make(map[string]string)
		for key, value := range query {
			if key != QueryType && key != Source {
				conditions[key.(string)] = value.(string)
			}
		}
		results := e.selectTokens(source, conditions)
		show(results)
	case CreateType:
		url := query[Url].(string)
		schedule := query[Schedule].(string)
		e.createSource(source, url, schedule)
	}
}

func show(results []*storage.Token) {
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

	source := &storage.Source{
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

func (e *Engine) selectTokens(sourceName string, conditions map[string]string) []*storage.Token {
	return e.storage.SelectTokens(sourceName, conditions)
}

func (e *Engine) getConditions(whereConditions *map[string]*querypb.BindVariable) map[string]interface{} {
	conditions := make(map[string]interface{})
	return conditions
}

func (e *Engine) Parse(input string) (map[interface{}]interface{}, error) {
	var query string
	var bindVars = make(map[string]*querypb.BindVariable)
	var queryVars = make(map[interface{}]interface{}, 0)

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

		query = fmt.Sprintf("insert into %s (url, schedule) values ('%s', '%s')", name, url, schedule)
	} else {
		query = input
	}

	statement, err := sqlparser.Parse(query)
	sqlparser.Normalize(statement, bindVars, "")

	if !common.IsError(err, "when parsing SQL input") {
		switch stmt := statement.(type) {
		case *sqlparser.Select:
			var preField = ""
			queryVars[QueryType] = SelectType

			whereBuffer := sqlparser.NewTrackedBuffer(nil)
			sourceBuffer := sqlparser.NewTrackedBuffer(nil)

			stmt.Where.Expr.Format(whereBuffer)
			stmt.From.Format(sourceBuffer)

			fields := whereBuffer.String()
			queryVars[Source] = sourceBuffer.String()

			for id, bindVar := range bindVars {
				fields = common.NewString(fields).
					ReplaceAll(fmt.Sprintf(" = :%s", id), "").
					ReplaceAll(" and ", " ").
					ReplaceAll(preField, "").
					Value()
				field := common.NewString(fields).
					TrimSpace().
					Split(" ").
					Values()[0]
				preField = field
				queryVars[field] = string(bindVar.Value)
			}

			//results := e.selectTokens("earth", queryVars)
			//println(results)

		case *sqlparser.Insert:
			tableBuffer := sqlparser.NewTrackedBuffer(nil)
			stmt.Table.Format(tableBuffer)
			if len(bindVars) == 2 {
				queryVars[QueryType] = CreateType
				queryVars[Source] = tableBuffer.String()
				queryVars[Url] = string(bindVars["1"].Value)
				queryVars[Schedule] = string(bindVars["2"].Value)
			} else {
				err = errors.New("wrong query variables size")
			}
		}
	}

	return queryVars, err

}

//func selectData(engine *engine.Engine, params map[Param]interface{}) {
//	name := params[sourceName].(string)
//	url := params[sourceUrl].(string)
//	cronExp := params[schedule].(string)
//	engine.AddSource(name, url, cronExp)
//}
