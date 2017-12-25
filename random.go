package main

import (
	"math/rand"
	"net/http"

	"github.com/Jleagle/go-helpers/helpers"
	"github.com/Jleagle/spotifyhelper/session"
	"github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
	"github.com/go-chi/chi"
	spot "github.com/zmb3/spotify"
)

func randomHandler(w http.ResponseWriter, r *http.Request) {

	vars := structs.TemplateVars{}

	err := session.Write(w, r, session.LastPage, "/random")
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

	client := spotify.GetClient(r)
	search := helpers.RandomString(1, "aeiou")
	offset := rand.Intn(1000)
	options := spotify.GetOptions(r, 3, offset, "")

	switch chi.URLParam(r, "type") {
	case "albums", "":
		vars.SearchAlbums, err = client.SearchOpt(search, spot.SearchTypeAlbum, options)
	case "artists":
		vars.SearchArtists, err = client.SearchOpt(search, spot.SearchTypeArtist, options)
	case "tracks":
		vars.SearchTracks, err = client.SearchOpt(search, spot.SearchTypeTrack, options)
	case "playlists":
		vars.SearchPlaylists, err = client.SearchOpt(search, spot.SearchTypePlaylist, options)
	}

	if err != nil {
		returnTemplate(w, r, "error", vars, err)
		return
	}

	returnTemplate(w, r, "random", vars, err)
	return
}
