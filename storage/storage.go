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

func (s *Storage) Select(name string, query *model.Query) []*model.Source {
	results := make([]*model.Source, 0)

	if name == model.Sources {
		for _, source := range s.sources {
			results = append(results, source)
		}
	} else if s.sources[name] != nil {

	} else {
		fmt.Println("source name not found!")
	}

	if s.sources[name] != nil {
		//var attribute *model.Attribute
		//for _, token := range s.sources[name].Tokens {
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
	}

	return results
}
