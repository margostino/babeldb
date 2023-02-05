package model

import "golang.org/x/net/html"

type Attributes struct {
	attributes []html.Attribute
}

func NewAttributes(attributes []html.Attribute) *Attributes {
	return &Attributes{
		attributes: attributes,
	}
}

func (s *Attributes) Get(key string) string {
	for _, attribute := range s.attributes {
		if attribute.Key == key {
			return attribute.Val
		}
	}
	return ""
}
