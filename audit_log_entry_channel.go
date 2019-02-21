package harmony

import (
	"strconv"

	"github.com/skwair/harmony/audit"
)

func channelCreateFromEntry(e *entry) (*audit.ChannelCreate, error) {
	chCreate := &audit.ChannelCreate{
		BaseEntry: audit.BaseEntry{
			ID:       e.ID,
			TargetID: e.TargetID,
			UserID:   e.UserID,
			Reason:   e.Reason,
		},
	}

	var err error
	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyName:
			chCreate.Name, err = stringValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyType:
			chCreate.Type, err = intValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyRateLimitPerUser:
			chCreate.RateLimitPerUser, err = intValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyNFSW:
			chCreate.NSFW, err = boolValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyPermissionOverwrites:
			chCreate.PermissionOverwrites, err = permissionOverwritesValue(ch.New)
			if err != nil {
				return nil, err
			}
		}
	}

	return chCreate, nil
}

func channelUpdateFromEntry(e *entry) (*audit.ChannelUpdate, error) {
	chUpdate := &audit.ChannelUpdate{
		BaseEntry: audit.BaseEntry{
			ID:       e.ID,
			TargetID: e.TargetID,
			UserID:   e.UserID,
			Reason:   e.Reason,
		},
	}

	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyName:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			chUpdate.Name = &audit.StringValues{Old: oldValue, New: newValue}

		case changeKeyTopic:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			chUpdate.Topic = &audit.StringValues{Old: oldValue, New: newValue}

		case changeKeyBitrate:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			chUpdate.Bitrate = &audit.IntValues{Old: oldValue, New: newValue}

		case changeKeyRateLimitPerUser:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			chUpdate.RateLimitPerUser = &audit.IntValues{Old: oldValue, New: newValue}

		case changeKeyNFSW:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			chUpdate.NSFW = &audit.BoolValues{Old: oldValue, New: newValue}

		case changeKeyApplicationID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			chUpdate.ApplicationID = &audit.StringValues{Old: oldValue, New: newValue}

		}
	}

	return chUpdate, nil
}

func channelDeleteFromEntry(e *entry) (*audit.ChannelDelete, error) {
	chDelete := &audit.ChannelDelete{
		BaseEntry: audit.BaseEntry{
			ID:       e.ID,
			TargetID: e.TargetID,
			UserID:   e.UserID,
			Reason:   e.Reason,
		},
	}

	var err error
	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyName:
			chDelete.Name, err = stringValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyType:
			chDelete.Type, err = intValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyRateLimitPerUser:
			chDelete.RateLimitPerUser, err = intValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyNFSW:
			chDelete.NSFW, err = boolValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyPermissionOverwrites:
			chDelete.PermissionOverwrites, err = permissionOverwritesValue(ch.Old)
			if err != nil {
				return nil, err
			}
		}
	}

	return chDelete, nil
}

func channelOverwriteCreateFromEntry(e *entry) (*audit.ChannelOverwriteCreate, error) {
	overwrite := &audit.ChannelOverwriteCreate{
		BaseEntry: audit.BaseEntry{
			ID:       e.ID,
			TargetID: e.TargetID,
			UserID:   e.UserID,
			Reason:   e.Reason,
		},
		RoleName: e.Options.RoleName,
	}

	var err error
	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyID:
			overwrite.ID, err = stringValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyType:
			overwrite.Type, err = stringValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyAllow:
			overwrite.Allow, err = intValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyDeny:
			overwrite.Deny, err = intValue(ch.New)
			if err != nil {
				return nil, err
			}
		}
	}

	return overwrite, nil
}

func channelOverwriteUpdateFromEntry(e *entry) (*audit.ChannelOverwriteUpdate, error) {
	overwrite := &audit.ChannelOverwriteUpdate{
		BaseEntry: audit.BaseEntry{
			ID:       e.ID,
			TargetID: e.TargetID,
			UserID:   e.UserID,
			Reason:   e.Reason,
		},
		Type:     e.Options.Type,
		ID:       e.Options.ID,
		RoleName: e.Options.RoleName,
	}

	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyAllow:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			overwrite.Allow = &audit.IntValues{Old: oldValue, New: newValue}

		case changeKeyDeny:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			overwrite.Deny = &audit.IntValues{Old: oldValue, New: newValue}
		}
	}

	return overwrite, nil
}

func channelOverwriteDeleteFromEntry(e *entry) (*audit.ChannelOverwriteDelete, error) {
	overwrite := &audit.ChannelOverwriteDelete{
		BaseEntry: audit.BaseEntry{
			ID:       e.ID,
			TargetID: e.TargetID,
			UserID:   e.UserID,
			Reason:   e.Reason,
		},
		RoleName: e.Options.RoleName,
	}

	var err error
	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyID:
			overwrite.ID, err = stringValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyType:
			overwrite.Type, err = stringValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyAllow:
			overwrite.Allow, err = intValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyDeny:
			overwrite.Deny, err = intValue(ch.Old)
			if err != nil {
				return nil, err
			}
		}
	}

	return overwrite, nil
}

func messageDeleteFromEntry(e *entry) (*audit.MessageDelete, error) {
	message := &audit.MessageDelete{
		BaseEntry: audit.BaseEntry{
			ID:       e.ID,
			TargetID: e.TargetID,
			UserID:   e.UserID,
			Reason:   e.Reason,
		},
		ChannelID: e.Options.ChannelID,
	}

	var err error
	message.Count, err = strconv.Atoi(e.Options.Count)
	if err != nil {
		return nil, err
	}

	return message, nil
}
