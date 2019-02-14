package harmony

import "net/http"

const auditLogReasonHeader = "X-Audit-Log-Reason"

func reasonHeader(r string) http.Header {
	h := http.Header{}

	if r != "" {
		h.Set(auditLogReasonHeader, r)
	}

	return h
}
