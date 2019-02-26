package audit

func roleCreateFromEntry(e *rawEntry) (*RoleCreate, error) {
	roleCreate := &RoleCreate{
		BaseEntry: baseEntryFromRaw(e),
	}

	var err error
	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyName:
			roleCreate.Name, err = stringValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyPermissions:
			roleCreate.Permissions, err = intValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyColor:
			roleCreate.Color, err = intValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyHoist:
			roleCreate.Hoist, err = boolValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyMentionable:
			roleCreate.Mentionable, err = boolValue(ch.New)
			if err != nil {
				return nil, err
			}
		}
	}

	return roleCreate, nil
}

func roleUpdateFromEntry(e *rawEntry) (*RoleUpdate, error) {
	roleUpdate := &RoleUpdate{
		BaseEntry: baseEntryFromRaw(e),
	}

	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyName:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			roleUpdate.Name = &StringValues{Old: oldValue, New: newValue}

		case changeKeyPermissions:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			roleUpdate.Permissions = &IntValues{Old: oldValue, New: newValue}

		case changeKeyColor:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			roleUpdate.Color = &IntValues{Old: oldValue, New: newValue}

		case changeKeyHoist:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			roleUpdate.Hoist = &BoolValues{Old: oldValue, New: newValue}

		case changeKeyMentionable:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			roleUpdate.Mentionable = &BoolValues{Old: oldValue, New: newValue}
		}
	}

	return roleUpdate, nil
}

func roleDeleteFromEntry(e *rawEntry) (*RoleDelete, error) {
	roleDelete := &RoleDelete{
		BaseEntry: baseEntryFromRaw(e),
	}

	var err error
	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyName:
			roleDelete.Name, err = stringValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyPermissions:
			roleDelete.Permissions, err = intValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyColor:
			roleDelete.Color, err = intValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyHoist:
			roleDelete.Hoist, err = boolValue(ch.Old)
			if err != nil {
				return nil, err
			}

		case changeKeyMentionable:
			roleDelete.Mentionable, err = boolValue(ch.Old)
			if err != nil {
				return nil, err
			}
		}
	}

	return roleDelete, nil
}
