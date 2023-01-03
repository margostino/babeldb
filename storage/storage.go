package storage

import (
	"github.com/xwb1989/sqlparser/dependency/querypb"
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

func (s *Storage) SelectTokens(name string, conditions map[string]*querypb.BindVariable) []*Token {
	var tokenType html.TokenType
	results := make([]*Token, 0)

	if conditions["type"] != nil {
		switch value := string(conditions["type"].Value); value {
		case "text":
			tokenType = html.TextToken
			break
		case "error":
			tokenType = html.TextToken
			break
		case "start_tag":
			tokenType = html.StartTagToken
			break
		case "end_tag":
			tokenType = html.EndTagToken
			break
		case "self_closing_tag":
			tokenType = html.SelfClosingTagToken
			break
		case "comment":
			tokenType = html.CommentToken
			break
		case "doc_type":
			tokenType = html.DoctypeToken
			break
		default:
			// TODO
		}
	}

	for _, token := range s.sources[name].Tokens {
		if token.Type == tokenType {
			results = append(results, token)
		}
	}
	return results
}
