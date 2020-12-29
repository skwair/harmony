package audit

import (
	"fmt"
)

func guildUpdateFromEntry(e *rawEntry) (*GuildUpdate, error) {
	guildUpdate := &GuildUpdate{
		BaseEntry: baseEntryFromRaw(e),
	}

	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyName:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyName, err)
			}
			guildUpdate.Name = &StringValues{Old: oldValue, New: newValue}

		case changeKeyIconHash:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyIconHash, err)
			}
			guildUpdate.IconHash = &StringValues{Old: oldValue, New: newValue}

		case changeKeySplashHash:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeySplashHash, err)
			}
			guildUpdate.SplashHash = &StringValues{Old: oldValue, New: newValue}

		case changeKeyOwnerID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyOwnerID, err)
			}
			guildUpdate.OwnerID = &StringValues{Old: oldValue, New: newValue}

		case changeKeyRegion:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyRegion, err)
			}
			guildUpdate.Region = &StringValues{Old: oldValue, New: newValue}

		case changeKeyAFKChannelID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyAFKChannelID, err)
			}
			guildUpdate.AFKChannelID = &StringValues{Old: oldValue, New: newValue}

		case changeKeyAFKTimeout:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyAFKTimeout, err)
			}
			guildUpdate.AFKTimeout = &IntValues{Old: oldValue, New: newValue}

		case changeKeyMFALevel:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyMFALevel, err)
			}
			guildUpdate.MFALevel = &IntValues{Old: oldValue, New: newValue}

		case changeKeyVerificationLevel:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyVerificationLevel, err)
			}
			guildUpdate.VerificationLevel = &IntValues{Old: oldValue, New: newValue}

		case changeKeyExplicitContentFilter:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyExplicitContentFilter, err)
			}
			guildUpdate.ExplicitContentFilter = &IntValues{Old: oldValue, New: newValue}

		case changeKeyDefaultMessageNotification:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyDefaultMessageNotification, err)
			}
			guildUpdate.DefaultMessageNotification = &IntValues{Old: oldValue, New: newValue}

		case changeKeyVanityURLCode:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyVanityURLCode, err)
			}
			guildUpdate.VanityURLCode = &StringValues{Old: oldValue, New: newValue}

		case changeKeyPruneDeleteDays:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyPruneDeleteDays, err)
			}
			guildUpdate.PruneDeleteDays = &IntValues{Old: oldValue, New: newValue}

		case changeKeyWidgetEnabled:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyWidgetEnabled, err)
			}
			guildUpdate.WidgetEnabled = &BoolValues{Old: oldValue, New: newValue}

		case changeKeyWidgetChannelID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyWidgetChannelID, err)
			}
			guildUpdate.WidgetChannelID = &StringValues{Old: oldValue, New: newValue}
		}
	}

	return guildUpdate, nil
}
