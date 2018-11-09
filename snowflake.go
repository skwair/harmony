package discord

import (
	"strconv"
	"time"
)

// CreationTimeOf returns the creation time of the given Discord ID (userID, guildID, channelID).
// For more information, see : https://discordapp.com/developers/docs/reference#snowflakes.
func CreationTimeOf(id string) (time.Time, error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	ts := (i >> 22) + 1420070400000

	return time.Unix(ts/1000, 0), nil
}
