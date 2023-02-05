package model

import "time"

type Meta struct {
	Title       string
	Url         string
	Description string
	Twitter     string
	Locale      string
	SiteMap     *SiteMap
}

type SiteMapUrl struct {
	Loc        string    `xml:"loc"`
	Lastmod    time.Time `xml:"lastmod"`
	ChangeFreq string    `xml:"changefreq"`
	Priority   float32   `xml:"priority"`
}

type SiteMap struct {
	Sites []*SiteMapUrl `xml:"url"`
}

type Section struct {
	Links []string
	Text  string
}

type Page struct {
	Meta     *Meta
	Sections []*Section
}

func NewPage() *Page {
	return &Page{
		Meta:     NewMeta(),
		Sections: make([]*Section, 0),
	}
}

func NewMeta() *Meta {
	return &Meta{
		SiteMap: NewSitemap(),
	}
}

func NewSitemap() *SiteMap {
	return &SiteMap{
		Sites: make([]*SiteMapUrl, 0),
	}
}

func NewSection() *Section {
	return &Section{
		Links: make([]string, 0),
	}
}

func (s *Section) AddLink(link string) {
	s.Links = append(s.Links, link)
}

func (p *Page) AddSection(section *Section) {
	p.Sections = append(p.Sections, section)
}

func (m *Meta) IsCompleted() bool {
	return m.Title != "" &&
		m.Locale != "" &&
		m.Twitter != "" &&
		m.Url != "" &&
		m.Description != ""
}

func (m *Meta) AddSitemapUrl(url *SiteMapUrl) {
	m.SiteMap.Sites = append(m.SiteMap.Sites, url)
}

func (m *Meta) AddSitemap(sitemap *SiteMap) {
	m.SiteMap = sitemap
}
