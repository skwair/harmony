package audit

import "github.com/skwair/harmony/permission"

// EntryType defines the type of event an entry describes.
type EntryType int

const (
	EntryTypeGuildUpdate            EntryType = 1
	EntryTypeChannelCreate          EntryType = 10
	EntryTypeChannelUpdate          EntryType = 11
	EntryTypeChannelDelete          EntryType = 12
	EntryTypeChannelOverwriteCreate EntryType = 13
	EntryTypeChannelOverwriteUpdate EntryType = 14
	EntryTypeChannelOverwriteDelete EntryType = 15
	EntryTypeMemberKick             EntryType = 20
	EntryTypeMemberPrune            EntryType = 21
	EntryTypeMemberBanAdd           EntryType = 22
	EntryTypeMemberBanRemove        EntryType = 23
	EntryTypeMemberUpdate           EntryType = 24
	EntryTypeMemberRoleUpdate       EntryType = 25
	EntryTypeRoleCreate             EntryType = 30
	EntryTypeRoleUpdate             EntryType = 31
	EntryTypeRoleDelete             EntryType = 32
	EntryTypeInviteCreate           EntryType = 40
	EntryTypeInviteUpdate           EntryType = 41
	EntryTypeInviteDelete           EntryType = 42
	EntryTypeWebhookCreate          EntryType = 50
	EntryTypeWebhookUpdate          EntryType = 51
	EntryTypeWebhookDelete          EntryType = 52
	EntryTypeEmojiCreate            EntryType = 60
	EntryTypeEmojiUpdate            EntryType = 61
	EntryTypeEmojiDelete            EntryType = 62
	EntryTypeMessageDelete          EntryType = 72
)

// ChannelUpdate is the log entry that describes a Channel creation.
// It contains settings this Channel was created with.
type ChannelCreate struct {
	BaseEntry

	Name                 string
	Type                 int
	RateLimitPerUser     int
	NSFW                 bool
	PermissionOverwrites []permission.Overwrite
}

// EntryType implements the LogEntry interface.
func (ChannelCreate) EntryType() EntryType { return EntryTypeChannelCreate }

// ChannelUpdate is the log entry that describes how a Channel was updated.
// It contains a list of settings that can be updated on a Channel.
// Settings that are not nil are those which were modified. They contain both
// their old value as well as the new one.
type ChannelUpdate struct {
	BaseEntry

	Name             *StringValues
	Topic            *StringValues
	Bitrate          *IntValues
	RateLimitPerUser *IntValues
	NSFW             *BoolValues
	ApplicationID    *StringValues // Application ID of the added or removed webhook or bot.
}

// EntryType implements the LogEntry interface.
func (ChannelUpdate) EntryType() EntryType { return EntryTypeChannelUpdate }

// ChannelUpdate is the log entry that describes a Channel deletion.
// It contains settings this Channel had before being deleted.
type ChannelDelete struct {
	BaseEntry

	Name                 string
	Type                 int
	RateLimitPerUser     int
	NSFW                 bool
	PermissionOverwrites []permission.Overwrite
}

// EntryType implements the LogEntry interface.
func (ChannelDelete) EntryType() EntryType { return EntryTypeChannelDelete }

// ChannelOverwriteCreate is the log entry that describes a Channel PermissionOverwrite creation.
// It contains settings this overwrite was created with.
type ChannelOverwriteCreate struct {
	BaseEntry

	Type  string
	ID    string
	Allow int
	Deny  int

	RoleName string // Name of the role if Type is "role".
}

// EntryType implements the LogEntry interface.
func (ChannelOverwriteCreate) EntryType() EntryType { return EntryTypeChannelOverwriteCreate }

// ChannelOverwriteCreate is the log entry that describes how a Channel PermissionOverwrite was updated.
// It contains a list of settings that can be updated on a Channel PermissionOverwrite.
// Settings that are not nil are those which were modified. They contain both
// their old value as well as the new one.
type ChannelOverwriteUpdate struct {
	BaseEntry

	Allow *IntValues
	Deny  *IntValues

	Type     string
	ID       string
	RoleName string // Name of the role if Type is "role".
}

// EntryType implements the LogEntry interface.
func (ChannelOverwriteUpdate) EntryType() EntryType { return EntryTypeChannelOverwriteUpdate }

// ChannelOverwriteDelete is the log entry that describes a Channel PermissionOverwrite deletion.
// It contains settings this overwrite had before being deleted.
type ChannelOverwriteDelete struct {
	BaseEntry

	Type  string
	ID    string
	Allow int
	Deny  int

	RoleName string // Name of the role if Type is "role".
}

// EntryType implements the LogEntry interface.
func (ChannelOverwriteDelete) EntryType() EntryType { return EntryTypeChannelOverwriteDelete }

// MemberKick is the log entry that describes a member kick.
type MemberKick struct {
	BaseEntry
}

// EntryType implements the LogEntry interface.
func (MemberKick) EntryType() EntryType { return EntryTypeMemberKick }

// MemberKick is the log entry that describes a member prune.
type MemberPrune struct {
	BaseEntry

	DeleteMemberDays int
	MembersRemoved   int
}

// EntryType implements the LogEntry interface.
func (MemberPrune) EntryType() EntryType { return EntryTypeMemberPrune }

// MemberBanAdd is the log entry that describes a member ban creation.
type MemberBanAdd struct {
	BaseEntry
}

// EntryType implements the LogEntry interface.
func (MemberBanAdd) EntryType() EntryType { return EntryTypeMemberBanAdd }

// MemberBanRemove is the log entry that describes a member ban deletion.
type MemberBanRemove struct {
	BaseEntry
}

// EntryType implements the LogEntry interface.
func (MemberBanRemove) EntryType() EntryType { return EntryTypeMemberBanRemove }

// MemberUpdate is the log entry that describes a member update.
type MemberUpdate struct {
	BaseEntry

	Nick *StringValues
	Deaf *BoolValues
	Mute *BoolValues
}

// EntryType implements the LogEntry interface.
func (MemberUpdate) EntryType() EntryType { return EntryTypeMemberUpdate }

// MemberRoleUpdate is the log entry that describes a member's roles update.
// It contains roles that were added as well as roles that were removed.
type MemberRoleUpdate struct {
	BaseEntry

	Added   []permission.Overwrite
	Removed []permission.Overwrite
}

// EntryType implements the LogEntry interface.
func (MemberRoleUpdate) EntryType() EntryType { return EntryTypeMemberRoleUpdate }

// RoleCreate is the log entry that describes a role creation.
// It contains the settings the role was created with.
type RoleCreate struct {
	BaseEntry

	Name        string
	Permissions int
	Color       int
	Mentionable bool
	Hoist       bool
}

// EntryType implements the LogEntry interface.
func (RoleCreate) EntryType() EntryType { return EntryTypeRoleCreate }

// RoleUpdate is the log entry that describes how a Role was updated.
// It contains a list of settings that can be updated on a Role.
// Settings that are not nil are those which were modified. They contain both
// their old value as well as the new one.
type RoleUpdate struct {
	BaseEntry

	Name        *StringValues
	Permissions *IntValues
	Color       *IntValues
	Mentionable *BoolValues
	Hoist       *BoolValues
}

// EntryType implements the LogEntry interface.
func (RoleUpdate) EntryType() EntryType { return EntryTypeRoleUpdate }

// RoleDelete is the log entry that describes a role deletion.
// It contains settings this role had before being deleted.
type RoleDelete struct {
	BaseEntry

	Name        string
	Permissions int
	Color       int
	Mentionable bool
	Hoist       bool
}

// EntryType implements the LogEntry interface.
func (RoleDelete) EntryType() EntryType { return EntryTypeRoleDelete }

// InviteCreate is the log entry that describes a channel invite creation.
// It contains the settings the invite was created with.
type InviteCreate struct {
	BaseEntry

	Code      string
	ChannelID string
	InviterID string
	MaxUses   int
	Uses      int
	MaxAge    int
	Temporary bool
}

// EntryType implements the LogEntry interface.
func (InviteCreate) EntryType() EntryType { return EntryTypeInviteCreate }

// InviteUpdate is the log entry that describes how a channel invite was updated.
// It contains a list of settings that can be updated on an invite.
// Settings that are not nil are those which were modified. They contain both
// their old value as well as the new one.
type InviteUpdate struct {
	BaseEntry

	Code      *StringValues
	ChannelID *StringValues
	InviterID *StringValues
	MaxUses   *IntValues
	Uses      *IntValues
	MaxAge    *IntValues
	Temporary *BoolValues
}

// EntryType implements the LogEntry interface.
func (InviteUpdate) EntryType() EntryType { return EntryTypeInviteUpdate }

// InviteDelete is the log entry that describes a channel invite deletion.
// It contains settings this invite had before being deleted.
type InviteDelete struct {
	BaseEntry

	Code      string
	ChannelID string
	InviterID string
	MaxUses   int
	Uses      int
	MaxAge    int
	Temporary bool
}

// EntryType implements the LogEntry interface.
func (InviteDelete) EntryType() EntryType { return EntryTypeInviteDelete }