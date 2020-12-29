package audit

import (
	"fmt"
	"strconv"
)

func memberKickFromEntry(e *rawEntry) (*MemberKick, error) {
	memberKick := &MemberKick{
		BaseEntry: baseEntryFromRaw(e),
	}

	return memberKick, nil
}

func memberPruneFromEntry(e *rawEntry) (*MemberPrune, error) {
	memberPrune := &MemberPrune{
		BaseEntry: baseEntryFromRaw(e),
	}

	var err error
	memberPrune.DeleteMemberDays, err = strconv.Atoi(e.Options.DeleteMemberDays)
	if err != nil {
		return nil, fmt.Errorf("delete member days: %w", err)
	}

	memberPrune.MembersRemoved, err = strconv.Atoi(e.Options.MembersRemoved)
	if err != nil {
		return nil, fmt.Errorf("members removed: %w", err)
	}

	return memberPrune, nil
}

func memberBanAddFromEntry(e *rawEntry) (*MemberBanAdd, error) {
	banAdd := &MemberBanAdd{
		BaseEntry: baseEntryFromRaw(e),
	}

	return banAdd, nil
}

func memberBanRemoveFromEntry(e *rawEntry) (*MemberBanRemove, error) {
	banRemove := &MemberBanRemove{
		BaseEntry: baseEntryFromRaw(e),
	}

	return banRemove, nil
}

func memberUpdateFromEntry(e *rawEntry) (*MemberUpdate, error) {
	memberUpdate := &MemberUpdate{
		BaseEntry: baseEntryFromRaw(e),
	}

	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyNick:
			oldValue, newValue, err := stringValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyNick, err)
			}
			memberUpdate.Nick = &StringValues{Old: oldValue, New: newValue}

		case changeKeyDeaf:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyDeaf, err)
			}
			memberUpdate.Deaf = &BoolValues{Old: oldValue, New: newValue}

		case changeKeyMute:
			oldValue, newValue, err := boolValues(ch.Old, ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyMute, err)
			}
			memberUpdate.Mute = &BoolValues{Old: oldValue, New: newValue}
		}
	}

	return memberUpdate, nil
}

func memberRoleUpdateFromEntry(e *rawEntry) (*MemberRoleUpdate, error) {
	roleUpdate := &MemberRoleUpdate{
		BaseEntry: baseEntryFromRaw(e),
	}

	var err error
	for _, ch := range e.Changes {
		switch changeKey(ch.Key) {
		case changeKeyAddRole:
			roleUpdate.Added, err = permissionOverwritesValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyAddRole, err)
			}

		case changeKeyRemoveRole:
			roleUpdate.Removed, err = permissionOverwritesValue(ch.New)
			if err != nil {
				return nil, fmt.Errorf("change key %q: %w", changeKeyRemoveRole, err)
			}
		}
	}

	return roleUpdate, nil
}
