package permission

// Set of permissions that can be assigned to Users and Roles.
const (
	None               = 0x00000000 // Allows nothing.
	CreateInvite       = 0x00000001 // Allows creation of instant invites.
	KickMembers        = 0x00000002 // Allows kicking members.
	BanMembers         = 0x00000004 // Allows banning members.
	Administrator      = 0x00000008 // Allows *all permissions* and bypasses channel permission overwrites.
	ManageChannels     = 0x00000010 // Allows management and editing of channels.
	ManageGuild        = 0x00000020 // Allows management and editing of the guild.
	AddReactions       = 0x00000040 // Allows for the addition of reactions to messages.
	ViewAuditLog       = 0x00000080 // Allows for viewing of audit logs.
	ViewChannel        = 0x00000400 // Allows guild members to view a channel, which includes reading messages in text channels.
	SendMessages       = 0x00000800 // Allows for sending messages in a channel.
	SendTTSMessages    = 0x00001000 // Allows for sending of /tts messages.
	ManageMessages     = 0x00002000 // Allows for deletion of other users messages.
	EmbedLinks         = 0x00004000 // Links sent by users with this permission will be auto-embedded.
	AttachFiles        = 0x00008000 // Allows for uploading images and files.
	ReadMessageHistory = 0x00010000 // Allows for reading of message history.
	MentionEveryone    = 0x00020000 // Allows for using the @everyone tag to notify all users in a channel, and the @here tag to notify all online users in a channel.
	UseExternalEmojis  = 0x00040000 // Allows the usage of custom emojis from other servers.
	Connect            = 0x00100000 // Allows for joining of a voice channel.
	Speak              = 0x00200000 // Allows for speaking in a voice channel.
	MuteMembers        = 0x00400000 // Allows for muting members in a voice channel.
	DeafenMembers      = 0x00800000 // Allows for deafening of members in a voice channel.
	MoveMembers        = 0x01000000 // Allows for moving of members between voice channels.
	UseVAD             = 0x02000000 // Allows for using voice-activity-detection in a voice channel.
	ChangeNickname     = 0x04000000 // Allows for modification of own nickname.
	ManageNicknames    = 0x08000000 // Allows for modification of other users nicknames.
	ManageRoles        = 0x10000000 // Allows management and editing of roles.
	ManageWebhooks     = 0x20000000 // Allows management and editing of webhooks.
	ManageEmojis       = 0x40000000 // Allows management and editing of emojis.
	All                = 0x7FF7FCFF // Equivalent to all permissions, OR'd.
)

// Overwrite describes a specific permission that overwrites
// server-wide permissions.
type Overwrite struct {
	ID    string
	Type  string // Either "role" or "member".
	Allow int
	Deny  int
}

// Contains returns whether the given permission is set in permissions.
func Contains(permissions, permission int) bool {
	return permissions&permission == permission
}
