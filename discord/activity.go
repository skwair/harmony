package discord

import "time"

// ActivityType describes what the user is doing.
type ActivityType int

const (
	// ActivityTypePlaying will display "Playing {name}".
	ActivityTypePlaying ActivityType = 0
	// ActivityTypeStreaming will display "Streaming {name}".
	ActivityTypeStreaming ActivityType = 1
	// ActivityTypeListening will display "Listening to {name}".
	ActivityTypeListening ActivityType = 2
	// ActivityTypeCustom will display "{emoji} {name}".
	ActivityTypeCustom ActivityType = 4
	// ActivityTypeCompeting will display "Competing in {name}".
	ActivityTypeCompeting ActivityType = 5
)

// ActivityFlag describes an Activity.
type ActivityFlag int

const (
	ActivityFlagInstance    ActivityFlag = 1 << 0
	ActivityFlagJoin        ActivityFlag = 1 << 1
	ActivityFlagSpectate    ActivityFlag = 1 << 2
	ActivityFlagJoinRequest ActivityFlag = 1 << 3
	ActivityFlagSync        ActivityFlag = 1 << 4
	ActivityFlagPlay        ActivityFlag = 1 << 5
)

// Activity represents a user activity (playing a game, streaming, etc.).
// Bots are only able to send Name, Type, and optionally URL.
type Activity struct {
	Name string       `json:"name"`
	Type ActivityType `json:"type"`
	// Stream url, is validated when type is Streaming.
	URL string `json:"url"`
	// Time at which the activity was added to the user's session.
	CreatedAt time.Time `json:"created_at"`
	// Unix timestamps for start and/or end of the game.
	Timestamps ActivityTimestamp `json:"timestamps"`
	// Application id for the game.
	ApplicationID string `json:"application_id"`
	// What the player is currently doing.
	Details string `json:"details"`
	// The user's current party status.
	State string `json:"state"`
	// The emoji used for a custom status.
	Emoji ActivityEmoji `json:"emoji"`
	// Information for the current party of the player.
	Party ActivityParty `json:"party"`
	// Images for the presence and their hover texts.
	Assets ActivityAssets `json:"assets"`
	// Secrets for Rich Presence joining and spectating.
	Secrets ActivitySecrets `json:"secrets"`
	// Whether or not the activity is an instanced game session.
	Instance bool `json:"instance"`
	// Activity flags ORd together, describes what the payload includes.
	Flags ActivityFlag `json:"flags"`
}

// ActivityEmoji is the emoji set in a custom activity status.
type ActivityEmoji struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Animated bool   `json:"animated"`
}

// ActivitySecrets are secrets than can be attached to an activity.
type ActivitySecrets struct {
	Join     string `json:"join"`
	Spectate string `json:"spectate"`
	Match    string `json:"match"`
}

// ActivityTimestamp is the unix time (in milliseconds) of when the
// activity starts and ends.
type ActivityTimestamp struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

// ActivityParty contains information for the current party of the player.
type ActivityParty struct {
	ID   string `json:"id"`
	Size []int  `json:"size"` // Array of two integers (current_size, max_size).
}

// ActivityAssets contains images for the presence and their hover texts.
type ActivityAssets struct {
	LargeImage string `json:"large_image"`
	LargeText  string `json:"large_text"`
	SmallImage string `json:"small_image"`
	SmallText  string `json:"small_text"`
}
