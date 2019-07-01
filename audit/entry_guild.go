package audit

func guildUpdateFromEntry(e *rawEntry) (*GuildUpdate, error) {
	guildUpdate := &GuildUpdate{
		BaseEntry: baseEntryFromRaw(e),
	}

	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyName:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guildUpdate.Name = &StringValues{Old: oldValue, New: newValue}

		case changeKeyIconHash:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guildUpdate.IconHash = &StringValues{Old: oldValue, New: newValue}

		case changeKeySplashHash:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guildUpdate.SplashHash = &StringValues{Old: oldValue, New: newValue}

		case changeKeyOwnerID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guildUpdate.OwnerID = &StringValues{Old: oldValue, New: newValue}

		case changeKeyRegion:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guildUpdate.Region = &StringValues{Old: oldValue, New: newValue}

		case changeKeyAFKChannelID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guildUpdate.AFKChannelID = &StringValues{Old: oldValue, New: newValue}

		case changeKeyAFKTimeout:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guildUpdate.AFKTimeout = &IntValues{Old: oldValue, New: newValue}

		case changeKeyMFALevel:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guildUpdate.MFALevel = &IntValues{Old: oldValue, New: newValue}

		case changeKeyVerificationLevel:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guildUpdate.VerificationLevel = &IntValues{Old: oldValue, New: newValue}

		case changeKeyExplicitContentFilter:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guildUpdate.ExplicitContentFilter = &IntValues{Old: oldValue, New: newValue}

		case changeKeyDefaultMessageNotification:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guildUpdate.DefaultMessageNotification = &IntValues{Old: oldValue, New: newValue}

		case changeKeyVanityURLCode:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guildUpdate.VanityURLCode = &StringValues{Old: oldValue, New: newValue}

		case changeKeyPruneDeleteDays:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guildUpdate.PruneDeleteDays = &IntValues{Old: oldValue, New: newValue}

		case changeKeyWidgetEnabled:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guildUpdate.WidgetEnabled = &BoolValues{Old: oldValue, New: newValue}

		case changeKeyWidgetChannelID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			guildUpdate.WidgetChannelID = &StringValues{Old: oldValue, New: newValue}
		}
	}

	return guildUpdate, nil
}
