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

func GetAttribute(attributes []*model.Attribute, key string) (*model.Attribute, bool) {
	for _, attribute := range attributes {
		if attribute.Key == key {
			return attribute, true
		}
	}
	return nil, false
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

func (s *Storage) Select(name string, query *model.Query) []*model.Section {
	sections := make([]*model.Section, 0)

	if s.sources[name] != nil {
		source := s.sources[name]

		for _, section := range source.Page.Sections {
			if query.Match(section) {
				sections = append(sections, section)
			}
		}
		//var attribute *model.Attribute
		//page := s.sources[name].Page
		//for _, section := range page.Sections {
		//	data := common.NewString(token.Data).TrimSpace().Value()
		//	if len(token.Attributes) > 0 {
		//		newAttribute, _ := GetAttribute(token.Attributes, "href")
		//		if newAttribute != nil {
		//			attribute = newAttribute
		//		}
		//	}
		//	if data != "" && data != "\t" {
		//		match := query.Match(token)
		//		if match {
		//			var exists bool
		//			if query.Distinct {
		//				for _, result := range results {
		//					if result.Data == token.Data {
		//						exists = true
		//						break
		//					}
		//				}
		//			}
		//			if !exists {
		//				if attribute != nil {
		//					token.Attributes = append(token.Attributes, attribute)
		//					attribute = nil
		//				}
		//				results = append(results, token)
		//			}
		//		}
		//	}
		//}
	} else {
		fmt.Println("source name not found!")
	}

	return sections
}
