package audit

func webhookCreateFromEntry(e *rawEntry) (*WebhookCreate, error) {
	webhookCreate := &WebhookCreate{
		BaseEntry: baseEntryFromRaw(e),
	}

	var err error
	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyName:
			webhookCreate.Name, err = stringValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyType:
			webhookCreate.Type, err = intValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyChannelID:
			webhookCreate.ChannelID, err = stringValue(ch.New)
			if err != nil {
				return nil, err
			}
		}
	}

	return webhookCreate, nil
}

func webhookUpdateFromEntry(e *rawEntry) (*WebhookUpdate, error) {
	webhookUpdate := &WebhookUpdate{
		BaseEntry: baseEntryFromRaw(e),
	}

	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyName:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			webhookUpdate.Name = &StringValues{Old: oldValue, New: newValue}

		case changeKeyChannelID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			webhookUpdate.ChannelID = &StringValues{Old: oldValue, New: newValue}

		case changeKeyAvatarHash:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			webhookUpdate.AvatarHash = &StringValues{Old: oldValue, New: newValue}
		}
	}

	return webhookUpdate, nil
}

func webhookDeleteFromEntry(e *rawEntry) (*WebhookDelete, error) {
	webhookDelete := &WebhookDelete{
		BaseEntry: baseEntryFromRaw(e),
	}

	var err error
	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyName:
			webhookDelete.Name, err = stringValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyType:
			webhookDelete.Type, err = intValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyChannelID:
			webhookDelete.ChannelID, err = stringValue(ch.Old)
			if err != nil {
				return nil, err
			}
		}
	}

	return webhookDelete, nil
}
