package harmony

import (
	"sync"
	"time"

	"github.com/skwair/harmony/discord"
	"github.com/skwair/harmony/voice"
)

// State is a cache of the state of the application that is updated in real-time
// as events are received from the Gateway.
// Objects returned by State methods are snapshots of original objects used
// internally by the State. This means they are safe to be used and modified
// but they won't be updated as new events are received.
type State struct {
	mu sync.RWMutex

	me                *discord.User
	users             map[string]*discord.User
	guilds            map[string]*discord.Guild
	presences         map[string]*discord.Presence // Presence by user ID.
	channels          map[string]*discord.Channel
	dms               map[string]*discord.Channel
	groups            map[string]*discord.Channel
	unavailableGuilds map[string]*discord.UnavailableGuild

	rtt time.Duration

	// NOTE: consider adding statistics such as the uptime, ping, number
	// of voice connections, etc... in the state.
}

// newState returns a new initialized state, ready to be used.
func newState() *State {
	return &State{
		users:             make(map[string]*discord.User),
		guilds:            make(map[string]*discord.Guild),
		presences:         make(map[string]*discord.Presence),
		channels:          make(map[string]*discord.Channel),
		dms:               make(map[string]*discord.Channel),
		groups:            make(map[string]*discord.Channel),
		unavailableGuilds: make(map[string]*discord.UnavailableGuild),
	}
}

// Me returns the current user from the state.
func (s *State) Me() *discord.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.me.Clone()
}

// User returns a user given its ID from the state.
func (s *State) User(id string) *discord.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.users[id].Clone()
}

// Guild returns a guild given its ID from the state.
func (s *State) Guild(id string) *discord.Guild {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.guilds[id].Clone()
}

// Channel returns a channel given its ID from the state.
func (s *State) Channel(id string) *discord.Channel {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.channels[id].Clone()
}

// GroupDM returns a group DM given its ID from the state.
func (s *State) GroupDM(id string) *discord.Channel {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.groups[id].Clone()
}

// DM returns a DM given its ID from the state.
func (s *State) DM(id string) *discord.Channel {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.dms[id].Clone()
}

// Presence returns a presence given a user ID from the state.
func (s *State) Presence(userID string) *discord.Presence {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.presences[userID].Clone()
}

// UnavailableGuild returns an unavailable guild given its ID from the state.
func (s *State) UnavailableGuild(id string) *discord.UnavailableGuild {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.unavailableGuilds[id].Clone()
}

// Users returns a map of user ID to user from the state.
func (s *State) Users() map[string]*discord.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	newMap := make(map[string]*discord.User)
	for k, v := range s.users {
		newMap[k] = v.Clone()
	}

	return newMap
}

// Guilds returns a map of guild ID to guild from the state.
func (s *State) Guilds() map[string]*discord.Guild {
	s.mu.RLock()
	defer s.mu.RUnlock()

	newMap := make(map[string]*discord.Guild)
	for k, v := range s.guilds {
		newMap[k] = v.Clone()
	}

	return newMap
}

// Channels returns a map of channels ID to channels from the state.
func (s *State) Channels() map[string]*discord.Channel {
	s.mu.RLock()
	defer s.mu.RUnlock()

	newMap := make(map[string]*discord.Channel)
	for k, v := range s.channels {
		newMap[k] = v.Clone()
	}

	return newMap
}

// GroupDMs returns a map of group DM ID to group DM from the state.
func (s *State) GroupDMs() map[string]*discord.Channel {
	s.mu.RLock()
	defer s.mu.RUnlock()

	newMap := make(map[string]*discord.Channel)
	for k, v := range s.groups {
		newMap[k] = v.Clone()
	}

	return newMap
}

// DMs returns a map of DM ID to DM from the state.
func (s *State) DMs() map[string]*discord.Channel {
	s.mu.RLock()
	defer s.mu.RUnlock()

	newMap := make(map[string]*discord.Channel)
	for k, v := range s.dms {
		newMap[k] = v.Clone()
	}

	return newMap
}

// Presences returns a map of user ID to presence from the state.
func (s *State) Presences() map[string]*discord.Presence {
	s.mu.RLock()
	defer s.mu.RUnlock()

	newMap := make(map[string]*discord.Presence)
	for k, v := range s.presences {
		newMap[k] = v.Clone()
	}

	return newMap
}

// UnavailableGuilds returns a map of guild ID to unavailable guild from the state.
func (s *State) UnavailableGuilds() map[string]*discord.UnavailableGuild {
	s.mu.RLock()
	defer s.mu.RUnlock()

	newMap := make(map[string]*discord.UnavailableGuild)
	for k, v := range s.unavailableGuilds {
		newMap[k] = v.Clone()
	}

	return newMap
}

// RTT returns the Round Trip Time between the client and Discord's Gateway.
// It is calculated and updated when sending heartbeat payloads (roughly
// every minute).
func (s *State) RTT() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.rtt
}

// setInitialState initializes the state with a Ready event received from the gateway.
func (s *State) setInitialState(r *Ready) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.me = r.User
	for i := 0; i < len(r.Guilds); i++ {
		g := &r.Guilds[i]
		s.guilds[g.ID] = &discord.Guild{
			ID:          g.ID,
			Name:        g.Name,
			Owner:       g.Owner,
			Permissions: g.Permissions,
		}
		if g.Icon != "" {
			s.guilds[g.ID].Icon = g.Icon
		}
	}
	for i := 0; i < len(r.PrivateChannels); i++ {
		dm := &r.PrivateChannels[i]
		if dm.Type == discord.ChannelTypeDM {
			s.dms[dm.ID] = dm
		}
		if dm.Type == discord.ChannelTypeGroupDM {
			s.groups[dm.ID] = dm
		}
	}
}

// updateGuild adds the given guild to the state. If it already
// exists, it merges its content with the existing guild.
// It also removes this guild from the UnavailableGuilds map if
// it was present.
func (s *State) updateGuild(g *discord.Guild) {
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
func (s *State) removeGuild(g *discord.UnavailableGuild) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.guilds, g.ID)
	s.unavailableGuilds[g.ID] = g
}

// updateGuildEmojis updates the emojis available in a guild if it
// is already tracked by the state, does nothing otherwise.
func (s *State) updateGuildEmojis(guildID string, emojis []discord.Emoji) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.guilds[guildID] != nil {
		s.guilds[guildID].Emojis = emojis
	}
}

// updateGuildVoiceStates updates the voice states in a guild if it is
// already tracked by the state, does nothing otherwise.
func (s *State) updateGuildVoiceStates(vsu *voice.StateUpdate) {
	s.mu.Lock()
	defer s.mu.Unlock()

	g := s.guilds[vsu.GuildID]
	if g == nil {
		return
	}

	// If we have a channel ID, then it means it is either a new voice
	// state or an update to an existing one.
	if vsu.ChannelID != nil {
		index := -1
		// Check if we already have a voice state for this user.
		// If we do, save the index of the voice state.
		for i, state := range g.VoiceStates {
			if state.UserID == vsu.UserID {
				index = i
			}
		}

		// This state is already tracked, update it.
		if index != -1 {
			g.VoiceStates[index] = vsu.State
		} else { // This is a new voice state, append it.
			g.VoiceStates = append(g.VoiceStates, vsu.State)
		}
	} else { // We have no channel ID, the user left the channel, remove it from the state.
		// Find the index of the voice state update to remove.
		var toRemove int
		for i, update := range g.VoiceStates {
			if update.UserID == vsu.UserID {
				toRemove = i
			}
		}

		// Remove it, without preserving the order of the slice.
		g.VoiceStates[toRemove] = g.VoiceStates[len(g.VoiceStates)-1]
		g.VoiceStates = g.VoiceStates[:len(g.VoiceStates)-1]
	}
}

// updatePresence updates a presence both in the presences map as
// well as in the Guilds map.
func (s *State) updatePresence(p *discord.Presence) {
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
func (s *State) updateUser(u *discord.User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if u.ID == s.me.ID {
		s.me = u
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
func (s *State) updateChannel(c *discord.Channel) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch c.Type {
	case discord.ChannelTypeDM:
		s.dms[c.ID] = c

	case discord.ChannelTypeGroupDM:
		s.groups[c.ID] = c

	case discord.ChannelTypeGuildText, discord.ChannelTypeGuildVoice, discord.ChannelTypeGuildCategory,
		discord.ChannelTypeGuildNews, discord.ChannelTypeGuildStore:
		guild := s.guilds[c.GuildID]
		var found bool
		for i := 0; i < len(guild.Channels); i++ {
			if guild.Channels[i].ID == c.ID {
				guild.Channels[i] = *c
				found = true
				break
			}
		}
		if !found {
			guild.Channels = append(guild.Channels, *c)
		}
	}

	s.channels[c.ID] = c
}

// removeChannel removes the given channel from the channels map as
// well as the guild it was in for guild text, voice and category channels.
func (s *State) removeChannel(c *discord.Channel) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch c.Type {
	case discord.ChannelTypeDM:
		delete(s.dms, c.ID)

	case discord.ChannelTypeGroupDM:
		delete(s.groups, c.ID)

	case discord.ChannelTypeGuildText, discord.ChannelTypeGuildVoice, discord.ChannelTypeGuildCategory:
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
	case discord.ChannelTypeDM:
		s.dms[p.ChannelID].LastPinTimestamp = p.LastPinTimestamp

	case discord.ChannelTypeGroupDM:
		s.groups[p.ChannelID].LastPinTimestamp = p.LastPinTimestamp

	case discord.ChannelTypeGuildText:
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
