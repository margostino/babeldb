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

func (s *Storage) SelectTokens(name string) []*Token {
	return s.sources[name].Tokens
}
