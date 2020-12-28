package discord

import (
	"fmt"
	"strconv"
)

// User in Discord is generally considered the base entity.
// Users can spawn across the entire platform, be members of guilds,
// participate in text and voice chat, and much more. Users are separated
// by a distinction of "bot" vs "normal." Although they are similar,
// bot users are automated users that are "owned" by another user.
// Unlike normal users, bot users do not have a limitation on the number
// of Guilds they can be a part of.
type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Bot           bool   `json:"bot"`
	MFAEnabled    bool   `json:"mfa_enabled"`
	Verified      bool   `json:"verified"`
	Email         string `json:"email"`
}

// AvatarURL returns the user's avatar URL.
func (u *User) AvatarURL() string {
	if u.Avatar == "" {
		d, _ := strconv.ParseInt(u.Discriminator, 10, 64)
		return fmt.Sprintf("https://cdn.discordapp.com/embed/avatars/%d.png", d%5)
	}
	return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", u.ID, u.Avatar)
}

// UserConnection that the user has attached.
type UserConnection struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Type         string             `json:"type"`
	Revoked      bool               `json:"revoked"`
	Integrations []GuildIntegration `json:"integrations"` // Partial server integrations.
}
