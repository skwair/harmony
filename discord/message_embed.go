package discord

import "time"

// MessageEmbed describes some rich content for a Discord message.
type MessageEmbed struct {
	Title       string                 `json:"title,omitempty"`
	Type        string                 `json:"type,omitempty"` // Type of embed (always "rich" for webhook embeds).
	Description string                 `json:"description,omitempty"`
	URL         string                 `json:"url,omitempty"`
	Timestamp   *time.Time              `json:"timestamp,omitempty"`
	Color       int                    `json:"color,omitempty"`
	Footer      *MessageEmbedFooter    `json:"footer,omitempty"`
	Image       *MessageEmbedImage     `json:"image,omitempty"`
	Thumbnail   *MessageEmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *MessageEmbedVideo     `json:"video,omitempty"`
	Provider    *MessageEmbedProvider  `json:"provider,omitempty"`
	Author      *MessageEmbedAuthor    `json:"author,omitempty"`
	Fields      []MessageEmbedField    `json:"fields,omitempty"`
}

// MessageEmbedFooter is a embedded footer in a Discord message.
type MessageEmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// MessageEmbedImage is an embedded image in a Discord Message.
type MessageEmbedImage struct {
	URL      string `json:"url,omitempty"`       // Source url of image (only supports http(s) and attachments).
	ProxyURL string `json:"proxy_url,omitempty"` // A proxied url of the image.
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

// MessageEmbedThumbnail is an embedded thumbnail in a Discord Message.
type MessageEmbedThumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

// MessageEmbedVideo is an embedded video in a Discord Message.
type MessageEmbedVideo struct {
	URL    string `json:"url,omitempty"`
	Height int    `json:"height,omitempty"`
	Width  int    `json:"width,omitempty"`
}

type MessageEmbedProvider struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// MessageEmbedAuthor is the embedded author in a Discord message.
type MessageEmbedAuthor struct {
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// MessageEmbedField is an embedded field in a Discord message.
type MessageEmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}
