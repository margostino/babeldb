package storage

import (
	"fmt"
	"github.com/margostino/babeldb/model"
	"strings"
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

func (s *Storage) SelectTokens(name string, query *model.Query) []*model.Token {
	results := make([]*model.Token, 0)
	if s.sources[name] != nil {
		var attribute *model.Attribute
		for _, token := range s.sources[name].Tokens {
			if len(token.Attributes) > 0 {
				newAttribute, _ := GetAttribute(token.Attributes, "href")
				if newAttribute != nil {
					attribute = newAttribute
				}
			}
			if strings.Contains(token.Data, "land use and household decisions in the province of Salta in the Gran Chaco") {
				println("")
			}
			match := query.Match(token)
			if match {
				var exists bool
				if query.Distinct {
					for _, result := range results {
						if result.Data == token.Data {
							exists = true
							break
						}
					}
				}
				if !exists {
					if attribute != nil {
						token.Attributes = append(token.Attributes, attribute)
						attribute = nil
					}
					results = append(results, token)
				}
			}
		}
	} else {
		fmt.Println("source name not found!")
	}

	return results
}
