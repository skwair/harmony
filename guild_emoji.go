package harmony

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/skwair/harmony/internal/endpoint"
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

// Emojis returns the list of emojis of the guild.
// Requires the MANAGE_EMOJIS permission.
func (r *GuildResource) Emojis(ctx context.Context) ([]Emoji, error) {
	e := endpoint.ListGuildEmojis(r.guildID)
	resp, err := r.client.doReq(ctx, e, nil)
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

// Emoji returns an emoji from the guild.
func (r *GuildResource) Emoji(ctx context.Context, emojiID string) (*Emoji, error) {
	e := endpoint.GetGuildEmoji(r.guildID, emojiID)
	resp, err := r.client.doReq(ctx, e, nil)
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

// NewEmoji is like NewEmojiWithReason but with no particular reason.
func (r *GuildResource) NewEmoji(ctx context.Context, name, image string, roles []string) (*Emoji, error) {
	return r.NewEmojiWithReason(ctx, name, image, roles, "")
}

// NewEmojiWithReason creates a new emoji for the guild. image is the base64 encoded data of a
// 128*128 image. Requires the 'MANAGE_EMOJIS' permission. Fires a Guild Emojis Update
// Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *GuildResource) NewEmojiWithReason(ctx context.Context, name, image string, roles []string, reason string) (*Emoji, error) {
	st := struct {
		Name  string   `json:"name"`
		Image string   `json:"image"`
		Roles []string `json:"roles"`
	}{
		Name:  name,
		Image: image,
		Roles: roles,
	}
	b, err := json.Marshal(st)
	if err != nil {
		return nil, err
	}

	e := endpoint.CreateGuildEmoji(r.guildID)
	resp, err := r.client.doReqWithHeader(ctx, e, jsonPayload(b), reasonHeader(reason))
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

// ModifyEmoji is like ModifyEmojiWithReason but with no particular reason.
func (r *GuildResource) ModifyEmoji(ctx context.Context, emojiID, name string, roles []string) (*Emoji, error) {
	return r.ModifyEmojiWithReason(ctx, emojiID, name, roles, "")
}

// ModifyEmojiWithReason modifies the given emoji for the guild. Requires
// the 'MANAGE_EMOJIS' permission. Fires a Guild Emojis Update Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *GuildResource) ModifyEmojiWithReason(ctx context.Context, emojiID, name string, roles []string, reason string) (*Emoji, error) {
	st := struct {
		Name  string   `json:"name"`
		Roles []string `json:"roles"`
	}{
		Name:  name,
		Roles: roles,
	}
	b, err := json.Marshal(st)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyGuildEmoji(r.guildID, emojiID)
	resp, err := r.client.doReqWithHeader(ctx, e, jsonPayload(b), reasonHeader(reason))
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

// DeleteEmoji is like DeleteEmojiWithReason but with no particular reason.
func (r *GuildResource) DeleteEmoji(ctx context.Context, emojiID string) error {
	return r.DeleteEmojiWithReason(ctx, emojiID, "")
}

// DeleteEmojiWithReason deletes the given emoji. Requires the 'MANAGE_EMOJIS'
// permission. Fires a Guild Emojis Update Gateway event.
// The given reason will be set in the audit log entry for this action.
func (r *GuildResource) DeleteEmojiWithReason(ctx context.Context, emojiID, reason string) error {
	e := endpoint.DeleteGuildEmoji(r.guildID, emojiID)
	resp, err := r.client.doReqWithHeader(ctx, e, nil, reasonHeader(reason))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}
