package audit

import "encoding/json"

// rawAuditLog is the raw audit log, as returned by Discord's API.
type rawAuditLog struct {
	Entries []rawEntry `json:"audit_log_entries"`
}

type rawEntry struct {
	ID         string `json:"id"`
	ActionType int    `json:"action_type"`
	TargetID   string `json:"target_id"`
	UserID     string `json:"user_id"`
	Reason     string `json:"reason"`
	Changes    []struct {
		Key string          `json:"key"`
		Old json.RawMessage `json:"old_value"`
		New json.RawMessage `json:"new_value"`
	} `json:"changes"`
	Options struct {
		// MEMBER_PRUNE actions.
		DeleteMemberDays string `json:"delete_member_days"` // Number of days after which inactive members were kicked.
		MembersRemoved   string `json:"members_removed"`    // Number of members removed by the prune.

		// MESSAGE_DELETE actions.
		ChannelID string `json:"channel_id"` // ID of the channel in which the messages were deleted.
		Count     string `json:"count"`      // Number of deleted messages.

		// CHANNEL_OVERWRITE_* actions.
		ID       string `json:"id"`        // ID of the overwritten entity.
		Type     string `json:"type"`      // Type of the overwritten entity ("member" or "role").
		RoleName string `json:"role_name"` // Name of the role if Type is "role".
	} `json:"options"`
}

// ParseRaw parses a raw, JSON-encoded audit log returned by Discord's API
// into a typed and structured Log.
// It is not intended to be used by end user and is only exposed so
// harmony.AuditLog can use it.
func ParseRaw(raw json.RawMessage) (*Log, error) {
	var log rawAuditLog
	if err := json.Unmarshal(raw, &log); err != nil {
		return nil, err
	}

	var res Log

	// For every "raw" entry in this audit log, generate the
	// typed audit entry that corresponds to the action type.
	for _, e := range log.Entries {
		var (
			entry LogEntry
			err   error
		)

		switch EntryType(e.ActionType) {
		case EntryTypeGuildUpdate:
			entry, err = guildUpdateFromEntry(&e)

		case EntryTypeChannelCreate:
			entry, err = channelCreateFromEntry(&e)

		case EntryTypeChannelUpdate:
			entry, err = channelUpdateFromEntry(&e)

		case EntryTypeChannelDelete:
			entry, err = channelDeleteFromEntry(&e)

		case EntryTypeChannelOverwriteCreate:
			entry, err = channelOverwriteCreateFromEntry(&e)

		case EntryTypeChannelOverwriteUpdate:
			entry, err = channelOverwriteUpdateFromEntry(&e)

		case EntryTypeChannelOverwriteDelete:
			entry, err = channelOverwriteDeleteFromEntry(&e)

		case EntryTypeMemberKick:
			entry, err = memberKickFromEntry(&e)

		case EntryTypeMemberPrune:
			entry, err = memberPruneFromEntry(&e)

		case EntryTypeMemberBanAdd:
			entry, err = memberBanAddFromEntry(&e)

		case EntryTypeMemberBanRemove:
			entry, err = memberBanRemoveFromEntry(&e)

		case EntryTypeMemberUpdate:
			entry, err = memberUpdateFromEntry(&e)

		case EntryTypeMemberRoleUpdate:
			entry, err = memberRoleUpdateFromEntry(&e)

		case EntryTypeRoleCreate:
			entry, err = roleCreateFromEntry(&e)

		case EntryTypeRoleUpdate:
			entry, err = roleUpdateFromEntry(&e)

		case EntryTypeRoleDelete:
			entry, err = roleDeleteFromEntry(&e)

		case EntryTypeInviteCreate:
			entry, err = inviteCreateFromEntry(&e)

		case EntryTypeInviteUpdate:
			entry, err = inviteUpdateFromEntry(&e)

		case EntryTypeInviteDelete:
			entry, err = inviteDeleteFromEntry(&e)

		case EntryTypeWebhookCreate:
			entry, err = webhookCreateFromEntry(&e)

		case EntryTypeWebhookUpdate:
			entry, err = webhookUpdateFromEntry(&e)

		case EntryTypeWebhookDelete:
			entry, err = webhookDeleteFromEntry(&e)

		case EntryTypeEmojiCreate:
			entry, err = emojiCreateFromEntry(&e)

		case EntryTypeEmojiUpdate:
			entry, err = emojiUpdateFromEntry(&e)

		case EntryTypeEmojiDelete:
			entry, err = emojiDeleteFromEntry(&e)

		case EntryTypeMessageDelete:
			entry, err = messageDeleteFromEntry(&e)
		}

		if err != nil {
			return nil, err
		}
		res.Entries = append(res.Entries, entry)
	}

	return &res, nil
}
