package model

type Meta struct {
	Title       string
	Url         string
	Description string
	Twitter     string
	Locale      string
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
		Meta:     &Meta{},
		Sections: make([]*Section, 0),
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
