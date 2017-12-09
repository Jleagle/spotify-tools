package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
	spot "github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
	"github.com/go-chi/chi"
)

func topHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, "/top")

	if !session.IsLoggedIn(r) {
		returnLoggedOutTemplate(w, r, nil)
		return
	}

	var err error
	vars := structs.TemplateVars{}

	client := spot.GetClient(r)

	switch chi.URLParam(r, "type") {
	case "artists", "":
		vars.FullArtistPage, err = client.CurrentUsersTopArtistsOpt(spot.GetOptions(r, 50, 0))
	case "tracks":
		vars.FullTrackPage, err = client.CurrentUsersTopTracksOpt(spot.GetOptions(r, 50, 0))
	}

	if err != nil {
		returnTemplate(w, r, "error", vars, err)
		return
	}

	returnTemplate(w, r, "top", vars, nil)
	return
}
