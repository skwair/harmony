package discord

// Status is the current presence status of a user.
type Status string

// Valid presence status:
const (
	StatusOnline    Status = "online"
	StatusIdle      Status = "idle"
	StatusDND       Status = "dnd"
	StatusInvisible Status = "invisible"
	StatusOffline   Status = "offline"
)

// Presence is a user's current state on a guild.
// This event is sent when a user's presence is updated for a guild.
type Presence struct {
	User         *User        `json:"user"`
	GuildID      string       `json:"guild_id"`
	Status       Status       `json:"status"`
	Activities   []Activity   `json:"activities"`
	ClientStatus ClientStatus `json:"client_status"`
}

// ClientStatus is a platform-specific client status.
type ClientStatus struct {
	Desktop string `json:"desktop"`
	Mobile  string `json:"mobile"`
	Web     string `json:"web"`
}

// BotStatus is sent by the client to indicate a presence or status update.
type BotStatus struct {
	Since      int                 `json:"since"`
	Activities []BotStatusActivity `json:"activities,omitempty"`
	Status     Status              `json:"status,omitempty"`
	AFK        bool                `json:"afk"`
}

// BotStatusActivity is the subset of allowed values for bots from regular Activities.
type BotStatusActivity struct {
	Type ActivityType `json:"type"`
	Name string       `json:"name"`
	URL  string       `json:"url,omitempty"` // Used when Type is ActivityTypeStreaming.
}
