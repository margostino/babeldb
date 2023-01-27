package collector

import (
	"fmt"
	"github.com/margostino/babeldb/common"
	"github.com/margostino/babeldb/model"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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
			c.parse(string(text))
		}
	}

}

func (c *Collector) parse(text string) {
	extractor := newExtractor()
	tokenizer := html.NewTokenizer(strings.NewReader(text))

	for {
		_ = tokenizer.Next()
		token := tokenizer.Token()

		extractor.flag(&token)
		extractor.addMeta(&token)
		extractor.addText(&token)
		extractor.addSection(&token)
		extractor.addLink(c.source.Url, &token)

		if token.Type == html.ErrorToken {
			c.source.Page = extractor.Page
			c.source.LastUpdate = time.Now().UTC()
			return
		}

	}
}
