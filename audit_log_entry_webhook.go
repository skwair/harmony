package harmony

import "github.com/skwair/harmony/audit"

func webhookCreateFromEntry(e *entry) (*audit.WebhookCreate, error) {
	webhook := &audit.WebhookCreate{
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
			webhook.Name, err = stringValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyType:
			webhook.Type, err = intValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyChannelID:
			webhook.ChannelID, err = stringValue(ch.New)
			if err != nil {
				return nil, err
			}
		}
	}

	return webhook, nil
}

func webhookUpdateFromEntry(e *entry) (*audit.WebhookUpdate, error) {
	webhook := &audit.WebhookUpdate{
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
			webhook.Name = &audit.StringValues{Old: oldValue, New: newValue}

		case changeKeyChannelID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			webhook.ChannelID = &audit.StringValues{Old: oldValue, New: newValue}

		case changeKeyAvatarHash:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			webhook.AvatarHash = &audit.StringValues{Old: oldValue, New: newValue}
		}
	}

	return webhook, nil
}

func webhookDeleteFromEntry(e *entry) (*audit.WebhookDelete, error) {
	webhook := &audit.WebhookDelete{
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
			webhook.Name, err = stringValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyType:
			webhook.Type, err = intValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyChannelID:
			webhook.ChannelID, err = stringValue(ch.Old)
			if err != nil {
				return nil, err
			}
		}
	}

	return webhook, nil
}
