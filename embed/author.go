package embed

// Author is the embedded author in a Discord message.
type Author struct {
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// AuthorBuilder creates authors for messages with rich embedded content.
type AuthorBuilder interface {
	Name(n string) AuthorBuilder
	URL(url string) AuthorBuilder
	IconURL(url string) AuthorBuilder
	Build() *Author
}

type authorBuilder struct {
	name    string
	url     string
	iconURL string
}

func (a *authorBuilder) Name(n string) AuthorBuilder {
	a.name = n
	return a
}

func (a *authorBuilder) URL(url string) AuthorBuilder {
	a.url = url
	return a
}

func (a *authorBuilder) IconURL(url string) AuthorBuilder {
	a.iconURL = url
	return a
}

func (a *authorBuilder) Build() *Author {
	return &Author{
		Name:    a.name,
		URL:     a.url,
		IconURL: a.iconURL,
	}
}

// NewAuthor returns a builder to create authors.
func NewAuthor() AuthorBuilder {
	return &authorBuilder{}
}
