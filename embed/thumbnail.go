package embed

// Thumbnail is an embedded thumbnail in a Discord Message.
type Thumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

// NewThumbnail creates an embedded thumbnail from its source URL.
// Supported formats are JPEG, PNG, WebP and GIF.
func NewThumbnail(url string) *Thumbnail {
	return &Thumbnail{URL: url}
}
