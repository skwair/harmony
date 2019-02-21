package harmony

import "github.com/skwair/harmony/audit"

func emojiCreateFromEntry(e *entry) (*audit.EmojiCreate, error) {
	emoji := &audit.EmojiCreate{
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
			emoji.Name, err = stringValue(ch.New)
			if err != nil {
				return nil, err
			}
		}
	}

	return emoji, nil
}

func emojiUpdateFromEntry(e *entry) (*audit.EmojiUpdate, error) {
	emoji := &audit.EmojiUpdate{
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
			emoji.Name = &audit.StringValues{Old: oldValue, New: newValue}
		}
	}

	return emoji, nil
}

func emojiDeleteFromEntry(e *entry) (*audit.EmojiDelete, error) {
	emoji := &audit.EmojiDelete{
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
			emoji.Name, err = stringValue(ch.Old)
			if err != nil {
				return nil, err
			}
		}
	}

	return emoji, nil
}
