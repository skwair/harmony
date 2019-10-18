package harmony

// Clone returns a clone of this User.
func (u *User) Clone() *User {
	if u == nil {
		return nil
	}

	return &User{
		ID:            u.ID,
		Username:      u.Username,
		Discriminator: u.Discriminator,
		Avatar:        u.Avatar,
		Bot:           u.Bot,
		MFAEnabled:    u.MFAEnabled,
		Verified:      u.Verified,
		Email:         u.Email,
	}
}

// Clone returns a clone of this Guild.
func (g *Guild) Clone() *Guild {
	if g == nil {
		return nil
	}

	guild := &Guild{
		ID:                          g.ID,
		Name:                        g.Name,
		Splash:                      g.Icon,
		Owner:                       g.Owner,
		OwnerID:                     g.OwnerID,
		Permissions:                 g.Permissions,
		Region:                      g.Region,
		AFKChannelID:                g.AFKChannelID,
		AFKTimeout:                  g.AFKTimeout,
		EmbedEnabled:                g.EmbedEnabled,
		EmbedChannelID:              g.EmbedChannelID,
		VerificationLevel:           g.VerificationLevel,
		DefaultMessageNotifications: g.DefaultMessageNotifications,
		ExplicitContentFilter:       g.ExplicitContentFilter,
		MFALevel:                    g.MFALevel,
		ApplicationID:               g.ApplicationID,
		WidgetEnabled:               g.WidgetEnabled,
		WidgetChannelID:             g.WidgetChannelID,
		SystemChannelID:             g.SystemChannelID,
		JoinedAt:                    g.JoinedAt,
		Large:                       g.Large,
		Unavailable:                 g.Unavailable,
		MemberCount:                 g.MemberCount,
	}

	for i := 0; i < len(g.Roles); i++ {
		role := g.Roles[i].Clone()
		guild.Roles = append(guild.Roles, *role)
	}

	for i := 0; i < len(g.Emojis); i++ {
		emoji := g.Emojis[i].Clone()
		guild.Emojis = append(guild.Emojis, *emoji)
	}

	for i := 0; i < len(g.VoiceStates); i++ {
		vs := g.VoiceStates[i].Clone()
		guild.VoiceStates = append(guild.VoiceStates, *vs)
	}

	for i := 0; i < len(g.Members); i++ {
		member := g.Members[i].Clone()
		guild.Members = append(guild.Members, *member)
	}

	for i := 0; i < len(g.Channels); i++ {
		ch := g.Channels[i].Clone()
		guild.Channels = append(guild.Channels, *ch)
	}

	for i := 0; i < len(g.Presences); i++ {
		presence := g.Presences[i].Clone()
		guild.Presences = append(guild.Presences, *presence)
	}

	guild.Features = append(guild.Features, g.Features...)

	return guild
}

// Clone returns a clone of this Role.
func (r *Role) Clone() *Role {
	if r == nil {
		return nil
	}

	return &Role{
		ID:          r.ID,
		Name:        r.Name,
		Color:       r.Color,
		Hoist:       r.Hoist,
		Position:    r.Position,
		Permissions: r.Permissions,
		Managed:     r.Managed,
		Mentionable: r.Mentionable,
	}
}

// Clone returns a clone of this Emoji.
func (e *Emoji) Clone() *Emoji {
	if e == nil {
		return nil
	}

	emoji := &Emoji{
		ID:            e.ID,
		Name:          e.Name,
		RequireColons: e.RequireColons,
		Managed:       e.Managed,
		Animated:      e.Animated,
	}

	for i := 0; i < len(e.Roles); i++ {
		role := e.Roles[i].Clone()
		emoji.Roles = append(emoji.Roles, *role)
	}

	emoji.User = e.User.Clone()

	return emoji
}

// Clone returns a clone of this GuildMember.
func (m *GuildMember) Clone() *GuildMember {
	if m == nil {
		return nil
	}

	return &GuildMember{
		User:     m.User,
		Nick:     m.Nick,
		Roles:    m.Roles,
		JoinedAt: m.JoinedAt,
		Deaf:     m.Deaf,
		Mute:     m.Mute,
	}
}

// Clone returns a clone of this Channel.
func (c *Channel) Clone() *Channel {
	if c == nil {
		return nil
	}

	channel := &Channel{
		ID:               c.ID,
		Type:             c.Type,
		GuildID:          c.GuildID,
		Position:         c.Position,
		Name:             c.Name,
		Topic:            c.Topic,
		NSFW:             c.NSFW,
		LastMessageID:    c.LastMessageID,
		Bitrate:          c.Bitrate,
		UserLimit:        c.UserLimit,
		Icon:             c.Icon,
		OwnerID:          c.OwnerID,
		ApplicationID:    c.ApplicationID,
		ParentID:         c.ParentID,
		LastPinTimestamp: c.LastPinTimestamp,
	}

	for i := 0; i < len(c.PermissionOverwrites); i++ {
		overwrite := c.PermissionOverwrites[i].Clone()
		channel.PermissionOverwrites = append(channel.PermissionOverwrites, *overwrite)
	}

	for i := 0; i < len(c.Recipients); i++ {
		recipient := c.Recipients[i].Clone()
		channel.Recipients = append(channel.Recipients, *recipient)
	}

	return channel
}

// Clone returns a clone of this Presence.
func (p *Presence) Clone() *Presence {
	if p == nil {
		return nil
	}

	presence := &Presence{
		User:    p.User,
		Game:    p.Game,
		GuildID: p.GuildID,
		Status:  p.Status,
	}

	presence.Roles = append(presence.Roles, p.Roles...)

	return presence
}

// Clone returns a clone of this UnavailableGuild.
func (g *UnavailableGuild) Clone() *UnavailableGuild {
	if g == nil {
		return nil
	}

	return &UnavailableGuild{
		ID:          g.ID,
		Unavailable: g.Unavailable,
	}
}
