package model

import "github.com/margostino/babeldb/common"

type Operator int32
type Type int32
type QueryType int32

const (
	CreateSource string = "create source"
)

const (
	Wildcard                    string = "*"
	TypeField                   string = "type"
	DataField                   string = "data"
	TokenField                  string = "token"
	HrefField                   string = "href"
	SourceName                  string = "name"
	SourceUrl                   string = "url"
	SourceTotalSections         string = "total_sections"
	SourceLastUpdate            string = "last_update"
	SourceMeta                  string = "meta"
	SourceMetaTitle             string = "meta_title"
	SourceMetaDescription       string = "meta_description"
	SourceMetaTwitter           string = "meta_twitter"
	SourceMetaUrl               string = "meta_url"
	SourceMetaLocale            string = "meta_locale"
	SourcePageText              string = "text"
	SourcePageLinks             string = "links"
	SourcePageLink              string = "link"
	SourcePageSitemap           string = "sitemap"
	SourcePageSitemapUrl        string = "sitemap_url"
	SourcePageSitemapLastMod    string = "sitemap_lastmod"
	SourcePageSitemapChangeFreq string = "sitemap_changefreq"
)

var AttributeFields = common.NewStringSlice(HrefField)
var Fields = common.NewStringSlice(
	Wildcard,
	SourceName,
	SourceUrl,
	SourceTotalSections,
	SourceLastUpdate,
	SourceMeta,
	SourceMetaTitle,
	SourceMetaDescription,
	SourceMetaTwitter,
	SourceMetaLocale,
	SourceMetaUrl,
	SourcePageText,
	SourcePageLinks,
	SourcePageLink,
	TypeField,
	DataField,
	TokenField,
	HrefField,
	SourcePageSitemapUrl,
	SourcePageSitemapLastMod,
	SourcePageSitemapChangeFreq,
	SourcePageSitemap,
)

const (
	EqualOperator Operator = iota
	LikeOperator
	NotLikeOperator
	AndOperator
	OrOperator
	StringType Type = iota
	TokenType
	TextTokenType
	ErrorTokenType
	StartTagTokenType
	EndTagTokenType
	SelfClosingTagTokenType
	CommentTokenType
	DocTypeTokenType
	SelectType QueryType = iota
	CreateType
	ShowSources
)
