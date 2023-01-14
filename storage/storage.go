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
	//var tokenType html.TokenType
	results := make([]*model.Token, 0)

	//if query.Vars conditions["type"] != "" {
	//	switch value := conditions["type"]; value {
	//	case "text":
	//		tokenType = html.TextToken
	//		break
	//	case "error":
	//		tokenType = html.TextToken
	//		break
	//	case "start_tag":
	//		tokenType = html.StartTagToken
	//		break
	//	case "end_tag":
	//		tokenType = html.EndTagToken
	//		break
	//	case "self_closing_tag":
	//		tokenType = html.SelfClosingTagToken
	//		break
	//	case "comment":
	//		tokenType = html.CommentToken
	//		break
	//	case "doc_type":
	//		tokenType = html.DoctypeToken
	//		break
	//	default:
	//		// TODO
	//	}
	//} else {
	//	tokenType = 100
	//}

	if s.sources[name] != nil {
		for _, token := range s.sources[name].Tokens {
			match := query.Match(token)
			if match {
				results = append(results, token)
			}
			//tokenTypeMatch := (tokenType == 100) || (tokenType != 100 && token.Type == tokenType)
			//dataMatch := conditions["data"] == "" || strings.Contains(lowerData, conditions["data"])
			//if tokenTypeMatch && dataMatch {
			//	if len(token.Attributes) > 0 {
			//		println("f")
			//	}
			//	results = append(results, token)
			//}
		}
	} else {
		fmt.Println("source name not found!")
	}

	return results
}
