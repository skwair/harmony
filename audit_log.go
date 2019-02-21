package harmony

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/skwair/harmony/audit"
	"github.com/skwair/harmony/internal/endpoint"
)

func reasonHeader(r string) http.Header {
	h := http.Header{}

	if r != "" {
		h.Set("X-Audit-Log-Reason", r)
	}

	return h
}

type auditLog struct {
	Entries []entry `json:"audit_log_entries"`
}

type entry struct {
	ID         string   `json:"id"`
	ActionType int      `json:"action_type"`
	TargetID   string   `json:"target_id"`
	UserID     string   `json:"user_id"`
	Reason     string   `json:"reason"`
	Changes    []change `json:"changes"`
	Options    options  `json:"options"`
}

type change struct {
	Key string          `json:"key"`
	Old json.RawMessage `json:"old_value"`
	New json.RawMessage `json:"new_value"`
}

type options struct {
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
}

// AuditLog returns the audit log of the given Guild. Requires the 'VIEW_AUDIT_LOG' permission.
func (c *Client) AuditLog(ctx context.Context, guildID, userID string, typ audit.EntryType, before string, limit int) (*audit.Log, error) {
	q := url.Values{}

	if userID != "" {
		q.Set("user_id", userID)
	}
	if typ != 0 {
		q.Set("action_type", strconv.Itoa(int(typ)))
	}
	if before != "" {
		q.Set("before", before)
	}
	if limit != 0 {
		q.Set("limit", strconv.Itoa(limit))
	}

	e := endpoint.GetAuditLog(guildID, q.Encode())
	resp, err := c.doReq(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var log auditLog
	if err = json.NewDecoder(resp.Body).Decode(&log); err != nil {
		return nil, err
	}

	return extractAuditLog(&log)
}

func extractAuditLog(log *auditLog) (*audit.Log, error) {
	var res audit.Log

	for _, e := range log.Entries {
		var (
			entry audit.LogEntry
			err   error
		)

		switch audit.EntryType(e.ActionType) {
		case audit.EntryTypeChannelCreate:
			entry, err = channelCreateFromEntry(&e)

		case audit.EntryTypeChannelUpdate:
			entry, err = channelUpdateFromEntry(&e)

		case audit.EntryTypeChannelDelete:
			entry, err = channelDeleteFromEntry(&e)

		case audit.EntryTypeChannelOverwriteCreate:
			entry, err = channelOverwriteCreateFromEntry(&e)

		case audit.EntryTypeChannelOverwriteUpdate:
			entry, err = channelOverwriteUpdateFromEntry(&e)

		case audit.EntryTypeChannelOverwriteDelete:
			entry, err = channelOverwriteDeleteFromEntry(&e)

		case audit.EntryTypeMemberKick:
			entry, err = memberKickFromEntry(&e)

		case audit.EntryTypeMemberPrune:
			entry, err = memberPruneFromEntry(&e)

		case audit.EntryTypeMemberBanAdd:
			entry, err = memberBanAddFromEntry(&e)

		case audit.EntryTypeMemberBanRemove:
			entry, err = memberBanRemoveFromEntry(&e)

		case audit.EntryTypeMemberUpdate:
			entry, err = memberUpdateFromEntry(&e)

		case audit.EntryTypeMemberRoleUpdate:
			entry, err = memberRoleUpdateFromEntry(&e)

		case audit.EntryTypeRoleCreate:
			entry, err = roleCreateFromEntry(&e)

		case audit.EntryTypeRoleUpdate:
			entry, err = roleUpdateFromEntry(&e)

		case audit.EntryTypeRoleDelete:
			entry, err = roleDeleteFromEntry(&e)

		case audit.EntryTypeInviteCreate:
			entry, err = inviteCreateFromEntry(&e)

		case audit.EntryTypeInviteUpdate:
			entry, err = inviteUpdateFromEntry(&e)

		case audit.EntryTypeInviteDelete:
			entry, err = inviteDeleteFromEntry(&e)

		case audit.EntryTypeWebhookCreate:
			entry, err = webhookCreateFromEntry(&e)

		case audit.EntryTypeWebhookUpdate:
			entry, err = webhookUpdateFromEntry(&e)

		case audit.EntryTypeWebhookDelete:
			entry, err = webhookDeleteFromEntry(&e)

		case audit.EntryTypeEmojiCreate:
			entry, err = emojiCreateFromEntry(&e)

		case audit.EntryTypeEmojiUpdate:
			entry, err = emojiUpdateFromEntry(&e)

		case audit.EntryTypeEmojiDelete:
			entry, err = emojiDeleteFromEntry(&e)

		case audit.EntryTypeMessageDelete:
			entry, err = messageDeleteFromEntry(&e)
		}

		if err != nil {
			return nil, err
		}
		res.Entries = append(res.Entries, entry)
	}

	return &res, nil
}
