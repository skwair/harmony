package discord

import (
	"sync"
	"time"

	"github.com/skwair/discord/channel"
)

// State is a cache of the state of the application that is updated in real-time
// as events are received from the Gateway.
// The state must be locked by calling its RLock method before being used, else data
// races can occur. It must then be unlocked with the RUnlock method so it can continue
// to be updated as events flows in.
type State struct {
	mu sync.RWMutex

	currentUser       *User
	users             map[string]*User
	guilds            map[string]*Guild
	presences         map[string]*Presence // Presence by user ID.
	channels          map[string]*Channel
	dms               map[string]*Channel
	groups            map[string]*Channel
	unavailableGuilds map[string]*UnavailableGuild

	rtt time.Duration

	// NOTE: consider adding statistics such as the uptime, ping, number
	// of voice connections, etc... in the state.
}

// newState returns a new initialized state, ready to be used.
func newState() *State {
	return &State{
		users:             make(map[string]*User),
		guilds:            make(map[string]*Guild),
		presences:         make(map[string]*Presence),
		channels:          make(map[string]*Channel),
		dms:               make(map[string]*Channel),
		groups:            make(map[string]*Channel),
		unavailableGuilds: make(map[string]*UnavailableGuild),
	}
}

// RLock locks the state for reading, preventing it from being updated.
func (s *State) RLock() {
	s.mu.RLock()
}

// RUnlock releases the read lock on the state, allowing it to be updated
// as events are received from the gateway.
func (s *State) RUnlock() {
	s.mu.RUnlock()
}

// CurrentUser returns the current user from the state.
func (s *State) CurrentUser() *User {
	return s.currentUser
}

// User returns a user given its ID from the state.
func (s *State) User(id string) *User {
	return s.users[id]
}

// Guild returns a guild given its ID from the state.
func (s *State) Guild(id string) *Guild {
	return s.guilds[id]
}

// Channel returns a channel given its ID from the state.
func (s *State) Channel(id string) *Channel {
	return s.channels[id]
}

// GroupDM returns a group DM given its ID from the state.
func (s *State) GroupDM(id string) *Channel {
	return s.groups[id]
}

// DM returns a DM given its ID from the state.
func (s *State) DM(id string) *Channel {
	return s.dms[id]
}

// Presence returns a presence given a user ID from the state.
func (s *State) Presence(userID string) *Presence {
	return s.presences[userID]
}

// UnavailableGuild returns an unavailable guild given its ID from the state.
func (s *State) UnavailableGuild(id string) *UnavailableGuild {
	return s.unavailableGuilds[id]
}

// Users returns a map of user ID to user from the state.
func (s *State) Users() map[string]*User {
	return s.users
}

// Guilds returns a map of guild ID to guild from the state.
func (s *State) Guilds() map[string]*Guild {
	return s.guilds
}

// Channels returns a map of channels ID to channels from the state.
func (s *State) Channels() map[string]*Channel {
	return s.channels
}

// GroupDMs returns a map of group DM ID to group DM from the state.
func (s *State) GroupDMs() map[string]*Channel {
	return s.groups
}

// DMs returns a map of DM ID to DM from the state.
func (s *State) DMs() map[string]*Channel {
	return s.dms
}

// Presences returns a map of user ID to presence from the state.
func (s *State) Presences() map[string]*Presence {
	return s.presences
}

// UnavailableGuilds returns a map of guild ID to unavailable guild from the state.
func (s *State) UnavailableGuilds() map[string]*UnavailableGuild {
	return s.unavailableGuilds
}

// RTT returns the Round Trip Time between the client and Discord's Gateway.
// It is calculated and updated when sending heartbeat payloads (roughly
// every minute).
func (s *State) RTT() time.Duration {
	s.mu.RLock()
	defer s.mu.Unlock()

	return s.rtt
}

// setInitialState initializes the state with a Ready event received from the gateway.
func (s *State) setInitialState(r *Ready) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.currentUser = r.User
	for i := 0; i < len(r.Guilds); i++ {
		g := &r.Guilds[i]
		s.guilds[g.ID] = &Guild{
			ID:          g.ID,
			Name:        g.Name,
			Owner:       g.Owner,
			Permissions: g.Permissions,
		}
		if g.Icon != "" {
			s.guilds[g.ID].Icon = &g.Icon
		}
	}
	for i := 0; i < len(r.PrivateChannels); i++ {
		dm := &r.PrivateChannels[i]
		if dm.Type == channel.TypeDM {
			s.dms[dm.ID] = dm
		}
		if dm.Type == channel.TypeGroupDM {
			s.groups[dm.ID] = dm
		}
	}
}

// updateGuild adds the given guild to the state. If it already
// exists, it merges its content with the existing guild.
// It also removes this guild from the UnavailableGuilds map if
// it was present.
func (s *State) updateGuild(g *Guild) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Make sure we do not overwrite fields
	// that were set before but not anymore.
	old := s.guilds[g.ID]
	if old != nil {
		if g.Roles == nil {
			g.Roles = old.Roles
		}
		if g.Emojis == nil {
			g.Emojis = old.Emojis
		}
		if g.VoiceStates == nil {
			g.VoiceStates = old.VoiceStates
		}
		if g.Members == nil {
			g.Members = old.Members
		}
		if g.Channels == nil {
			g.Channels = old.Channels
		}
		if g.Presences == nil {
			g.Presences = old.Presences
		}
	}

	for i := 0; i < len(g.Channels); i++ {
		ch := &g.Channels[i]
		ch.GuildID = g.ID
		s.channels[ch.ID] = ch
	}

	for i := 0; i < len(g.Members); i++ {
		m := &g.Members[i]
		s.users[m.User.ID] = m.User
	}

	for i := 0; i < len(g.Presences); i++ {
		p := &g.Presences[i]
		s.presences[p.User.ID] = p
	}

	s.guilds[g.ID] = g
	delete(s.unavailableGuilds, g.ID)
}

// removeGuild removes a guild from the Guilds map, adding it to
// the UnavailableGuilds map.
func (s *State) removeGuild(g *UnavailableGuild) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.guilds, g.ID)
	s.unavailableGuilds[g.ID] = g
}

// updateGuildEmojis updates the emojis available in a guild if it
// is already tracked by the state, does nothing otherwise.
func (s *State) updateGuildEmojis(guildID string, emojis []Emoji) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.guilds[guildID] != nil {
		s.guilds[guildID].Emojis = emojis
	}
}

// updatePresence updates a presence both in the presences map as
// well as in the Guilds map.
func (s *State) updatePresence(p *Presence) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check that the concerned user exists in the state.
	if s.users[p.User.ID] == nil {
		return
	}

	// NOTE: consider removing the presence from the presence map
	// if the user goes offline.
	s.presences[p.User.ID] = p

	// Check that his guild exists in the state.
	if s.guilds[p.GuildID] == nil {
		return
	}

	for i := 0; i < len(s.guilds[p.GuildID].Presences); i++ {
		if s.guilds[p.GuildID].Presences[i].User.ID == p.User.ID {
			s.guilds[p.GuildID].Presences[i] = *p
		}
	}
}

// updateUser updates a user in the Users map (or the User) as well
// as in all the guilds this user is.
func (s *State) updateUser(u *User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if u.ID == s.currentUser.ID {
		s.currentUser = u
	} else {
		s.users[u.ID] = u
	}

	for id := range s.guilds {
		g := s.guilds[id]
		for i := 0; i < len(g.Members); i++ {
			if g.Members[i].User.ID == u.ID {
				g.Members[i].User = u
			}
		}
	}
}

// updateChannel updates a channel in the channel map as well as in
// the guild this channel is for guild text, voice and category channels.
// If the channel does not exist yet, it is added.
func (s *State) updateChannel(c *Channel) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch c.Type {
	case channel.TypeDM:
		s.dms[c.ID] = c

	case channel.TypeGroupDM:
		s.groups[c.ID] = c

	case channel.TypeGuildText, channel.TypeGuildVoice, channel.TypeGuildCategory:
		chs := s.guilds[c.GuildID].Channels
		var found bool
		for i := 0; i < len(chs); i++ {
			if chs[i].ID == c.ID {
				chs[i] = *c
				found = true
				break
			}
		}
		if !found {
			chs = append(chs, *c)
		}
	}

	s.channels[c.ID] = c
}

// removeChannel removes the given channel from the channels map as
// well as the guild it was in for guild text, voice and category channels.
func (s *State) removeChannel(c *Channel) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch c.Type {
	case channel.TypeDM:
		delete(s.dms, c.ID)

	case channel.TypeGroupDM:
		delete(s.groups, c.ID)

	case channel.TypeGuildText, channel.TypeGuildVoice, channel.TypeGuildCategory:
		g := s.guilds[c.GuildID]
		if g == nil {
			return
		}
		chs := g.Channels
		for i := 0; i < len(chs); i++ {
			if chs[i].ID == c.ID {
				chs = append(chs[:i], chs[i+1:]...)
				s.guilds[c.GuildID].Channels = chs
				break
			}
		}
	}

	delete(s.channels, c.ID)
}

// updatePins updates the LastPinTimestamp of a channel in the Channel map
// and in the DM, group DM or guild this channel is in.
func (s *State) updatePins(p *ChannelPinsUpdate) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ch := s.channels[p.ChannelID]
	if ch == nil {
		return
	}

	ch.LastPinTimestamp = p.LastPinTimestamp
	switch ch.Type {
	case channel.TypeDM:
		s.dms[p.ChannelID].LastPinTimestamp = p.LastPinTimestamp

	case channel.TypeGroupDM:
		s.groups[p.ChannelID].LastPinTimestamp = p.LastPinTimestamp

	case channel.TypeGuildText:
		chs := s.guilds[ch.GuildID].Channels
		for i := 0; i < len(chs); i++ {
			if chs[i].ID == p.ChannelID {
				chs[i].LastPinTimestamp = p.LastPinTimestamp
				break
			}
		}
	}
}

func (s *State) guildMemberAdd(m *GuildMemberAdd) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// This is a new user, add it to the users map.
	if s.users[m.User.ID] == nil {
		s.users[m.User.ID] = m.User
	}

	g := s.guilds[m.GuildID]
	if g == nil {
		return
	}

	g.MemberCount++
	g.Members = append(g.Members, *m.GuildMember)
}

func (s *State) guildMemberUpdate(m *GuildMemberUpdate) {
	s.mu.Lock()
	defer s.mu.Unlock()

	g := s.guilds[m.GuildID]
	if g == nil {
		return
	}

	for i := 0; i < len(g.Members); i++ {
		if g.Members[i].User.ID == m.User.ID {
			g.Members[i].Roles = m.Roles
			g.Members[i].User = m.User
			g.Members[i].Nick = m.Nick
		}
	}
}

func (s *State) guildMemberRemove(r *GuildMemberRemove) {
	s.mu.Lock()
	defer s.mu.Unlock()

	g := s.guilds[r.GuildID]
	if g == nil {
		return
	}

	g.MemberCount--
	for i := 0; i < len(g.Members); i++ {
		if g.Members[i].User.ID == r.User.ID {
			g.Members = append(g.Members[:i], g.Members[i+1:]...)
			break
		}
	}

	var found bool
	for id := range s.guilds {
		g := s.guilds[id]
		for _, m := range g.Members {
			if m.User.ID == r.User.ID {
				found = true
				break
			}
		}
	}
	// If this user is in no other guild, remove it from the state.
	if !found {
		delete(s.users, r.User.ID)
	}
}

// guildRoleCreate adds a role to a guild.
func (s *State) guildRoleCreate(gr *GuildRole) {
	s.mu.Lock()
	defer s.mu.Unlock()

	g := s.guilds[gr.GuildID]
	if g == nil {
		return
	}

	g.Roles = append(g.Roles, *gr.Role)
}

// guildRoleUpdate updates a role in a guild.
func (s *State) guildRoleUpdate(gr *GuildRole) {
	g := s.guilds[gr.GuildID]
	if g == nil {
		return
	}

	for i := 0; i < len(g.Roles); i++ {
		if g.Roles[i].ID == gr.Role.ID {
			g.Roles[i] = *gr.Role
			break
		}
	}
}

// guildRoleRemove removes a role from a guild.
func (s *State) guildRoleRemove(gr *GuildRoleDelete) {
	s.mu.Lock()
	defer s.mu.Unlock()

	g := s.guilds[gr.GuildID]
	if g == nil {
		return
	}

	for i := 0; i < len(g.Roles); i++ {
		if g.Roles[i].ID == gr.RoleID {
			g.Roles = append(g.Roles[:i], g.Roles[i+1:]...)
			break
		}
	}
}

// setRTT sets the Round Trip Time. See the RTT method for more information
// on how it is calculated.
func (s *State) setRTT(d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.rtt = d
}
