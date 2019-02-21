package audit

func inviteCreateFromEntry(e *rawEntry) (*InviteCreate, error) {
	invite := &InviteCreate{
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

func inviteUpdateFromEntry(e *rawEntry) (*InviteUpdate, error) {
	invite := &InviteUpdate{
		BaseEntry: BaseEntry{
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
			invite.Code = &StringValues{Old: oldValue, New: newValue}

		case changeKeyChannelID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			invite.ChannelID = &StringValues{Old: oldValue, New: newValue}

		case changeKeyInviterID:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			invite.InviterID = &StringValues{Old: oldValue, New: newValue}

		case changeKeyMaxUses:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			invite.MaxUses = &IntValues{Old: oldValue, New: newValue}

		case changeKeyUses:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			invite.Uses = &IntValues{Old: oldValue, New: newValue}

		case changeKeyMaxAge:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			invite.MaxAge = &IntValues{Old: oldValue, New: newValue}

		case changeKeyTemporary:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			invite.Temporary = &BoolValues{Old: oldValue, New: newValue}
		}
	}

	return invite, nil
}

func inviteDeleteFromEntry(e *rawEntry) (*InviteCreate, error) {
	invite := &InviteCreate{
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
