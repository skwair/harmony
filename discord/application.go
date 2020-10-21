package discord

// Application represents a Discord developer application.
type Application struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Icon                string   `json:"icon"`
	Description         string   `json:"description"`
	RPCOrigins          []string `json:"rpc_origins"`
	BotPublic           bool     `json:"bot_public"`
	BotRequireCodeGrant bool     `json:"bot_require_code_grant"`
	Owner               *User    `json:"owner"`
	Summary             string   `json:"summary"`
	VerifyKey           string   `json:"verify_key"`
	Team                Team     `json:"team"`
	GuildID             string   `json:"guild_id"`
	PrimarySKUID        string   `json:"primary_sku_id"`
	Slug                string   `json:"slug"`
	CoverImage          string   `json:"cover_image"`
}

// Team represents a Discord developer team.
type Team struct {
	ID          string       `json:"id"`
	Members     []TeamMember `json:"members"`
	OwnerUserID string       `json:"owner_member_id"`
	Icon        string       `json:"icon"`
}

// TeamMember is a member part of a Discord developer team.
type TeamMember struct {
	TeamID          string   `json:"team_id"`
	User            User     `json:"user"`
	MembershipState int      `json:"membership_state"`
	Permissions     []string `json:"permissions"`
}
