package discord

import (
	"encoding/json"
	"net/http"

	"github.com/skwair/discord/internal/endpoint"
)

// Emoji represents a Discord emoji (both standard and custom).
type Emoji struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Roles         []Role `json:"roles"`
	User          *User  `json:"user"` // The user that created this emoji.
	RequireColons bool   `json:"require_colons"`
	Managed       bool   `json:"managed"`
	Animated      bool   `json:"animated"`
}

// GetGuildEmojis returns the list of emojis for the given guild.
// Requires the MANAGE_EMOJIS permission.
func (c *Client) GetGuildEmojis(guildID string) ([]Emoji, error) {
	e := endpoint.GetGuildEmojis(guildID)
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var emojis []Emoji
	if err = json.NewDecoder(resp.Body).Decode(&emojis); err != nil {
		return nil, err
	}
	return emojis, nil
}

// GetGuildEmoji returns an emoji from a guild.
func (c *Client) GetGuildEmoji(guildID, emojiID string) (*Emoji, error) {
	e := endpoint.GetGuildEmoji(guildID, emojiID)
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var emoji Emoji
	if err = json.NewDecoder(resp.Body).Decode(&emoji); err != nil {
		return nil, err
	}
	return &emoji, nil

}

// CreateGuildEmoji creates a new emoji for the guild. image is the base64 encoded data of a
// 128*128 image. Requires the 'MANAGE_EMOJIS' permission. Fires a Guild Emojis Update
// Gateway event.
func (c *Client) CreateGuildEmoji(guildID, name, image string, roles []string) (*Emoji, error) {
	s := struct {
		Name  string   `json:"name"`
		Image string   `json:"image"`
		Roles []string `json:"roles"`
	}{
		Name:  name,
		Image: image,
		Roles: roles,
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	e := endpoint.CreateGuildEmoji(guildID)
	resp, err := c.doReq(http.MethodPost, e, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, apiError(resp)
	}

	var emoji Emoji
	if err = json.NewDecoder(resp.Body).Decode(&emoji); err != nil {
		return nil, err
	}
	return &emoji, nil
}

// ModifyGuildEmoji modifies the given emoji for the given guild. Requires the
// 'MANAGE_EMOJIS' permission. Fires a Guild Emojis Update Gateway event.
func (c *Client) ModifyGuildEmoji(guildID, emojiID, name string, roles []string) (*Emoji, error) {
	s := struct {
		Name  string   `json:"name"`
		Roles []string `json:"roles"`
	}{
		Name:  name,
		Roles: roles,
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyGuildEmoji(guildID, emojiID)
	resp, err := c.doReq(http.MethodPatch, e, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var emoji Emoji
	if err = json.NewDecoder(resp.Body).Decode(&emoji); err != nil {
		return nil, err
	}
	return &emoji, nil
}

// DeleteGuildEmoji deletes the given emoji. Requires the
// 'MANAGE_EMOJIS' permission. Fires a Guild Emojis Update Gateway event.
func (c *Client) DeleteGuildEmoji(guildID, emojiID string) error {
	e := endpoint.DeleteGuildEmoji(guildID, emojiID)
	resp, err := c.doReq(http.MethodDelete, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}
