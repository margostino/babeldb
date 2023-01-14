package collector

import (
	"fmt"
	"github.com/margostino/babeldb/common"
	"github.com/margostino/babeldb/model"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"strings"
)

type Collector struct {
	source *model.Source
}

func New(source *model.Source) *Collector {
	return &Collector{
		source: source,
	}
}

func (c *Collector) Collect() {
	url := c.source.Url
	res, err := http.Get(url)

	if !common.IsError(err, fmt.Sprintf("error when collecting data from %s", url)) {
		text, err := ioutil.ReadAll(res.Body)

		if !common.IsError(err, fmt.Sprintf("error when parsing response from %s", url)) {
			tokens := parse(string(text))
			c.source.Tokens = tokens
		}
	}

}

func parse(text string) []*model.Token {

	var tokens = make([]*model.Token, 0)
	tkn := html.NewTokenizer(strings.NewReader(text))

	for {
		tokenType := tkn.Next()
		currentToken := tkn.Token()

		if isValidTokenType(tokenType) && isValidData(currentToken) {
			attrs := make([]*model.Attributes, 0)
			for _, attr := range currentToken.Attr {
				att := &model.Attributes{
					Key:   attr.Key,
					Value: attr.Val,
				}
				attrs = append(attrs, att)
			}
			token := &model.Token{
				Type:       tokenType,
				Data:       currentToken.Data,
				Attributes: attrs,
			}
			tokens = append(tokens, token)
		}

		if tokenType == html.ErrorToken {
			return tokens
		}

	}
}

func isValidTokenType(tokenType html.TokenType) bool {
	return tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken || tokenType == html.TextToken
}

func isValidData(token html.Token) bool {
	return !strings.Contains(token.Data, "\n")
}
