package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/internal/endpoint"
	"github.com/skwair/harmony/internal/rest"
)

// Get returns the current user.
func (r *Resource) Get(ctx context.Context) (*discord.User, error) {
	e := endpoint.GetUser(r.userID)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var u discord.User
	if err = json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

// Modify modifies the current user account settings.
// Avatar is a Data URI scheme that supports JPG, GIF, and PNG formats.
// An example Data URI format is:
//
//     data:image/jpeg;base64,BASE64_ENCODED_JPEG_IMAGE_DATA
//
// Ensure you use the proper header type (image/jpeg, image/png, image/gif)
// that matches the image data being provided.
func (r *Resource) Modify(ctx context.Context, username, avatar string) (*discord.User, error) {
	if r.userID != "@me" {
		return nil, discord.ErrNotCurrentUser
	}

	st := struct {
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
	}{
		Username: username,
		Avatar:   avatar,
	}
	b, err := json.Marshal(st)
	if err != nil {
		return nil, err
	}

	e := endpoint.ModifyCurrentUser()
	resp, err := r.client.Do(ctx, e, rest.JSONPayload(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var u discord.User
	if err = json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

// Guilds returns a list of partial guilds the current
// user is a member of. This endpoint returns at most 100 guilds by
// default, which is the maximum number of guilds a non-bot user can
// join. Therefore, pagination is not needed for integrations that need
// to get a list of users' guilds.
func (r *Resource) Guilds(ctx context.Context) ([]discord.PartialGuild, error) {
	if r.userID != "@me" {
		return nil, discord.ErrNotCurrentUser
	}

	e := endpoint.GetCurrentUserGuilds()
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var guilds []discord.PartialGuild
	if err = json.NewDecoder(resp.Body).Decode(&guilds); err != nil {
		return nil, err
	}
	return guilds, nil
}

// LeaveGuild make the current user leave a guild given its ID.
func (r *Resource) LeaveGuild(ctx context.Context, id string) error {
	e := endpoint.LeaveGuild(id)
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return discord.NewAPIError(resp)
	}
	return nil
}

// DMs returns the list of direct message channels the current user is in.
// This endpoint does not seem to be available for Bot users, always returning
// an empty list of channels.
func (r *Resource) DMs(ctx context.Context) ([]discord.Channel, error) {
	if r.userID != "@me" {
		return nil, discord.ErrNotCurrentUser
	}

	e := endpoint.GetUserDMs()
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var channels []discord.Channel
	if err = json.NewDecoder(resp.Body).Decode(&channels); err != nil {
		return nil, err
	}
	return channels, nil
}

// NewDM creates a new DM channel with a user. Returns the created channel.
// If a DM channel already exist with this recipient, it does not create a new
// one and returns the existing one instead.
func (r *Resource) NewDM(ctx context.Context, recipientID string) (*discord.Channel, error) {
	if r.userID != "@me" {
		return nil, discord.ErrNotCurrentUser
	}

	st := struct {
		RecipientID string `json:"recipient_id"`
	}{
		RecipientID: recipientID,
	}
	b, err := json.Marshal(st)
	if err != nil {
		return nil, err
	}

	e := endpoint.CreateDM()
	resp, err := r.client.Do(ctx, e, rest.JSONPayload(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var ch discord.Channel
	if err = json.NewDecoder(resp.Body).Decode(&ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

// Connections returns a list of connections for the connected user.
func (r *Resource) Connections(ctx context.Context) ([]discord.UserConnection, error) {
	if r.userID != "@me" {
		return nil, discord.ErrNotCurrentUser
	}

	e := endpoint.GetUserConnections()
	resp, err := r.client.Do(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, discord.NewAPIError(resp)
	}

	var conns []discord.UserConnection
	if err = json.NewDecoder(resp.Body).Decode(&conns); err != nil {
		return nil, err
	}
	return conns, nil
}
