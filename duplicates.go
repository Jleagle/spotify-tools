package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
	spot "github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
)

func duplicatesHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, "/duplicates")

	vars := structs.TemplateVars{}

	if session.IsLoggedIn(r) {

		playlists, err := spot.CurrentUsersPlaylists(r)
		if err != nil {
			returnTemplate(w, r, "error", structs.TemplateVars{}, err)
			return
		}
		vars.Playlists = playlists
	}

	returnTemplate(w, r, "duplicates", vars, nil)
	return
}
