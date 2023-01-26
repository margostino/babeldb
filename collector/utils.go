package collector

import (
	"golang.org/x/net/html"
	"log"
	"net/url"
	"strings"
)

func isStartToken(token *html.Token) bool {
	return token.Type == html.StartTagToken || token.Type == html.SelfClosingTagToken
}

func isEndToken(token *html.Token) bool {
	return token.Type == html.EndTagToken || token.Type == html.SelfClosingTagToken
}

func getHostname(href string) string {
	url, err := url.Parse(href)
	if err != nil {
		log.Fatal(err)
	}
	hostname := strings.TrimPrefix(url.Hostname(), "www.")
	return hostname
}
