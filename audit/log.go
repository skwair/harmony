package audit

// Log is the audit log of a Guild. It contains a list of entries that
// map to every admin actions performed on a Guild.
type Log struct {
	Entries  []LogEntry
	Users    []PartialUser
	Webhooks []PartialWebhook
}

// LogEntry represents a single rawEntry in the audit log.
// Entries are defined by the EntryType they describe.
type LogEntry interface {
	EntryType() EntryType
}

// BaseEntry contains the shared fields of every log entries.
type BaseEntry struct {
	ID       string // ID of the rawEntry.
	UserID   string // ID of the User that performed the action logged by the rawEntry.
	TargetID string // ID of the entity affected by this action.
	Reason   string // Reason why this action was performed.
}

// PartialUser contains a subset of the regular harmony.User type.
type PartialUser struct {
	ID            string `json:"id,omitempty"`
	Username      string `json:"username,omitempty"`
	Discriminator string `json:"discriminator,omitempty"`
	Avatar        string `json:"avatar,omitempty"`
	Bot           bool   `json:"bot,omitempty"`
}

// PartialWebhook contains a subset of the regular harmony.Webhook type.
type PartialWebhook struct {
	ID        string `json:"id,omitempty"`
	GuildID   string `json:"guild_id,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
}
