package audit

import (
	"fmt"
	"strconv"
)

func channelCreateFromEntry(e *rawEntry) (*ChannelCreate, error) {
	chCreate := &ChannelCreate{
		BaseEntry: baseEntryFromRaw(e),
	}

	var err error
	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyName:
			chCreate.Name, err = stringValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyName, err)
			}

		case changeKeyType:
			chCreate.Type, err = intValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyType, err)
			}

		case changeKeyRateLimitPerUser:
			chCreate.RateLimitPerUser, err = intValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyRateLimitPerUser, err)
			}

		case changeKeyNFSW:
			chCreate.NSFW, err = boolValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyNFSW, err)
			}

		case changeKeyPermissionOverwrites:
			chCreate.PermissionOverwrites, err = permissionOverwritesValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyPermissionOverwrites, err)
			}
		}
	}

	return chCreate, nil
}

func channelUpdateFromEntry(e *rawEntry) (*ChannelUpdate, error) {
	chUpdate := &ChannelUpdate{
		BaseEntry: baseEntryFromRaw(e),
	}

	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyName:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyName, err)
			}
			chUpdate.Name = &StringValues{Old: oldValue, New: newValue}

		case changeKeyTopic:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyTopic, err)
			}
			chUpdate.Topic = &StringValues{Old: oldValue, New: newValue}

		case changeKeyBitrate:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyBitrate, err)
			}
			chUpdate.Bitrate = &IntValues{Old: oldValue, New: newValue}

		case changeKeyRateLimitPerUser:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyRateLimitPerUser, err)
			}
			chUpdate.RateLimitPerUser = &IntValues{Old: oldValue, New: newValue}

		case changeKeyNFSW:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyNFSW, err)
			}
			chUpdate.NSFW = &BoolValues{Old: oldValue, New: newValue}

		case changeKeyApplicationID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyApplicationID, err)
			}
			chUpdate.ApplicationID = &StringValues{Old: oldValue, New: newValue}

		case changeKeyPosition:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyPosition, err)
			}
			chUpdate.Position = &IntValues{Old: oldValue, New: newValue}
		}
	}

	return chUpdate, nil
}

func channelDeleteFromEntry(e *rawEntry) (*ChannelDelete, error) {
	chDelete := &ChannelDelete{
		BaseEntry: baseEntryFromRaw(e),
	}

	var err error
	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyName:
			chDelete.Name, err = stringValue(ch.Old)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyName, err)
			}

		case changeKeyType:
			chDelete.Type, err = intValue(ch.Old)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyType, err)
			}

		case changeKeyRateLimitPerUser:
			chDelete.RateLimitPerUser, err = intValue(ch.Old)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyRateLimitPerUser, err)
			}

		case changeKeyNFSW:
			chDelete.NSFW, err = boolValue(ch.Old)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyNFSW, err)
			}

		case changeKeyPermissionOverwrites:
			chDelete.PermissionOverwrites, err = permissionOverwritesValue(ch.Old)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyPermissionOverwrites, err)
			}
		}
	}

	return chDelete, nil
}

func channelOverwriteCreateFromEntry(e *rawEntry) (*ChannelOverwriteCreate, error) {
	overwriteCreate := &ChannelOverwriteCreate{
		BaseEntry: baseEntryFromRaw(e),
		RoleName:  e.Options.RoleName,
	}

	var err error
	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyID:
			overwriteCreate.ID, err = stringValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyID, err)
			}

		case changeKeyType:
			overwriteCreate.Type, err = intValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyType, err)
			}

		case changeKeyAllow:
			var allowStr string
			allowStr, err = stringValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q (as string): %w", changeKeyAllow, err)
			}
			overwriteCreate.Allow, err = strconv.Atoi(allowStr)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyAllow, err)
			}

		case changeKeyDeny:
			var denyStr string
			denyStr, err = stringValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q (as string): %w", changeKeyAllow, err)
			}
			overwriteCreate.Deny, err = strconv.Atoi(denyStr)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyDeny, err)
			}
		}
	}

	return overwriteCreate, nil
}

func channelOverwriteUpdateFromEntry(e *rawEntry) (*ChannelOverwriteUpdate, error) {
	overwriteUpdate := &ChannelOverwriteUpdate{
		BaseEntry: baseEntryFromRaw(e),
		ID:        e.Options.ID,
		RoleName:  e.Options.RoleName,
	}

	var err error
	overwriteUpdate.Type, err = strconv.Atoi(e.Options.Type)
	if err != nil {
		return nil, fmt.Errorf("type: %w", err)
	}

	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyAllow:
			oldValueStr, newValueStr, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyAllow, err)
			}

			var oldValue, newValue int
			oldValue, err = strconv.Atoi(oldValueStr)
			if err != nil {
				return nil, fmt.Errorf("change key %q (old value): %w", changeKeyAllow, err)
			}
			newValue, err = strconv.Atoi(newValueStr)
			if err != nil {
				return nil, fmt.Errorf("change key %q (new value): %w", changeKeyAllow, err)
			}
			overwriteUpdate.Allow = &IntValues{Old: oldValue, New: newValue}

		case changeKeyDeny:
			oldValueStr, newValueStr, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyDeny, err)
			}

			var oldValue, newValue int
			oldValue, err = strconv.Atoi(oldValueStr)
			if err != nil {
				return nil, fmt.Errorf("change key %q (old value): %w", changeKeyDeny, err)
			}
			newValue, err = strconv.Atoi(newValueStr)
			if err != nil {
				return nil, fmt.Errorf("change key %q (new value): %w", changeKeyDeny, err)
			}
			overwriteUpdate.Deny = &IntValues{Old: oldValue, New: newValue}
		}
	}

	return overwriteUpdate, nil
}

func channelOverwriteDeleteFromEntry(e *rawEntry) (*ChannelOverwriteDelete, error) {
	overwriteDelete := &ChannelOverwriteDelete{
		BaseEntry: baseEntryFromRaw(e),
		RoleName:  e.Options.RoleName,
	}

	var err error
	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyID:
			overwriteDelete.ID, err = stringValue(ch.Old)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyID, err)
			}

		case changeKeyType:
			overwriteDelete.Type, err = intValue(ch.Old)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyType, err)
			}

		case changeKeyAllow:
			var allowStr string
			allowStr, err = stringValue(ch.Old)
			if err != nil {
				return nil, fmt.Errorf("change key %q (as string): %w", changeKeyAllow, err)
			}
			overwriteDelete.Allow, err = strconv.Atoi(allowStr)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyAllow, err)
			}

		case changeKeyDeny:
			var denyStr string
			denyStr, err = stringValue(ch.Old)
			if err != nil {
				return nil, fmt.Errorf("change key %q (as string): %w", changeKeyDeny, err)
			}
			overwriteDelete.Deny, err = strconv.Atoi(denyStr)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyDeny, err)
			}
		}
	}

	return overwriteDelete, nil
}

func messageDeleteFromEntry(e *rawEntry) (*MessageDelete, error) {
	msgDelete := &MessageDelete{
		BaseEntry: baseEntryFromRaw(e),
		ChannelID: e.Options.ChannelID,
	}

	var err error
	msgDelete.Count, err = strconv.Atoi(e.Options.Count)
	if err != nil {
		return nil, err
	}

	return msgDelete, nil
}
