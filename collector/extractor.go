package collector

import (
	"fmt"
	"github.com/margostino/babeldb/model"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"log"
	"net/url"
	"strings"
)

type Extractor struct {
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

func newExtractor() *Extractor {
	meta := &model.Meta{}
	return &Extractor{
		flags: &Flags{},
		Page: &model.Page{
			Meta:     meta,
			Sections: make([]*model.Section, 0),
		},
		section: &model.Section{
			Links: make([]string, 0),
		},
	}
}

func (e *Extractor) addSection(token *html.Token) {
	if !e.flags.isSectionToken && token.DataAtom == atom.Div && token.Type == html.StartTagToken {
		e.flags.isSectionToken = true
	}
	if e.flags.isSectionToken && token.DataAtom == atom.Div && token.Type == html.EndTagToken {
		e.flags.isSectionToken = false
		e.section.Text = sanitize(e.section.Text)
		e.Page.AddSection(e.section)
	}
}

func (e *Extractor) addMeta(token *html.Token) {
	e.Page.Meta.Title = e.getText(token)
	if token.DataAtom == atom.Meta && token.Type == html.StartTagToken {
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

func (e *Extractor) getText(token *html.Token) string {
	if e.flags.isTitleToken && token.Type == html.TextToken {
		e.flags.isTitleToken = false
		return token.Data
	}
	return ""
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

func (e *Extractor) addLink(url string, token *html.Token) {
	if e.flags.isSectionToken && token.DataAtom == atom.A {
		href := e.attributes.Get("href")
		hostname := getHostname(url)
		if href != "" && strings.Contains(href, hostname) {
			if strings.Contains(href, "members.earth.org") {
				println()
			}
			e.section.AddLink(href)
		}
	}
}

func (e *Extractor) mark(token *html.Token) {
	if token.DataAtom == atom.P && token.Type == html.StartTagToken {
		e.flags.isParagraphToken = true
	}
	if token.DataAtom == atom.Span && token.Type == html.StartTagToken {
		e.flags.isSpanToken = true
	}
	if (token.DataAtom == atom.H1 || token.DataAtom == atom.H2 || token.DataAtom == atom.H3 || token.DataAtom == atom.H4 || token.DataAtom == atom.H5) && token.Type == html.StartTagToken {
		e.flags.isHeadlineToken = true
	}
}

func (e *Extractor) start(token *html.Token) {
	e.attributes = model.NewAttributes(token.Attr)
	if token.Type == html.StartTagToken {
		e.mark(token)
	}
}

func getHostname(href string) string {
	url, err := url.Parse(href)
	if err != nil {
		log.Fatal(err)
	}
	hostname := strings.TrimPrefix(url.Hostname(), "www.")
	return hostname
}
