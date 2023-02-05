package model

import (
	"github.com/margostino/babeldb/common"
	"golang.org/x/net/html"
	"time"
)

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

type Attribute struct {
	Key   string
	Value string
}

type Token struct {
	Type       html.TokenType
	Data       string
	Attributes []*Attribute
}

type Source struct {
	Name       string
	Url        string
	Page       *Page
	LastUpdate time.Time
}

type ExpressionTree struct {
	Root *ExpressionNode
}

type ExpressionNodeKey struct {
	Value    string
	Type     Type
	Operator Operator
}

type ExpressionNode struct {
	Key   *ExpressionNodeKey
	Left  *ExpressionNode
	Right *ExpressionNode
}

type Query struct {
	Source     string
	Url        string
	Schedule   string
	Fields     *common.StringSlice
	QueryType  QueryType
	Distinct   bool
	Expression *ExpressionTree
}

type Attributes struct {
	attributes []html.Attribute
}

func NewAttributes(attributes []html.Attribute) *Attributes {
	return &Attributes{
		attributes: attributes,
	}
}

func NewSource() *Source {
	return &Source{
		Page: NewPage(),
	}
}

func (q *Query) InOrderPrint() {
	q.Expression.InOrderPrint(q.Expression.Root)
}

func (q *Query) Match(section *Section) bool {
	return q.Expression.Root == nil || (q.Expression.Root != nil && q.Expression.Match(q.Expression.Root, section))
}

func (s *Attributes) Get(key string) string {
	for _, attribute := range s.attributes {
		if attribute.Key == key {
			return attribute.Val
		}
	}
	return ""
}
