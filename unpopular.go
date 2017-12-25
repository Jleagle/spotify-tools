package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
	"github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
	spot "github.com/zmb3/spotify"
)

func unpopularHandler(w http.ResponseWriter, r *http.Request) {

	vars := structs.TemplateVars{}

	err := session.Write(w, r, session.LastPage, "/unpopular")
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

	// Get albums
	client := spotify.GetClient(r)

	vars.SearchAlbums, err = client.SearchOpt("tag:new", spot.SearchTypeAlbum, spotify.GetOptions(r, 3, 0, ""))
	if err != nil {
		returnTemplate(w, r, "error", vars, err)
		return
	}

	returnTemplate(w, r, "unpopular", vars, err)
	return
}
