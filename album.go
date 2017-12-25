package main

import (
	"net/http"
	"strings"

	"github.com/Jleagle/spotifyhelper/session"
	"github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
	"github.com/go-chi/chi"
	spot "github.com/zmb3/spotify"
)

func albumHandler(w http.ResponseWriter, r *http.Request) {

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

	// Get album
	client := spotify.GetClient(r)

	vars.Album, err = client.GetAlbum(spot.ID(chi.URLParam(r, "album")))
	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = err.Error() //todo, copy this to other places
		returnTemplate(w, r, "error", vars, err)
		return
	}
	vars.Album.AlbumType = strings.Title(vars.Album.AlbumType)

	returnTemplate(w, r, "album", vars, err)
	return
}
