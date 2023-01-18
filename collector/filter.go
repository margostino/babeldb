package collector

import (
	"github.com/margostino/babeldb/common"
	"golang.org/x/net/html"
	"strings"
)

// TODO: redefine, not hardcode rules, build the dom(?) and smartly delete components like styles (?)
func shouldFilter(token *html.Token) bool {
	data := common.NewString(token.Data).
		ReplaceAll("\n", "").
		ReplaceAll("\t", "").
		TrimSpace().
		Value()

	if strings.Contains(data, "Our best stories direct to your inbox every fortnight.") {
		println("")
	}

	//for _, e := range token.Attr {
	//	if strings.Contains(e.Val, "https://www.linkedin.com/shareArticle?mini") {
	//		println("")
	//	}
	//}

	return data == "" ||
		data == "\n" ||
		data == "\t" ||
		data == "html" ||
		data == "head" ||
		data == "meta" ||
		data == "title" ||
		data == "style" ||
		data == "link" ||
		data == "script" ||
		data == "noscript" ||
		data == "body" ||
		data == "header" ||
		data == "div" ||
		data == "svg" ||
		data == "path" ||
		data == "span" ||
		strings.Contains(data, "gtag(") ||
		strings.Contains(data, "!function(") ||
		strings.Contains(data, "display: block") ||
		strings.Contains(data, "<![CDATA[") ||
		strings.Contains(data, "@font-face") ||
		strings.Contains(data, "<style") ||
		strings.Contains(data, "<link") ||
		strings.Contains(data, "<img") ||
		strings.Contains(data, "jQuery") ||
		strings.Contains(data, "window.attachEvent") ||
		strings.Contains(data, "background-color") ||
		strings.Contains(data, "padding-bottom") ||
		strings.Contains(data, ":after {")
}
