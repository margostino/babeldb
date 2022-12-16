package engine

import (
	"github.com/margostino/babeldb/collector"
	"github.com/margostino/babeldb/storage"
	"github.com/robfig/cron/v3"
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

func (e *Engine) AddSource(name string, url string, cronExp string) {
	source := &storage.Source{
		Name: name,
		Url:  url,
	}

	job := cron.New(cron.WithSeconds())
	collector := collector.New(source)

	e.jobs = append(e.jobs, job)
	e.storage.AddSource(source)

	job.AddFunc(cronExp, collector.Collect)
	job.Start()
}
