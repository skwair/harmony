package audit

func roleCreateFromEntry(e *rawEntry) (*RoleCreate, error) {
	role := &RoleCreate{
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
			role.Name, err = stringValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyPermissions:
			role.Permissions, err = intValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyColor:
			role.Color, err = intValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyHoist:
			role.Hoist, err = boolValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyMentionable:
			role.Mentionable, err = boolValue(ch.New)
			if err != nil {
				return nil, err
			}
		}
	}

	return role, nil
}

func roleUpdateFromEntry(e *rawEntry) (*RoleUpdate, error) {
	role := &RoleUpdate{
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
			role.Name = &StringValues{Old: oldValue, New: newValue}

		case changeKeyPermissions:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			role.Permissions = &IntValues{Old: oldValue, New: newValue}

		case changeKeyColor:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			role.Color = &IntValues{Old: oldValue, New: newValue}

		case changeKeyHoist:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			role.Hoist = &BoolValues{Old: oldValue, New: newValue}

		case changeKeyMentionable:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			role.Mentionable = &BoolValues{Old: oldValue, New: newValue}
		}
	}

	return role, nil
}

func roleDeleteFromEntry(e *rawEntry) (*RoleDelete, error) {
	role := &RoleDelete{
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
			role.Name, err = stringValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyPermissions:
			role.Permissions, err = intValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyColor:
			role.Color, err = intValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyHoist:
			role.Hoist, err = boolValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyMentionable:
			role.Mentionable, err = boolValue(ch.Old)
			if err != nil {
				return nil, err
			}
		}
	}

	return role, nil
}
