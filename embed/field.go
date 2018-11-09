package embed

// Field is an embedded field in a Discord message.
type Field struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

// FieldBuilder creates fields for messages with rich embedded content.
type FieldBuilder interface {
	Name(url string) FieldBuilder
	Value(v string) FieldBuilder
	Inline(yes bool) FieldBuilder
	Build() *Field
}

type fieldBuilder struct {
	name   string
	value  string
	inline bool
}

func (f *fieldBuilder) Name(n string) FieldBuilder {
	f.name = n
	return f
}

func (f *fieldBuilder) Value(v string) FieldBuilder {
	f.value = v
	return f
}

func (f *fieldBuilder) Inline(yes bool) FieldBuilder {
	f.inline = yes
	return f
}

func (f *fieldBuilder) Build() *Field {
	return &Field{
		Name:   f.name,
		Value:  f.value,
		Inline: f.inline,
	}
}

// NewField returns a builder to create fields.
func NewField() FieldBuilder {
	return &fieldBuilder{}
}
