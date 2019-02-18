package harmony

import (
	"strconv"

	"github.com/skwair/harmony/audit"
)

func memberKickFromEntry(e *entry) (*audit.MemberKick, error) {
	kick := &audit.MemberKick{
		BaseEntry: audit.BaseEntry{
			ID:       e.ID,
			TargetID: e.TargetID,
			UserID:   e.UserID,
			Reason:   e.Reason,
		},
	}

	return kick, nil
}

func memberPruneFromEntry(e *entry) (*audit.MemberPrune, error) {
	prune := &audit.MemberPrune{
		BaseEntry: audit.BaseEntry{
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

func memberBanAddFromEntry(e *entry) (*audit.MemberBanAdd, error) {
	ban := &audit.MemberBanAdd{
		BaseEntry: audit.BaseEntry{
			ID:       e.ID,
			TargetID: e.TargetID,
			UserID:   e.UserID,
			Reason:   e.Reason,
		},
	}

	return ban, nil
}

func memberBanRemoveFromEntry(e *entry) (*audit.MemberBanRemove, error) {
	ban := &audit.MemberBanRemove{
		BaseEntry: audit.BaseEntry{
			ID:       e.ID,
			TargetID: e.TargetID,
			UserID:   e.UserID,
			Reason:   e.Reason,
		},
	}

	return ban, nil
}

func memberUpdateFromEntry(e *entry) (*audit.MemberUpdate, error) {
	member := &audit.MemberUpdate{
		BaseEntry: audit.BaseEntry{
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
			member.Nick = &audit.StringValues{Old: oldValue, New: newValue}

		case changeKeyDeaf:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			member.Deaf = &audit.BoolValues{Old: oldValue, New: newValue}

		case changeKeyMute:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, err
			}
			member.Mute = &audit.BoolValues{Old: oldValue, New: newValue}
		}
	}

	return member, nil
}

func memberRoleUpdateFromEntry(e *entry) (*audit.MemberRoleUpdate, error) {
	member := &audit.MemberRoleUpdate{
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
