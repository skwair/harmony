package debug

import (
	"encoding/json"
	"net/http"

	"github.com/skwair/harmony"
)

type httpDebugger struct {
	state *harmony.State
}

func NewHTTP(state *harmony.State) {
	d := httpDebugger{state}

	http.HandleFunc("/debug/state/index", d.index)
	http.HandleFunc("/debug/state/all", d.all)
}

func (d *httpDebugger) index(w http.ResponseWriter, r *http.Request) {
	state := struct {
		UsersCount             int `json:"users_count"`
		GuildsCount            int `json:"guilds_count"`
		PresencesCount         int `json:"presences_count"`
		ChannelsCount          int `json:"channels_count"`
		DMsCount               int `json:"dms_count"`
		GroupsCount            int `json:"groups_count"`
		UnavailableGuildsCount int `json:"unavailable_guilds_count"`
	}{
		UsersCount:             len(d.state.Users()),
		GuildsCount:            len(d.state.Guilds()),
		PresencesCount:         len(d.state.Presences()),
		ChannelsCount:          len(d.state.Channels()),
		DMsCount:               len(d.state.DMs()),
		GroupsCount:            len(d.state.GroupDMs()),
		UnavailableGuildsCount: len(d.state.UnavailableGuilds()),
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(state); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (d *httpDebugger) all(w http.ResponseWriter, r *http.Request) {
	state := struct {
		CurrentUser       *harmony.User                        `json:"current_user"`
		Users             map[string]*harmony.User             `json:"users"`
		Guilds            map[string]*harmony.Guild            `json:"guilds"`
		Presences         map[string]*harmony.Presence         `json:"presences"`
		Channels          map[string]*harmony.Channel          `json:"channels"`
		DMs               map[string]*harmony.Channel          `json:"dms"`
		Groups            map[string]*harmony.Channel          `json:"groups"`
		UnavailableGuilds map[string]*harmony.UnavailableGuild `json:"unavailable_guilds"`
	}{
		CurrentUser:       d.state.CurrentUser(),
		Users:             d.state.Users(),
		Guilds:            d.state.Guilds(),
		Presences:         d.state.Presences(),
		Channels:          d.state.Channels(),
		DMs:               d.state.DMs(),
		Groups:            d.state.GroupDMs(),
		UnavailableGuilds: d.state.UnavailableGuilds(),
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(state); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
