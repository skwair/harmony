package discord

import (
	"fmt"
	"strconv"
)

// UserFlag are specific attributes a User can have.
type UserFlag int

const (
	UserFlagNone                      UserFlag = 0
	UserFlagDiscordEmployee           UserFlag = 1 << 0
	UserFlagPartneredServerOwner      UserFlag = 1 << 1
	UserFlagHypeSquadEvents           UserFlag = 1 << 2
	UserFlagBugHunterLevel1           UserFlag = 1 << 3
	UserFlagHouseBravery              UserFlag = 1 << 6
	UserFlagHouseBrilliance           UserFlag = 1 << 7
	UserFlagHouseBalance              UserFlag = 1 << 8
	UserFlagEarlySupporter            UserFlag = 1 << 9
	UserFlagTeamUser                  UserFlag = 1 << 10
	UserFlagSystem                    UserFlag = 1 << 12
	UserFlagBugHunterLevel2           UserFlag = 1 << 14
	UserFlagVerifiedBot               UserFlag = 1 << 16
	UserFlagEarlyVerifiedBotDeveloper UserFlag = 1 << 17
)

// PremiumType denotes the level of premium a user has.
type PremiumType int

const (
	PremiumTypeNone         PremiumType = 0
	PremiumTypeNitroClassic PremiumType = 1
	PremiumTypeNitro        PremiumType = 2
)

type Visibility int

const (
	VisibilityNone     Visibility = 0
	VisibilityEveryone Visibility = 1
)

// User in Discord is generally considered the base entity.
// Users can spawn across the entire platform, be members of guilds,
// participate in text and voice chat, and much more. Users are separated
// by a distinction of "bot" vs "normal." Although they are similar,
// bot users are automated users that are "owned" by another user.
// Unlike normal users, bot users do not have a limitation on the number
// of Guilds they can be a part of.
type User struct {
	ID            string      `json:"id"`
	Username      string      `json:"username"`
	Discriminator string      `json:"discriminator"`
	Avatar        string      `json:"avatar"`
	Locale        string      `json:"locale"`
	Email         string      `json:"email"`
	Verified      bool        `json:"verified"`
	MFAEnabled    bool        `json:"mfa_enabled"`
	Bot           bool        `json:"bot"`
	System        bool        `json:"system"`
	PremiumType   PremiumType `json:"premium_type"`
	Flags         UserFlag    `json:"flags"`
	PublicFlags   UserFlag    `json:"public_flags"`
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
	Verified     bool               `json:"verified"`
	FriendSync   bool               `json:"friend_sync"`
	ShowActivity bool               `json:"show_activity"`
	Visibility   Visibility         `json:"visibility"`
}
