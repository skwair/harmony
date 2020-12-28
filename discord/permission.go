package discord

// Set of permissions that can be assigned to Users and Roles.
const (
	PermissionNone               = 0x00000000 // Allows nothing.
	PermissionCreateInvite       = 0x00000001 // Allows creation of instant invites.
	PermissionKickMembers        = 0x00000002 // Allows kicking members.
	PermissionBanMembers         = 0x00000004 // Allows banning members.
	PermissionAdministrator      = 0x00000008 // Allows *all permissions* and bypasses channel permission overwrites.
	PermissionManageChannels     = 0x00000010 // Allows management and editing of channels.
	PermissionManageGuild        = 0x00000020 // Allows management and editing of the guild.
	PermissionAddReactions       = 0x00000040 // Allows for the addition of reactions to messages.
	PermissionViewAuditLog       = 0x00000080 // Allows for viewing of audit logs.
	PermissionViewChannel        = 0x00000400 // Allows guild members to view a channel, which includes reading messages in text channels.
	PermissionSendMessages       = 0x00000800 // Allows for sending messages in a channel.
	PermissionSendTTSMessages    = 0x00001000 // Allows for sending of /tts messages.
	PermissionManageMessages     = 0x00002000 // Allows for deletion of other users messages.
	PermissionEmbedLinks         = 0x00004000 // Links sent by users with this permission will be auto-embedded.
	PermissionAttachFiles        = 0x00008000 // Allows for uploading images and files.
	PermissionReadMessageHistory = 0x00010000 // Allows for reading of message history.
	PermissionMentionEveryone    = 0x00020000 // Allows for using the @everyone tag to notify all users in a channel, and the @here tag to notify all online users in a channel.
	PermissionUseExternalEmojis  = 0x00040000 // Allows the usage of custom emojis from other servers.
	PermissionConnect            = 0x00100000 // Allows for joining of a voice channel.
	PermissionSpeak              = 0x00200000 // Allows for speaking in a voice channel.
	PermissionMuteMembers        = 0x00400000 // Allows for muting members in a voice channel.
	PermissionDeafenMembers      = 0x00800000 // Allows for deafening of members in a voice channel.
	PermissionMoveMembers        = 0x01000000 // Allows for moving of members between voice channels.
	PermissionUseVAD             = 0x02000000 // Allows for using voice-activity-detection in a voice channel.
	PermissionPrioritySpeaker    = 0x00000100 // Allows for using priority speaker in a voice channel.
	PermissionChangeNickname     = 0x04000000 // Allows for modification of own nickname.
	PermissionManageNicknames    = 0x08000000 // Allows for modification of other users nicknames.
	PermissionManageRoles        = 0x10000000 // Allows management and editing of roles.
	PermissionManageWebhooks     = 0x20000000 // Allows management and editing of webhooks.
	PermissionManageEmojis       = 0x40000000 // Allows management and editing of emojis.
)

// PermissionOverwrite describes a specific permission that overwrites
// server-wide permissions.
type PermissionOverwrite struct {
	Type  int `json:"type"` // Either 0 for "role" or 1 for "member".
	ID    string `json:"id"`   // ID of the role or member, depending on Type.
	Allow int    `json:"allow,string"`
	Deny  int    `json:"deny,string"`
}

// PermissionsContains returns whether the given permission is set in permissions.
func PermissionsContains(permissions, permission int) bool {
	return permissions&permission == permission
}

// computeBasePermissions returns the base permissions a member has in a given guild.
func computeBasePermissions(g *Guild, m *GuildMember) (permissions int) {
	if g.OwnerID == m.User.ID {
		return PermissionAdministrator
	}

	// Role '@everyone' has the same ID as the guild ID
	// and is always present, so no need to check for nil.
	roleEveryone := roleByID(g.Roles, g.ID)
	permissions = roleEveryone.Permissions

	for _, id := range m.Roles {
		role := roleByID(g.Roles, id)
		if role != nil {
			permissions |= role.Permissions
		}
	}

	if PermissionsContains(permissions, PermissionAdministrator) {
		return PermissionAdministrator
	}
	return permissions
}

func computeOverwrites(ch *Channel, m *GuildMember, basePermissions int) (permissions int) {
	// Administrator can not be overridden.
	if PermissionsContains(basePermissions, PermissionAdministrator) {
		return PermissionAdministrator
	}

	permissions = basePermissions

	po := overwriteByID(ch.PermissionOverwrites, ch.GuildID)
	if po != nil {
		permissions &= ^po.Deny
		permissions |= po.Allow
	}

	pos := ch.PermissionOverwrites
	allow := PermissionNone
	deny := PermissionNone
	for _, id := range m.Roles {
		por := overwriteByID(pos, id)
		if por != nil {
			allow |= por.Allow
			deny |= por.Deny
		}
	}
	permissions &= ^deny
	permissions |= allow

	pom := overwriteByID(ch.PermissionOverwrites, m.User.ID)
	if pom != nil {
		permissions &= ^pom.Deny
		permissions |= pom.Allow
	}

	return permissions
}

func roleByID(roles []Role, id string) *Role {
	for i := 0; i < len(roles); i++ {
		if roles[i].ID == id {
			return &roles[i]
		}
	}
	return nil
}

func overwriteByID(po []PermissionOverwrite, id string) *PermissionOverwrite {
	for i := 0; i < len(po); i++ {
		if po[i].ID == id {
			return &po[i]
		}
	}
	return nil
}
