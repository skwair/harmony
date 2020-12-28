package audit

import (
	"fmt"
)

func inviteCreateFromEntry(e *rawEntry) (*InviteCreate, error) {
	inviteCreate := &InviteCreate{
		BaseEntry: baseEntryFromRaw(e),
	}

	var err error
	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyCode:
			inviteCreate.Code, err = stringValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyCode, err)
			}

		case changeKeyChannelID:
			inviteCreate.ChannelID, err = stringValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyChannelID, err)
			}

		case changeKeyInviterID:
			inviteCreate.InviterID, err = stringValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyInviterID, err)
			}

		case changeKeyMaxUses:
			inviteCreate.MaxUses, err = intValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyMaxUses, err)
			}

		case changeKeyUses:
			inviteCreate.Uses, err = intValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyUses, err)
			}

		case changeKeyMaxAge:
			inviteCreate.MaxAge, err = intValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyMaxAge, err)
			}

		case changeKeyTemporary:
			inviteCreate.Temporary, err = boolValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyTemporary, err)
			}
		}
	}

	return inviteCreate, nil
}

func inviteUpdateFromEntry(e *rawEntry) (*InviteUpdate, error) {
	inviteUpdate := &InviteUpdate{
		BaseEntry: baseEntryFromRaw(e),
	}

	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyCode:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyCode, err)
			}
			inviteUpdate.Code = &StringValues{Old: oldValue, New: newValue}

		case changeKeyChannelID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyChannelID, err)
			}
			inviteUpdate.ChannelID = &StringValues{Old: oldValue, New: newValue}

		case changeKeyInviterID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyInviterID, err)
			}
			inviteUpdate.InviterID = &StringValues{Old: oldValue, New: newValue}

		case changeKeyMaxUses:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyMaxUses, err)
			}
			inviteUpdate.MaxUses = &IntValues{Old: oldValue, New: newValue}

		case changeKeyUses:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyUses, err)
			}
			inviteUpdate.Uses = &IntValues{Old: oldValue, New: newValue}

		case changeKeyMaxAge:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyMaxAge, err)
			}
			inviteUpdate.MaxAge = &IntValues{Old: oldValue, New: newValue}

		case changeKeyTemporary:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyTemporary, err)
			}
			inviteUpdate.Temporary = &BoolValues{Old: oldValue, New: newValue}
		}
	}

	return inviteUpdate, nil
}

func inviteDeleteFromEntry(e *rawEntry) (*InviteCreate, error) {
	inviteDelete := &InviteCreate{
		BaseEntry: baseEntryFromRaw(e),
	}

	var err error
	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyCode:
			inviteDelete.Code, err = stringValue(ch.Old)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyCode, err)
			}

		case changeKeyChannelID:
			inviteDelete.ChannelID, err = stringValue(ch.Old)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyChannelID, err)
			}

		case changeKeyInviterID:
			inviteDelete.InviterID, err = stringValue(ch.Old)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyInviterID, err)
			}

		case changeKeyMaxUses:
			inviteDelete.MaxUses, err = intValue(ch.Old)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyMaxUses, err)
			}

		case changeKeyUses:
			inviteDelete.Uses, err = intValue(ch.Old)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyUses, err)
			}

		case changeKeyMaxAge:
			inviteDelete.MaxAge, err = intValue(ch.Old)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyMaxAge, err)
			}

		case changeKeyTemporary:
			inviteDelete.Temporary, err = boolValue(ch.Old)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyTemporary, err)
			}
		}
	}

	return inviteDelete, nil
}
