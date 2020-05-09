package voiceutil

import "github.com/skwair/harmony/voice"

// FindUser tries to find the given user among the given voice states.
// Returns the voice channel ID the user is in if found, empty string if not.
func FindUser(states []voice.State, userID string) string {
	for _, state := range states {
		if state.UserID == userID && state.ChannelID != nil {
			return *state.ChannelID
		}
	}

	return ""
}
