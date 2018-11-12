package discord

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/skwair/discord/embed"
	"github.com/skwair/discord/internal/endpoint"
	"github.com/skwair/discord/webhook"
)

// Webhook is a low-effort way to post messages to channels in Discord.
// It do not require a bot user or authentication to use.
type Webhook struct {
	ID        string `json:"id,omitempty"`
	GuildID   string `json:"guild_id,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
	User      *User  `json:"user,omitempty"`
	Name      string `json:"name,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	Token     string `json:"token,omitempty"`
}

// CreateWebhook creates a new webhook. Requires the 'MANAGE_WEBHOOKS' permission.
// name must contain between 2 and 32 characters. avatar is an avatar data string,
// see https://discordapp.com/developers/docs/resources/user#avatar-data for more info.
// It can be left empty to have the default avatar.
func (c *Client) CreateWebhook(channelID, name, avatar string) (*Webhook, error) {
	s := struct {
		Name   string `json:"name,omitempty"`
		Avatar string `json:"avatar,omitempty"`
	}{
		Name:   name,
		Avatar: avatar,
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	e := endpoint.CreateWebhook(channelID)
	resp, err := c.doReq(http.MethodPost, e, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var w Webhook
	if err = json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return nil, err
	}
	return &w, nil
}

// GetChannelWebhooks returns the list of webhooks in the given channel.
// Requires the 'MANAGE_WEBHOOKS' permission.
func (c *Client) GetChannelWebhooks(channelID string) ([]Webhook, error) {
	e := endpoint.GetChannelWebhooks(channelID)
	return c.getWebhooks(e)
}

// GetGuildWebhooks returns the list of webhooks in the given guild.
// Requires the 'MANAGE_WEBHOOKS' permission.
func (c *Client) GetGuildWebhooks(guildID string) ([]Webhook, error) {
	e := endpoint.GetGuildWebhooks(guildID)
	return c.getWebhooks(e)
}

func (c *Client) getWebhooks(e *endpoint.Endpoint) ([]Webhook, error) {
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var webhooks []Webhook
	if err = json.NewDecoder(resp.Body).Decode(&webhooks); err != nil {
		return nil, err
	}
	return webhooks, nil
}

// GetWebhook returns a webhook given its ID.
func (c *Client) GetWebhook(id string) (*Webhook, error) {
	e := endpoint.GetWebhook(id)
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var w Webhook
	if err = json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return nil, err
	}
	return &w, nil
}

// GetWebhookWithToken is like GetWebhook except this call does not require
// authentication and returns no user in the webhook.
func GetWebhookWithToken(id, token string) (*Webhook, error) {
	url := fmt.Sprintf("/webhooks/%s/%s", id, token)
	resp, err := doReqNoAuth(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var w Webhook
	if err = json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return nil, err
	}
	return &w, nil
}

// ModifyWebhook modifies a webhook. Requires the 'MANAGE_WEBHOOKS' permission.
func (c *Client) ModifyWebhook(id string, s *webhook.Settings) (*Webhook, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyWebhook(id)
	resp, err := c.doReq(http.MethodPatch, e, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var w Webhook
	if err = json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return nil, err
	}
	return &w, nil
}

// ModifyWebhookWithToken is like ModifyWebhook except this call does not require
// authentication, does not allow to change the channel_id parameter in the webhook settings,
// and does not return a user in the webhook.
func ModifyWebhookWithToken(id, token string, s *webhook.Settings) (*Webhook, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("/webhooks/%s/%s", id, token)
	resp, err := doReqNoAuth(http.MethodPatch, url, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var w Webhook
	if err = json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return nil, err
	}
	return &w, nil
}

// DeleteWebhook deletes a webhook given its ID.
func (c *Client) DeleteWebhook(id string) error {
	e := endpoint.DeleteWebhook(id)
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

// DeleteWebhookWithToken is like DeleteWebhook except it does not require authentication.
func DeleteWebhookWithToken(id, token string) error {
	url := fmt.Sprintf("/webhooks/%s/%s", id, token)
	resp, err := doReqNoAuth(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// WebhookParameters are the parameters available when executing a
// webhook with ExecuteWebhook.
type WebhookParameters struct {
	Content   string        `json:"content,omitempty"`
	Username  string        `json:"username,omitempty"`
	AvatarURL string        `json:"avatar_url,omitempty"`
	TTS       bool          `json:"tts,omitempty"`
	Embeds    []embed.Embed `json:"embeds,omitempty"`
	Files     []File        `json:"-"`
}

// json implements the multipartPayload interface so WebhookParameters can be used as
// a payload with the multipartFromFiles method.
func (p *WebhookParameters) json() ([]byte, error) {
	return json.Marshal(p)
}

// ExecuteWebhook executes the webhook with the id id given its token and some
// execution parameters. wait indicates if we should wait for server confirmation
// of message send before response. If wait is set to false, the returned Message
// will be nil even if there is no error.
func ExecuteWebhook(id, token string, p *WebhookParameters, wait bool) (*Message, error) {
	if p == nil {
		return nil, errors.New("p is nil")
	}

	var (
		b   []byte
		h   http.Header
		err error
	)
	if len(p.Files) > 0 {
		b, h, err = multipartFromFiles(p, p.Files...)
		if err != nil {
			return nil, err
		}
	} else {
		b, err = json.Marshal(p)
		if err != nil {
			return nil, err
		}
		h = http.Header{}
		h.Set("Content-Type", "application/json")
	}

	q := url.Values{}
	q.Set("wait", strconv.FormatBool(wait))
	url := fmt.Sprintf("/webhooks/%s/%s?%s", id, token, q.Encode())
	resp, err := doReqNoAuthWithHeader(http.MethodPost, url, b, h)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !wait {
		if resp.StatusCode != http.StatusNoContent {
			return nil, apiError(resp)
		}
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var m Message
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}
