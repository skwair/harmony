package harmony

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/skwair/harmony/audit"
	"github.com/skwair/harmony/internal/endpoint"
)

type auditLogQuery struct {
	userID    string
	entryType audit.EntryType
	before    string
	limit     int
}

// AuditLogOption allows to customize a query to the audit log.
type AuditLogOption func(*auditLogQuery)

// WithUserID sets the user ID of the audit log query.
// It make the query only return audit log entries that have been
// creating for actions performed by this user.
func WithUserID(id string) AuditLogOption {
	return AuditLogOption(func(q *auditLogQuery) {
		q.userID = id
	})
}

// WithEntryType sets the entry type the query must return.
func WithEntryType(typ audit.EntryType) AuditLogOption {
	return AuditLogOption(func(q *auditLogQuery) {
		q.entryType = typ
	})
}

// WithBefore is used to paginate the audit log. The before parameter is the
// ID of the last audit log entry of our previous query.
func WithBefore(before string) AuditLogOption {
	return AuditLogOption(func(q *auditLogQuery) {
		q.before = before
	})
}

// WithLimit sets the limit the audit log query should return.
// It must be between 1 and 100 and defaults to 50 if not specified.
func WithLimit(limit int) AuditLogOption {
	return AuditLogOption(func(q *auditLogQuery) {
		q.limit = limit
	})
}

// AuditLog returns the audit log of the given Guild. Requires the 'VIEW_AUDIT_LOG' permission.
func (r *GuildResource) AuditLog(ctx context.Context, opts ...AuditLogOption) (*audit.Log, error) {
	query := &auditLogQuery{}

	for _, opt := range opts {
		opt(query)
	}

	q := url.Values{}

	if query.userID != "" {
		q.Set("user_id", query.userID)
	}
	if query.entryType != 0 {
		q.Set("action_type", strconv.Itoa(int(query.entryType)))
	}
	if query.before != "" {
		q.Set("before", query.before)
	}
	if query.limit != 0 {
		q.Set("limit", strconv.Itoa(query.limit))
	}

	e := endpoint.GetAuditLog(r.guildID, q.Encode())
	resp, err := r.client.doReq(ctx, e, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return audit.ParseRaw(b)
}

func reasonHeader(r string) http.Header {
	h := http.Header{}

	if r != "" {
		h.Set("X-Audit-Log-Reason", r)
	}

	return h
}
