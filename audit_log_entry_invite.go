package harmony

import "github.com/skwair/harmony/audit"

func inviteCreateFromEntry(e *entry) (*audit.InviteCreate, error) {
	invite := &audit.InviteCreate{
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
		case changeKeyCode:
			invite.Code, err = stringValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyChannelID:
			invite.ChannelID, err = stringValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyInviterID:
			invite.InviterID, err = stringValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyMaxUses:
			invite.MaxUses, err = intValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyUses:
			invite.Uses, err = intValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyMaxAge:
			invite.MaxAge, err = intValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyTemporary:
			invite.Temporary, err = boolValue(ch.New)
			if err != nil {
				return nil, err
			}
		}
	}

	return invite, nil
}

func inviteUpdateFromEntry(e *entry) (*audit.InviteUpdate, error) {
	invite := &audit.InviteUpdate{
		BaseEntry: audit.BaseEntry{
			ID:       e.ID,
			TargetID: e.TargetID,
			UserID:   e.UserID,
			Reason:   e.Reason,
		},
	}

	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyCode:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			invite.Code = &audit.StringValues{Old: oldValue, New: newValue}

		case changeKeyChannelID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			invite.ChannelID = &audit.StringValues{Old: oldValue, New: newValue}

		case changeKeyInviterID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			invite.InviterID = &audit.StringValues{Old: oldValue, New: newValue}

		case changeKeyMaxUses:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			invite.MaxUses = &audit.IntValues{Old: oldValue, New: newValue}

		case changeKeyUses:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			invite.Uses = &audit.IntValues{Old: oldValue, New: newValue}

		case changeKeyMaxAge:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			invite.MaxAge = &audit.IntValues{Old: oldValue, New: newValue}

		case changeKeyTemporary:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			invite.Temporary = &audit.BoolValues{Old: oldValue, New: newValue}
		}
	}

	return invite, nil
}

func inviteDeleteFromEntry(e *entry) (*audit.InviteCreate, error) {
	invite := &audit.InviteCreate{
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
		case changeKeyCode:
			invite.Code, err = stringValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyChannelID:
			invite.ChannelID, err = stringValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyInviterID:
			invite.InviterID, err = stringValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyMaxUses:
			invite.MaxUses, err = intValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyUses:
			invite.Uses, err = intValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyMaxAge:
			invite.MaxAge, err = intValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyTemporary:
			invite.Temporary, err = boolValue(ch.Old)
			if err != nil {
				return nil, err
			}
		}
	}

	return invite, nil
}
