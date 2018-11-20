package harmony

// Connection that the user has attached.
type Connection struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	Revoked      bool          `json:"revoked"`
	Integrations []Integration `json:"integrations"` // Partial server integrations.
}
