package audit

func emojiCreateFromEntry(e *rawEntry) (*EmojiCreate, error) {
	emojiCreate := &EmojiCreate{
		BaseEntry: baseEntryFromRaw(e),
	}

	var err error
	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyName:
			emojiCreate.Name, err = stringValue(ch.New)
			if err != nil {
				return nil, err
			}
		}
	}

	return emojiCreate, nil
}

func emojiUpdateFromEntry(e *rawEntry) (*EmojiUpdate, error) {
	emojiUpdate := &EmojiUpdate{
		BaseEntry: baseEntryFromRaw(e),
	}

	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyName:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			emojiUpdate.Name = &StringValues{Old: oldValue, New: newValue}
		}
	}

	return emojiUpdate, nil
}

func emojiDeleteFromEntry(e *rawEntry) (*EmojiDelete, error) {
	emojiDelete := &EmojiDelete{
		BaseEntry: baseEntryFromRaw(e),
	}

	var err error
	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyName:
			emojiDelete.Name, err = stringValue(ch.Old)
			if err != nil {
				return nil, err
			}
		}
	}

	return emojiDelete, nil
}
