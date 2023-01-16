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

func (s *Storage) SelectTokens(name string, query *model.Query) []*model.Token {
	results := make([]*model.Token, 0)
	if s.sources[name] != nil {
		var attribute *model.Attribute
		for i, token := range s.sources[name].Tokens {
			if len(token.Attributes) > 0 {
				newAttribute, _ := GetAttribute(token.Attributes, "href")
				if newAttribute != nil {
					attribute = newAttribute
				}
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

					if token.Data == "2022 has been a year of tremendous climate extremes. Humanity is learning the extent of the existential threats posed by climate change and ecological destruction the hard way. […]" {
						println(i)
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
