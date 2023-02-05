package collector

import (
	"encoding/xml"
	"fmt"
	"github.com/margostino/babeldb/common"
	"github.com/margostino/babeldb/model"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Extractor struct {
	url        string
	flags      *Flags
	attributes *model.Attributes
	section    *model.Section
	Page       *model.Page
}

type Flags struct {
	isTitleToken     bool
	isSectionToken   bool
	isParagraphToken bool
	isHeadlineToken  bool
	isSpanToken      bool
}

type SiteMap struct {
	Index []*model.SiteMapUrl `xml:"sitemap"`
	Urls  []*model.SiteMapUrl `xml:"url"`
}

func newExtractor(url string) *Extractor {
	return &Extractor{
		url:     url,
		flags:   &Flags{},
		Page:    model.NewPage(),
		section: newSection(),
	}
}

func newSection() *model.Section {
	return &model.Section{
		Links: make([]string, 0),
	}
}

func (e *Extractor) addSection(token *html.Token) {
	if !e.flags.isSectionToken && token.DataAtom == atom.Div && isStartToken(token) {
		e.flags.isSectionToken = true
	}
	if e.flags.isSectionToken && token.DataAtom == atom.Div && isEndToken(token) {
		e.flags.isSectionToken = false
		e.section.Text = sanitize(e.section.Text)
		if e.section.Text != "" || len(e.section.Links) > 0 {
			e.section.Text = sanitizeText(e.section.Text)
			e.Page.AddSection(e.section)
		}
		e.section = newSection()
	}
}

func (e *Extractor) addMeta(token *html.Token) {
	if e.flags.isTitleToken && token.Type == html.TextToken {
		e.flags.isTitleToken = false
		e.Page.Meta.Title = token.Data
	}
	if token.DataAtom == atom.Meta && isStartToken(token) {
		content := e.attributes.Get("content")
		if e.attributes.Get("name") == "description" {
			e.Page.Meta.Description = content
		}
		if e.attributes.Get("name") == "twitter:site" {
			e.Page.Meta.Twitter = content
		}
		if e.attributes.Get("property") == "og:title" && e.Page.Meta.Title == "" {
			e.Page.Meta.Title = content
		}
		if e.attributes.Get("property") == "og:url" {
			e.Page.Meta.Url = content
		}
		if e.attributes.Get("property") == "og:locale" {
			e.Page.Meta.Locale = content
		}
		if e.attributes.Get("property") == "og:description" {
			if e.Page.Meta.Description == "" {
				e.Page.Meta.Description = content
			}
		}
	}
}

func getSiteMapUrls(baseUrl string) []string {
	sitemapIndex := make([]string, 0)
	robotsUrl := fmt.Sprintf("%s/robots.txt", baseUrl)
	res, err := http.Get(robotsUrl)
	defer res.Body.Close()

	if !common.IsError(err, "when calling robots URL") {
		body, err := io.ReadAll(res.Body)
		if !common.IsError(err, "when reading robots call response") {
			text := string(body)
			tokens := strings.Split(text, "\n")
			for _, token := range tokens {
				if strings.HasPrefix(token, "Sitemap: ") {
					sitemapTokens := common.NewString(token).
						Split(" ").
						Values()
					if len(sitemapTokens) == 2 {
						sitemapIndex = append(sitemapIndex, sitemapTokens[1])
					}
				}
			}
		}
	}
	return sitemapIndex
}

// TODO: do it async
func (e *Extractor) addSitemap() {
	if e.Page.Meta != nil && e.Page.Meta.SiteMap == nil {
		urls := getSiteMapUrls(e.url)

		sites := make([]*model.SiteMapUrl, 0)
		appendSites(&sites, urls)
		e.Page.Meta.SiteMap = &model.SiteMap{
			Sites: sites,
		}
	}
}

func (e *Extractor) addText(token *html.Token) {
	if e.flags.isSectionToken && e.flags.isParagraphToken && token.Type == html.TextToken {
		e.section.Text += fmt.Sprintf("%s\n", token.Data)
		e.flags.isParagraphToken = false
	}
	if e.flags.isSectionToken && e.flags.isSpanToken && token.Type == html.TextToken {
		e.section.Text += fmt.Sprintf("%s\n", token.Data)
		e.flags.isSpanToken = false
	}
	if e.flags.isSectionToken && e.flags.isHeadlineToken && token.Type == html.TextToken {
		e.section.Text += fmt.Sprintf("%s\n", token.Data)
		e.flags.isHeadlineToken = false
	}
}

func (e *Extractor) addLink(token *html.Token) {
	url := e.url
	if e.flags.isSectionToken && token.DataAtom == atom.A {
		href := e.attributes.Get("href")
		hostname := getHostname(url)
		if href != "" && strings.Contains(href, hostname) {
			e.section.AddLink(href)
		}
	}
}

func (e *Extractor) mark(token *html.Token) {
	if token.DataAtom == atom.P {
		e.flags.isParagraphToken = true
	}
	if token.DataAtom == atom.Title {
		e.flags.isTitleToken = true
	}
	if token.DataAtom == atom.Span {
		e.flags.isSpanToken = true
	}
	if token.DataAtom == atom.H1 || token.DataAtom == atom.H2 || token.DataAtom == atom.H3 || token.DataAtom == atom.H4 || token.DataAtom == atom.H5 || token.DataAtom == atom.H6 {
		e.flags.isHeadlineToken = true
	}
}

func (e *Extractor) flag(token *html.Token) {
	e.attributes = model.NewAttributes(token.Attr)
	if isStartToken(token) {
		e.mark(token)
	}
}

func appendSites(sites *[]*model.SiteMapUrl, urls []string) {
	for _, url := range urls {
		res, err := http.Get(url)
		if !common.IsError(err, fmt.Sprintf("error when collecting data from %s", url)) {
			data, err := ioutil.ReadAll(res.Body)

			if !common.IsError(err, fmt.Sprintf("error when parsing response from %s", url)) {
				var sitemapResp SiteMap
				xml.Unmarshal(data, &sitemapResp)

				for _, index := range sitemapResp.Index {
					appendSites(sites, []string{index.Loc})
				}

				for _, sitemapUrl := range sitemapResp.Urls {
					*sites = append(*sites, sitemapUrl)
				}
			}
		}
	}
}

func sanitizeText(value string) string {
	partial := regExpLeadClose.ReplaceAllString(value, "")
	return regExpInsideClose.ReplaceAllString(partial, " ")
}
