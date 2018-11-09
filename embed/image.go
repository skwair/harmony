package embed

// Image is an embedded image in a Discord Message.
type Image struct {
	URL      string `json:"url,omitempty"`       // Source url of image (only supports http(s) and attachments).
	ProxyURL string `json:"proxy_url,omitempty"` // A proxied url of the image.
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

// NewImage creates an embedded image from its source URL.
// Supported formats are JPEG, PNG, WebP and GIF.
func NewImage(url string) *Image {
	return &Image{URL: url}
}
