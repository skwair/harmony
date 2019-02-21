package audit

func emojiCreateFromEntry(e *rawEntry) (*EmojiCreate, error) {
	emoji := &EmojiCreate{
		BaseEntry: BaseEntry{
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
			emoji.Name, err = stringValue(ch.New)
			if err != nil {
				return nil, err
			}
		}
	}

	return emoji, nil
}

func emojiUpdateFromEntry(e *rawEntry) (*EmojiUpdate, error) {
	emoji := &EmojiUpdate{
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
			emoji.Name = &StringValues{Old: oldValue, New: newValue}
		}
	}

	return emoji, nil
}

func emojiDeleteFromEntry(e *rawEntry) (*EmojiDelete, error) {
	emoji := &EmojiDelete{
		BaseEntry: BaseEntry{
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
			emoji.Name, err = stringValue(ch.Old)
			if err != nil {
				return nil, err
			}
		}
	}

	return emoji, nil
}
