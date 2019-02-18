package harmony

import "github.com/skwair/harmony/audit"

func roleCreateFromEntry(e *entry) (*audit.RoleCreate, error) {
	role := &audit.RoleCreate{
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

func roleUpdateFromEntry(e *entry) (*audit.RoleUpdate, error) {
	role := &audit.RoleUpdate{
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
			role.Name = &audit.StringValues{Old: oldValue, New: newValue}

		case changeKeyPermissions:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			role.Permissions = &audit.IntValues{Old: oldValue, New: newValue}

		case changeKeyColor:
			oldValue, newValue, err := intValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			role.Color = &audit.IntValues{Old: oldValue, New: newValue}

		case changeKeyHoist:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			role.Hoist = &audit.BoolValues{Old: oldValue, New: newValue}

		case changeKeyMentionable:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			role.Mentionable = &audit.BoolValues{Old: oldValue, New: newValue}
		}
	}

	return role, nil
}

func roleDeleteFromEntry(e *entry) (*audit.RoleDelete, error) {
	role := &audit.RoleDelete{
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
