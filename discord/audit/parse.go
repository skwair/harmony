package audit

import (
	"encoding/json"
	"fmt"
)

// rawAuditLog is the raw audit log, as returned by Discord's API.
type rawAuditLog struct {
	Entries []rawEntry `json:"audit_log_entries"`
}

// rawEntry represents a single audit log entry, as returned by Discord's API.
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

func baseEntryFromRaw(e *rawEntry) BaseEntry {
	return BaseEntry{
		ID:       e.ID,
		TargetID: e.TargetID,
		UserID:   e.UserID,
		Reason:   e.Reason,
	}
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
			if err != nil {
				return nil, fmt.Errorf("guild update: %w", err)
			}

		case EntryTypeChannelCreate:
			entry, err = channelCreateFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("channel create: %w", err)
			}

		case EntryTypeChannelUpdate:
			entry, err = channelUpdateFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("channel update: %w", err)
			}

		case EntryTypeChannelDelete:
			entry, err = channelDeleteFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("channel delete: %w", err)
			}

		case EntryTypeChannelOverwriteCreate:
			entry, err = channelOverwriteCreateFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("channel overwrite create: %w", err)
			}

		case EntryTypeChannelOverwriteUpdate:
			entry, err = channelOverwriteUpdateFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("channel overwrite update: %w", err)
			}

		case EntryTypeChannelOverwriteDelete:
			entry, err = channelOverwriteDeleteFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("channel overwrite delete: %w", err)
			}

		case EntryTypeMemberKick:
			entry, err = memberKickFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("member kick: %w", err)
			}

		case EntryTypeMemberPrune:
			entry, err = memberPruneFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("member prune: %w", err)
			}

		case EntryTypeMemberBanAdd:
			entry, err = memberBanAddFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("member ban add: %w", err)
			}

		case EntryTypeMemberBanRemove:
			entry, err = memberBanRemoveFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("member ban remove: %w", err)
			}

		case EntryTypeMemberUpdate:
			entry, err = memberUpdateFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("member update: %w", err)
			}

		case EntryTypeMemberRoleUpdate:
			entry, err = memberRoleUpdateFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("member role update: %w", err)
			}

		case EntryTypeRoleCreate:
			entry, err = roleCreateFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("member role create: %w", err)
			}

		case EntryTypeRoleUpdate:
			entry, err = roleUpdateFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("role update: %w", err)
			}

		case EntryTypeRoleDelete:
			entry, err = roleDeleteFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("role delete: %w", err)
			}

		case EntryTypeInviteCreate:
			entry, err = inviteCreateFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("invite create: %w", err)
			}

		case EntryTypeInviteUpdate:
			entry, err = inviteUpdateFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("invite update: %w", err)
			}

		case EntryTypeInviteDelete:
			entry, err = inviteDeleteFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("invite delete: %w", err)
			}

		case EntryTypeWebhookCreate:
			entry, err = webhookCreateFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("webhook create: %w", err)
			}

		case EntryTypeWebhookUpdate:
			entry, err = webhookUpdateFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("webhook update: %w", err)
			}

		case EntryTypeWebhookDelete:
			entry, err = webhookDeleteFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("webhook delete: %w", err)
			}

		case EntryTypeEmojiCreate:
			entry, err = emojiCreateFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("emoji create: %w", err)
			}

		case EntryTypeEmojiUpdate:
			entry, err = emojiUpdateFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("emoji update: %w", err)
			}

		case EntryTypeEmojiDelete:
			entry, err = emojiDeleteFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("emoji delete: %w", err)
			}

		case EntryTypeMessageDelete:
			entry, err = messageDeleteFromEntry(&e)
			if err != nil {
				return nil, fmt.Errorf("message delete: %w", err)
			}
		}

		res.Entries = append(res.Entries, entry)
	}

	return &res, nil
}
