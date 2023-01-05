package storage

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
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

func (s *Storage) SelectTokens(name string, conditions map[string]string) []*Token {
	var tokenType html.TokenType
	results := make([]*Token, 0)

	if conditions["type"] != "" {
		switch value := conditions["type"]; value {
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
	} else {
		tokenType = 100
	}

	if s.sources[name] != nil {
		for _, token := range s.sources[name].Tokens {
			lowerData := strings.ToLower(token.Data)
			tokenTypeMatch := (tokenType == 100) || (tokenType != 100 && token.Type == tokenType)
			dataMatch := conditions["data"] == "" || strings.Contains(lowerData, conditions["data"])
			if tokenTypeMatch && dataMatch {
				if len(token.Attributes) > 0 {
					println("f")
				}
				results = append(results, token)
			}
		}
	} else {
		fmt.Println("source name not found!")
	}

	return results
}
