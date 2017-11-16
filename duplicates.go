package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
	spot "github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
)

func duplicatesHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, "/duplicates")

	if !session.IsLoggedIn(r) {
		returnLoggedOutTemplate(w, r, nil)
		return
	}

	var err error
	vars := structs.TemplateVars{}

	vars.Playlists, err = spot.CurrentUsersPlaylists(r)
	if err != nil {
		returnTemplate(w, r, "error", structs.TemplateVars{}, err)
		return
	}

	returnTemplate(w, r, "duplicates", vars, nil)
	return
}

func duplicatesActionHandler(w http.ResponseWriter, r *http.Request) {

	if !session.IsLoggedIn(r) {
		returnLoggedOutTemplate(w, r, nil)
		return
	}

	//playlistID := chi.URLParam(r, "playlist")
	//createNew := chi.URLParam(r, "new")

}
