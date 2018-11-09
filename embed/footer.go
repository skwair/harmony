package embed

// Footer is a embedded footer in a Discord message.
type Footer struct {
	Text    string `json:"text,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

// FooterBuilder creates footers for messages with rich embedded content.
type FooterBuilder interface {
	Text(t string) FooterBuilder
	Icon(url string) FooterBuilder
	Build() *Footer
}

type footerBuilder struct {
	text         string
	iconURL      string
	proxyIconURL string
}

func (f *footerBuilder) Text(t string) FooterBuilder {
	f.text = t
	return f
}

func (f *footerBuilder) Icon(url string) FooterBuilder {
	f.iconURL = url
	return f
}

func (f *footerBuilder) Build() *Footer {
	return &Footer{
		Text:    f.text,
		IconURL: f.iconURL,
	}
}

// NewFooter returns a builder to create footers.
func NewFooter() FooterBuilder {
	return &footerBuilder{}
}
