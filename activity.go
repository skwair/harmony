package discord

// ActivityType describes what the user is doing.
type ActivityType int

const (
	// ActivityPlaying will display "Playing {name}".
	ActivityPlaying ActivityType = iota
	// ActivityStreaming will display "Streaming {name}".
	ActivityStreaming
	// ActivityListening will display "Listening to {name}".
	ActivityListening
)

// Activity represents a user activity (playing a game, streaming, etc.).
type Activity struct {
	Name string       `json:"name"`
	Type ActivityType `json:"type"`
	// Stream url, is validated when type is Streaming.
	URL string `json:"url,omitempty"`
	// Unix timestamps for start and/or end of the game.
	Timestamps *ActivityTimestamp `json:"timestamps,omitempty"`
	// Application id for the game.
	ApplicationID string `json:"application_id,omitempty"`
	// What the player is currently doing.
	Details string `json:"details,omitempty"`
	// The user's current party status.
	State string `json:"state,omitempty"`
	// Information for the current party of the player.
	Party *ActivityParty `json:"party,omitempty"`
	// Images for the presence and their hover texts.
	Assets *ActivityAssets `json:"assets,omitempty"`
}

// ActivityTimestamp is the unix time (in milliseconds) of when the
// activity starts and ends.
type ActivityTimestamp struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

// ActivityParty contains information for the current party of the player.
type ActivityParty struct {
	ID   string `json:"id,omitempty"`
	Size []int  `json:"size,omitempty"` // Array of two integers (current_size, max_size).
}

// ActivityAssets contains images for the presence and their hover texts.
type ActivityAssets struct {
	LargeImage string `json:"large_image,omitempty"`
	LargeText  string `json:"large_text,omitempty"`
	SmallImage string `json:"small_image,omitempty"`
	SmallText  string `json:"small_text,omitempty"`
}
