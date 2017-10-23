package main

import (
	"net/http"

	"github.com/Jleagle/spotify-tools/session"
	spot "github.com/Jleagle/spotify-tools/spotify"
)

func shuffleHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, "shuffle")

	vars := templateVars{}

	if session.IsLoggedIn(r) {
		client := spot.GetClient(r)
		playlist, err := client.GetPlaylistsForUser("jleagle")
		if err != nil {
			//todo
		}

		// todo, sort playlists?
		vars.Playlists = playlist.Playlists
	}

	returnTemplate(w, r, "shuffle", vars, nil)
	return
}
