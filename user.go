package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
	"github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
	"github.com/go-chi/chi"
	spot "github.com/zmb3/spotify"
)

func userHandler(w http.ResponseWriter, r *http.Request) {

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

	id := chi.URLParam(r, "user")

	client := spotify.GetClient(r)

	// Get profile
	vars.User, err = client.GetUsersPublicProfile(spot.ID(id))
	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = "Can't find user"
		returnTemplate(w, r, "error", vars, err)
		return
	}
	if vars.User.DisplayName == "" {
		vars.User.DisplayName = vars.User.ID
	}

	// Get playlists
	vars.UserPlaylists, err = client.GetPlaylistsForUser(id)
	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = "Can't find user playlists"
		returnTemplate(w, r, "error", vars, err)
		return
	}

	returnTemplate(w, r, "user", vars, err)
	return
}
