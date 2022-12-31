package engine

import (
	"github.com/margostino/babeldb/collector"
	"github.com/margostino/babeldb/common"
	"github.com/margostino/babeldb/storage"
	"github.com/robfig/cron/v3"
	"github.com/xwb1989/sqlparser"
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

func (e *Engine) selectTokens(sourceName string) {
	e.storage.SelectTokens(sourceName)
}

func (e *Engine) Parse(input string) (*sqlparser.Statement, error) {
	var statement *sqlparser.Statement

	if strings.HasPrefix(input, "create source") {
		parts := common.NewString(input).
			ReplaceAll("create source", "").
			ReplaceAll(" from ", "&").
			ReplaceAll(" when ", "&").
			ReplaceAll(";", "").
			ReplaceAll("'", "").
			Split("&").
			Values()

		e.createSource(parts[0], parts[1], parts[2])
	} else {
		statement, err := sqlparser.Parse(input)

		if !common.IsError(err, "when parsing SQL input") {
			switch stmt := statement.(type) {
			case *sqlparser.Select:
				println("do select")
				e.selectTokens("earth")
				_ = stmt
			case *sqlparser.Insert:
			}
		}
		return &statement, err
	}

	return statement, nil

}

//func selectData(engine *engine.Engine, params map[Param]interface{}) {
//	name := params[sourceName].(string)
//	url := params[sourceUrl].(string)
//	cronExp := params[schedule].(string)
//	engine.AddSource(name, url, cronExp)
//}
