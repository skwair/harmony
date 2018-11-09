package discord

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/skwair/discord/internal/endpoint"
)

// Reaction is a reaction on a Discord message.
type Reaction struct {
	Count int    `json:"count"`
	Me    bool   `json:"me"`
	Emoji *Emoji `json:"emoji"`
}

// CreateReaction adds a reaction to a message in a channel. This endpoint requires
// the 'READ_MESSAGE_HISTORY' permission to be present on the current user. Additionally,
// if nobody else has reacted to the message using this emoji, this endpoint requires
// the'ADD_REACTIONS' permission to be present on the current user.
func (c *Client) CreateReaction(channelID, messageID, emoji string) error {
	e := endpoint.CreateReaction(channelID, messageID, url.PathEscape(emoji))
	resp, err := c.doReq(http.MethodPut, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// DeleteReaction deletes a reaction the current user has made for the message.
func (c *Client) DeleteReaction(channelID, messageID, emoji string) error {
	e := endpoint.DeleteReaction(channelID, messageID, url.PathEscape(emoji))
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

// DeleteUserReaction deletes another user's reaction. This endpoint requires the
// 'MANAGE_MESSAGES' permission to be present on the current user.
func (c *Client) DeleteUserReaction(channelID, messageID, userID, emoji string) error {
	e := endpoint.DeleteUserReaction(channelID, messageID, userID, url.PathEscape(emoji))
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

// DeleteAllReactions deletes all reactions on a message. This endpoint requires the
// 'MANAGE_MESSAGES' permission to be present on the current user.
func (c *Client) DeleteAllReactions(channelID, messageID string) error {
	e := endpoint.DeleteAllReactions(channelID, messageID)
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

// GetReactions returns a list of users that reacted to a message with the given emoji.
// limit is the number of users to return and can be set to any value ranging from 1 to 100.
// If set to 0, it defaults to 25. If more than 100 users reacted with the given emoji,
// the before and after parameters can be used to fetch more users.
func (c *Client) GetReactions(channelID, messageID, emoji string, limit int, before, after string) ([]User, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if before != "" {
		q.Set("before", before)
	}
	if after != "" {
		q.Set("after", after)
	}

	e := endpoint.GetReactions(channelID, messageID, url.PathEscape(emoji), q.Encode())
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var users []User
	if err = json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}
	return users, nil
}
