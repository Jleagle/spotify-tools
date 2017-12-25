package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
	"github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
	"github.com/go-chi/chi"
	spot "github.com/zmb3/spotify"
)

func playlistHandler(w http.ResponseWriter, r *http.Request) {

	vars := structs.TemplateVars{}

	err := session.Write(w, r, session.LastPage, r.URL.Path)
	if err != nil {
		returnTemplate(w, r, "error", vars, err)
		return
	}

	// Check if logged in
	loggedIn, err := session.IsLoggedIn(r)
	if err != nil {
		returnTemplate(w, r, "error", vars, err)
		return
	}
	if !loggedIn {
		returnLoggedOutTemplate(w, r, nil)
		return
	}

	// Get playlist
	client := spotify.GetClient(r)

	vars.Playlist, err = client.GetPlaylist(chi.URLParam(r, "user"), spot.ID(chi.URLParam(r, "playlist")))
	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = "Can't get playlist"
		returnTemplate(w, r, "error", vars, err)
		return
	}

	returnTemplate(w, r, "playlist", vars, err)
	return
}
