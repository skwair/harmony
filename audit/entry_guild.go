package audit

func guildUpdateFromEntry(e *rawEntry) (*GuildUpdate, error) {
	guild := &GuildUpdate{
		BaseEntry: BaseEntry{
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
			guild.Name = &StringValues{Old: oldValue, New: newValue}

		case changeKeyIconHash:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.IconHash = &StringValues{Old: oldValue, New: newValue}

		case changeKeySplashHash:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.SplashHash = &StringValues{Old: oldValue, New: newValue}

		case changeKeyOwnerID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.OwnerID = &StringValues{Old: oldValue, New: newValue}

		case changeKeyRegion:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.Region = &StringValues{Old: oldValue, New: newValue}

		case changeKeyAFKChannelID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.AFKChannelID = &StringValues{Old: oldValue, New: newValue}

		case changeKeyAFKTimeout:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.AFKTimeout = &IntValues{Old: oldValue, New: newValue}

		case changeKeyMFALevel:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.MFALevel = &IntValues{Old: oldValue, New: newValue}

		case changeKeyVerificationLevel:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.VerificationLevel = &IntValues{Old: oldValue, New: newValue}

		case changeKeyExplicitContentFilter:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.ExplicitContentFilter = &IntValues{Old: oldValue, New: newValue}

		case changeKeyDefaultMessageNotification:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.DefaultMessageNotification = &IntValues{Old: oldValue, New: newValue}

		case changeKeyVanityURLCode:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.VanityURLCode = &StringValues{Old: oldValue, New: newValue}

		case changeKeyPruneDeleteDays:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.PruneDeleteDays = &IntValues{Old: oldValue, New: newValue}

		case changeKeyWidgetEnabled:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.WidgetEnabled = &BoolValues{Old: oldValue, New: newValue}

		case changeKeyWidgetChannelID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guild.WidgetChannelID = &StringValues{Old: oldValue, New: newValue}
		}
	}

	return guild, nil
}
