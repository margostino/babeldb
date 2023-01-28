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
		sections := e.storage.Select(source, query)
		showPage(query.Fields, sections)
	case model.ShowSources:
		sources := e.storage.Show()
		showSources(sources)
	case model.CreateType:
		url := query.Url
		schedule := query.Schedule
		go e.createSource(source, url, schedule)
	}
}

func showSources(sources []*model.Source) {
	if len(sources) == 0 {
		fmt.Println("no results!")
	} else {
		// TODO: pretty format
		fmt.Println()
		fmt.Println("---------------------------")
		for _, source := range sources {
			fmt.Printf("Name:  %s\n", source.Name)
			fmt.Printf("URL:  %s\n", source.Url)
			fmt.Printf("Last update: %s\n", source.LastUpdate)
			fmt.Printf("Title: %s\n", source.Page.Meta.Title)
			fmt.Printf("Description: %s\n", source.Page.Meta.Description)
			fmt.Printf("Twitter: %s\n", source.Page.Meta.Twitter)
			fmt.Printf("Locale: %s\n", source.Page.Meta.Locale)
			fmt.Printf("Total sections: %d\n", len(source.Page.Sections))
			fmt.Println("---------------------------")
		}
		fmt.Println()
	}
}

func showPage(fields *common.StringSlice, sections []*model.Section) {
	if len(sections) == 0 {
		fmt.Println("no results!")
	} else {
		// TODO: pretty format
		fmt.Println()
		fmt.Println("\n---------------------------")
		for _, section := range sections {
			if fields.Contains(model.SourcePageText) || fields.Contains(model.Wildcard) {
				fmt.Printf("Text:  %s\n", section.Text)
			}
			if fields.Contains(model.SourcePageLinks) || fields.Contains(model.Wildcard) {
				fmt.Printf("Links:  %s\n", section.Links)
			}
			if fields.Contains(model.SourcePageLink) || fields.Contains(model.Wildcard) {
				fmt.Printf("Last update: %s\n", section.Links)
			}
			fmt.Println("---------------------------")
		}
		fmt.Println()
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

func (e *Engine) Parse(input string) (*model.Query, error) {
	var queryInput string
	var query *model.Query
	var bindVars = make(map[string]*querypb.BindVariable)
	var params = make(map[string]*model.ExpressionNode, 0)

	if input == "show sources" {
		query = &model.Query{
			Source:    model.Wildcard,
			QueryType: model.ShowSources,
		}
		return query, nil
	}

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
			var whereClauses string
			whereBuffer := sqlparser.NewTrackedBuffer(nil)
			sourceBuffer := sqlparser.NewTrackedBuffer(nil)
			selectBuffer := sqlparser.NewTrackedBuffer(nil)

			stmt.SelectExprs.Format(selectBuffer)
			stmt.From.Format(sourceBuffer)

			if stmt.Where.Expr != nil {
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

			query = &model.Query{
				Source:     sourceBuffer.String(),
				QueryType:  model.SelectType,
				Fields:     common.NewStringSlice(fields...),
				Distinct:   strings.HasPrefix(queryInput, "select distinct"),
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

	//if query.QueryType == model.SelectType {
	//	query.InOrderPrint()
	//}

	return query, err
}
