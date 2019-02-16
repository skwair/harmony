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