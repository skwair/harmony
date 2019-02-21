package harmony

import "github.com/skwair/harmony/audit"

func guildUpdateFromEntry(e *entry) (*audit.GuildUpdate, error) {
	guild := &audit.GuildUpdate{
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
			guild.Name = &audit.StringValues{Old: oldValue, New: newValue}

		case changeKeyIconHash:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.IconHash = &audit.StringValues{Old: oldValue, New: newValue}

		case changeKeySplashHash:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.SplashHash = &audit.StringValues{Old: oldValue, New: newValue}

		case changeKeyOwnerID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.OwnerID = &audit.StringValues{Old: oldValue, New: newValue}

		case changeKeyRegion:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.Region = &audit.StringValues{Old: oldValue, New: newValue}

		case changeKeyAFKChannelID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.AFKChannelID = &audit.StringValues{Old: oldValue, New: newValue}

		case changeKeyAFKTimeout:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.AFKTimeout = &audit.IntValues{Old: oldValue, New: newValue}

		case changeKeyMFALevel:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.MFALevel = &audit.IntValues{Old: oldValue, New: newValue}

		case changeKeyVerificationLevel:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.VerificationLevel = &audit.IntValues{Old: oldValue, New: newValue}

		case changeKeyExplicitContentFilter:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.ExplicitContentFilter = &audit.IntValues{Old: oldValue, New: newValue}

		case changeKeyDefaultMessageNotification:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.DefaultMessageNotification = &audit.IntValues{Old: oldValue, New: newValue}

		case changeKeyVanityURLCode:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.VanityURLCode = &audit.StringValues{Old: oldValue, New: newValue}

		case changeKeyPruneDeleteDays:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.PruneDeleteDays = &audit.IntValues{Old: oldValue, New: newValue}

		case changeKeyWidgetEnabled:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.WidgetEnabled = &audit.BoolValues{Old: oldValue, New: newValue}

		case changeKeyWidgetChannelID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.WidgetChannelID = &audit.StringValues{Old: oldValue, New: newValue}
		}
	}

	return guild, nil
}
