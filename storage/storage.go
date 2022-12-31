package storage

import (
	"golang.org/x/net/html"
)

type Storage struct {
	sources map[string]*Source
}

type Attributes struct {
	Key   string
	Value string
}

type Token struct {
	Type       html.TokenType
	Data       string
	Attributes []*Attributes
}

type Source struct {
	Name   string
	Url    string
	Tokens []*Token
}

func New() *Storage {
	return &Storage{
		sources: make(map[string]*Source),
	}
}

func (s *Storage) AddSource(source *Source) {
	s.sources[source.Name] = source
}

func (s *Storage) SelectTokens(name string, conditions map[string]interface{}) []*Token {
	results := make([]*Token, 0)
	tokenType := conditions[""].(html.TokenType)
	for _, token := range s.sources[name].Tokens {
		if token.Type == tokenType {
			results = append(results, token)
		}
	}
	return results
}
