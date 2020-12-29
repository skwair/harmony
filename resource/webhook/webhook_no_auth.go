package webhook

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/internal/endpoint"
	"github.com/skwair/harmony/internal/rest"
)

// GetWithToken returns a webhook given its ID an a token. The user field in
// the returned webhook will be nil.
func GetWithToken(ctx context.Context, id, token string) (*discord.Webhook, error) {
	e := endpoint.GetWebhookWithToken(id, token)
	resp, err := rest.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var w discord.Webhook
	if err = json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return nil, err
	}
	return &w, nil
}

// ModifyWithToken is like Modify on a Webhook resource except this call does not require
// authentication, does not allow to change the channel_id parameter in the webhook settings,
// and does not return a user in the webhook.
func ModifyWithToken(ctx context.Context, id, token string, s *discord.WebhookSettings) (*discord.Webhook, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyWebhookWithToken(id, token)
	resp, err := rest.Do(ctx, e, rest.JSONPayload(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var w discord.Webhook
	if err = json.NewDecoder(resp.Body).Decode(&w); err != nil {
		return nil, err
	}
	return &w, nil
}

// DeleteWithToken is like Delete on a webhook resource except it does not require authentication.
func DeleteWithToken(ctx context.Context, id, token string) error {
	e := endpoint.DeleteWebhookWithToken(id, token)
	resp, err := rest.Do(ctx, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return discord.NewAPIError(resp)
	}
	return nil
}

// Exec executes the webhook with the id id given its token and some
// execution parameters. wait indicates if we should wait for server confirmation
// of message send before response. If wait is set to false, the returned Message
// will be nil even if there is no error.
func Exec(ctx context.Context, id, token string, p *discord.WebhookParameters, wait bool) (*discord.Message, error) {
	if p == nil {
		return nil, errors.New("nil webhook parameters")
	}

	var payload *rest.Payload
	if len(p.Files) > 0 {
		b, contentType, err := rest.MultipartFromFiles(p, p.Files...)
		if err != nil {
			return nil, err
		}
		payload = rest.CustomPayload(b, contentType)
	} else {
		b, err := json.Marshal(p)
		if err != nil {
			return nil, err
		}
		payload = rest.JSONPayload(b)
	}

	q := url.Values{}
	q.Set("wait", strconv.FormatBool(wait))
	e := endpoint.ExecuteWebhook(id, token, q.Encode())
	resp, err := rest.Do(ctx, e, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !wait {
		if resp.StatusCode != http.StatusNoContent {
			return nil, discord.NewAPIError(resp)
		}
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var m discord.Message
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}
