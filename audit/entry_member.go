package audit

import (
	"strconv"
)

func memberKickFromEntry(e *rawEntry) (*MemberKick, error) {
	kick := &MemberKick{
		BaseEntry: BaseEntry{
			ID:       e.ID,
			TargetID: e.TargetID,
			UserID:   e.UserID,
			Reason:   e.Reason,
		},
	}

	return kick, nil
}

func memberPruneFromEntry(e *rawEntry) (*MemberPrune, error) {
	prune := &MemberPrune{
		BaseEntry: BaseEntry{
			ID:       e.ID,
			TargetID: e.TargetID,
			UserID:   e.UserID,
			Reason:   e.Reason,
		},
	}

	var err error
	prune.DeleteMemberDays, err = strconv.Atoi(e.Options.DeleteMemberDays)
	if err != nil {
		return nil, err
	}

	prune.MembersRemoved, err = strconv.Atoi(e.Options.MembersRemoved)
	if err != nil {
		return nil, err
	}

	return prune, nil
}

func memberBanAddFromEntry(e *rawEntry) (*MemberBanAdd, error) {
	ban := &MemberBanAdd{
		BaseEntry: BaseEntry{
			ID:       e.ID,
			TargetID: e.TargetID,
			UserID:   e.UserID,
			Reason:   e.Reason,
		},
	}

	return ban, nil
}

func memberBanRemoveFromEntry(e *rawEntry) (*MemberBanRemove, error) {
	ban := &MemberBanRemove{
		BaseEntry: BaseEntry{
			ID:       e.ID,
			TargetID: e.TargetID,
			UserID:   e.UserID,
			Reason:   e.Reason,
		},
	}

	return ban, nil
}

func memberUpdateFromEntry(e *rawEntry) (*MemberUpdate, error) {
	member := &MemberUpdate{
		BaseEntry: BaseEntry{
			ID:       e.ID,
			TargetID: e.TargetID,
			UserID:   e.UserID,
			Reason:   e.Reason,
		},
	}

	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyNick:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			member.Nick = &StringValues{Old: oldValue, New: newValue}

		case changeKeyDeaf:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			member.Deaf = &BoolValues{Old: oldValue, New: newValue}

		case changeKeyMute:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			member.Mute = &BoolValues{Old: oldValue, New: newValue}
		}
	}

	return member, nil
}

func memberRoleUpdateFromEntry(e *rawEntry) (*MemberRoleUpdate, error) {
	member := &MemberRoleUpdate{
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
		case changeKeyAddRole:
			member.Added, err = permissionOverwritesValue(ch.New)
			if err != nil {
				return nil, err
			}

		case changeKeyRemoveRole:
			member.Removed, err = permissionOverwritesValue(ch.New)
			if err != nil {
				return nil, err
			}
		}
	}

	return member, nil
}
