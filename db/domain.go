package db

import "golang.org/x/net/html"

type Attributes struct {
	Key   string
	Value string
}

type Token struct {
	Type       html.TokenType
	Data       string
	Attributes []*Attributes
}

type Source struct {
	Name     string
	Url      string
	Schedule string
	Tokens   []*Token
}
