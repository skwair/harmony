package harmony

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/skwair/harmony/internal/endpoint"
)

// User in Discord is generally considered the base entity.
// Users can spawn across the entire platform, be members of guilds,
// participate in text and voice chat, and much more. Users are separated
// by a distinction of "bot" vs "normal." Although they are similar,
// bot users are automated users that are "owned" by another user.
// Unlike normal users, bot users do not have a limitation on the number
// of Guilds they can be a part of.
type User struct {
	ID            string `json:"id,omitempty"`
	Username      string `json:"username,omitempty"`
	Discriminator string `json:"discriminator,omitempty"`
	Avatar        string `json:"avatar,omitempty"`
	Bot           bool   `json:"bot,omitempty"`
	MFAEnabled    bool   `json:"mfa_enabled,omitempty"`
	Verified      bool   `json:"verified,omitempty"`
	Email         string `json:"email,omitempty"`
}

// AvatarURL returns the user's avatar URL.
func (u *User) AvatarURL() string {
	if u.Avatar == "" {
		d, _ := strconv.ParseInt(u.Discriminator, 10, 64)
		return fmt.Sprintf("https://cdn.discordapp.com/embed/avatars/%d.png", d%5)
	}
	return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", u.ID, u.Avatar)
}

// GetUser returns a user  given its ID. Use "@me" as the ID to fetch information
// about the connected user. For every other IDs, this endpoint can only be used by bots.
func (c *Client) GetUser(id string) (*User, error) {
	e := endpoint.GetUser(id)
	resp, err := c.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var u User
	if err = json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

// CurrentUserResource is a resource that allows to perform various actions on the
// current user. Create one with Client.Channel.
type CurrentUserResource struct {
	client *Client
}

// CurrentUser returns a new resource to manage the current user.
func (c *Client) CurrentUser() *CurrentUserResource {
	return &CurrentUserResource{
		client: c,
	}
}

// Get returns the current user.
func (r *CurrentUserResource) Get() (*User, error) {
	e := endpoint.GetUser("@me")
	resp, err := r.client.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var u User
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
func (r *CurrentUserResource) Modify(username, avatar string) (*User, error) {
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
	resp, err := r.client.doReq(http.MethodPatch, e, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var u User
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
func (r *CurrentUserResource) Guilds() ([]PartialGuild, error) {
	e := endpoint.GetCurrentUserGuilds()
	resp, err := r.client.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var guilds []PartialGuild
	if err = json.NewDecoder(resp.Body).Decode(&guilds); err != nil {
		return nil, err
	}
	return guilds, nil
}

// LeaveGuild make the current user leave a guild given its ID.
func (r *CurrentUserResource) LeaveGuild(id string) error {
	e := endpoint.LeaveGuild(id)
	resp, err := r.client.doReq(http.MethodDelete, e, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return apiError(resp)
	}
	return nil
}

// DMs returns the list of direct message channels the current user is in.
// This endpoint does not seem to be available for Bot users, always returning
// an empty list of channels.
func (r *CurrentUserResource) DMs(id string) ([]Channel, error) {
	e := endpoint.GetUserDMs()
	resp, err := r.client.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var channels []Channel
	if err = json.NewDecoder(resp.Body).Decode(&channels); err != nil {
		return nil, err
	}
	return channels, nil
}

// NewDM creates a new DM channel with a user. Returns the created channel.
// If a DM channel already exist with this recipient, it does not create a new
// one and returns the existing one instead.
func (r *CurrentUserResource) NewDM(recipientID string) (*Channel, error) {
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
	resp, err := r.client.doReq(http.MethodPost, e, b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var ch Channel
	if err = json.NewDecoder(resp.Body).Decode(&ch); err != nil {
		return nil, err
	}
	return &ch, nil
}

// Connections returns a list of connections for the connected user.
func (r *CurrentUserResource) Connections() ([]Connection, error) {
	e := endpoint.GetUserConnections()
	resp, err := r.client.doReq(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var conns []Connection
	if err = json.NewDecoder(resp.Body).Decode(&conns); err != nil {
		return nil, err
	}
	return conns, nil
}

// SetStatus sets the current user's status. You need to be connected to the
// Gateway to call this method, else it will return ErrGatewayNotConnected.
func (r *CurrentUserResource) SetStatus(status *Status) error {
	if !r.client.isConnected() {
		return ErrGatewayNotConnected
	}

	return r.client.sendPayload(gatewayOpcodeStatusUpdate, status)
}
