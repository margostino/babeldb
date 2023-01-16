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

func (s *Storage) SelectTokens(name string, query *model.Query) []*model.Token {
	results := make([]*model.Token, 0)
	if s.sources[name] != nil {
		for i, token := range s.sources[name].Tokens {
			if token.Data == "While the climate crisis has many factors that play a role in the exacerbation of the environment, there are some that warrant more attention than others. Here are […]" {
				fmt.Printf("%d\n", i)
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
					results = append(results, token)
				}
			}
		}
	} else {
		fmt.Println("source name not found!")
	}

	return results
}
