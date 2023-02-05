package model

import (
	"time"
)

type Source struct {
	Name       string
	Url        string
	Page       *Page
	LastUpdate time.Time
}
