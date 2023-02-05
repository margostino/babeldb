package storage

import (
	"fmt"
	"github.com/margostino/babeldb/model"
)

type Storage struct {
	sources map[string]*model.Source
}

func New() *Storage {
	return &Storage{
		sources: make(map[string]*model.Source),
	}
}

func (s *Storage) AddSource(source *model.Source) {
	s.sources[source.Name] = source
}

func (s *Storage) SelectSources(query *model.Query) map[string]*model.Source {
	// TODO
	return s.sources
}

func (s *Storage) Show() []*model.Source {
	sources := make([]*model.Source, 0)
	for _, source := range s.sources {
		sources = append(sources, source)
	}
	return sources
}

func (s *Storage) Select(name string, query *model.Query) (*model.Meta, []*model.Section) {
	sections := make([]*model.Section, 0)
	meta := &model.Meta{}

	if s.sources[name] != nil {
		source := s.sources[name]
		meta = source.Page.Meta
		for _, section := range source.Page.Sections {
			if query.Match(section) {
				sections = append(sections, section)
			}
		}

	} else {
		fmt.Println("source name not found!")
	}

	return meta, sections
}
