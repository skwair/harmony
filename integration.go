package discord

import "time"

type Integration struct {
	ID                string              `json:"id"`
	Name              string              `json:"name"`
	Type              string              `json:"type"`
	Enabled           bool                `json:"enabled"`
	Syncing           bool                `json:"syncing"`
	RoleID            string              `json:"role_id"`
	ExpireBehavior    int                 `json:"expire_behavior"`
	ExpireGravePeriod int                 `json:"expire_grave_period"`
	User              *User               `json:"user"`
	Account           *IntegrationAccount `json:"account"`
	SyncedAt          time.Time           `json:"synced_at"`
}

type IntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
