package embed

import "time"

// Embed describes some rich content for a Discord message.
type Embed struct {
	Title       string     `json:"title,omitempty"`
	Type        string     `json:"type,omitempty"` // Type of embed (always "rich" for webhook embeds).
	Description string     `json:"description,omitempty"`
	URL         string     `json:"url,omitempty"`
	Timestamp   *time.Time `json:"timestamp,omitempty"`
	Color       int        `json:"color,omitempty"`
	Footer      *Footer    `json:"footer,omitempty"`
	Image       *Image     `json:"image,omitempty"`
	Thumbnail   *Thumbnail `json:"thumbnail,omitempty"`
	Video       *Video     `json:"video,omitempty"`
	Provider    *Provider  `json:"provider,omitempty"`
	Author      *Author    `json:"author,omitempty"`
	Fields      []Field    `json:"fields,omitempty"`
}

// Builder allows to add rich embedded content to Discord messages.
type Builder interface {
	Title(t string) Builder
	Description(t string) Builder
	URL(u string) Builder
	Timestamp(t time.Time) Builder
	Color(c int) Builder
	Footer(f *Footer) Builder
	Image(i *Image) Builder
	Thumbnail(i *Thumbnail) Builder
	Author(i *Author) Builder
	Fields(fields ...*Field) Builder
	Build() *Embed
}

type builder struct {
	title       string
	description string
	url         string
	timestamp   *time.Time
	color       int
	footer      *Footer
	image       *Image
	thumbnail   *Thumbnail
	author      *Author
	fields      []Field
}

func (e *builder) Title(t string) Builder {
	e.title = t
	return e
}

func (e *builder) Description(d string) Builder {
	e.description = d
	return e
}

func (e *builder) URL(u string) Builder {
	e.url = u
	return e
}

func (e *builder) Timestamp(t time.Time) Builder {
	e.timestamp = &t
	return e
}

func (e *builder) Color(c int) Builder {
	e.color = c
	return e
}

func (e *builder) Footer(f *Footer) Builder {
	e.footer = f
	return e
}

func (e *builder) Image(i *Image) Builder {
	e.image = i
	return e
}

func (e *builder) Thumbnail(t *Thumbnail) Builder {
	e.thumbnail = t
	return e
}

func (e *builder) Author(a *Author) Builder {
	e.author = a
	return e
}

func (e *builder) Fields(fields ...*Field) Builder {
	for i := 0; i < len(fields); i++ {
		e.fields = append(e.fields, *fields[i])
	}
	return e
}

func (e *builder) Build() *Embed {
	return &Embed{
		Title:       e.title,
		Description: e.description,
		URL:         e.url,
		Timestamp:   e.timestamp,
		Color:       e.color,
		Footer:      e.footer,
		Image:       e.image,
		Thumbnail:   e.thumbnail,
		Author:      e.author,
		Fields:      e.fields,
	}
}

// New returns a builder to create embeds.
func New() Builder {
	return &builder{}
}
