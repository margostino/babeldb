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
	source     *model.Source
	httpClient *http.Client
}

func New(source *model.Source) *Collector {
	return &Collector{
		source:     source,
		httpClient: &http.Client{},
	}
}

func (c *Collector) Collect() {
	url := c.source.Url
	req, _ := http.NewRequest("GET", url, nil)
	res, err := c.httpClient.Do(req)

	if !common.IsError(err, fmt.Sprintf("error when collecting data from %s", url)) {
		text, err := ioutil.ReadAll(res.Body)

		if !common.IsError(err, fmt.Sprintf("error when parsing response from %s", url)) {
			c.parse(string(text))
		}
	}

}

func (c *Collector) parse(text string) {
	extractor := newExtractor(c.source.Url)
	tokenizer := html.NewTokenizer(strings.NewReader(text))

	go extractor.addSitemap()
	for {
		_ = tokenizer.Next()
		token := tokenizer.Token()

		extractor.flag(&token)
		extractor.addLink(&token)
		extractor.addMeta(&token)
		extractor.addText(&token)
		extractor.addSection(&token)

		if token.Type == html.ErrorToken {
			c.source.Page = extractor.Page
			c.source.LastUpdate = time.Now().UTC()
			return
		}

	}
}
