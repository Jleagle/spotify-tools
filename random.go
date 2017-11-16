package main

import (
	"math/rand"
	"net/http"

	"github.com/Jleagle/go-helpers/helpers"
	"github.com/Jleagle/spotifyhelper/session"
	spot "github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
	"github.com/go-chi/chi"
	"github.com/kr/pretty"
	"github.com/zmb3/spotify"
)

func randomHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, "/random")

	if !session.IsLoggedIn(r) {
		returnLoggedOutTemplate(w, r, nil)
		return
	}

	var err error
	vars := structs.TemplateVars{}

	client := spot.GetClient(r)
	search := helpers.RandomString(1, "aeiou")
	offset := rand.Intn(1000)

	switch chi.URLParam(r, "type") {
	case "albums", "":
		vars.SearchAlbums, err = client.SearchOpt(search, spotify.SearchTypeAlbum, spot.GetOptions(r, 3, offset))
	case "artists":
		vars.SearchArtists, err = client.SearchOpt(search, spotify.SearchTypeArtist, spot.GetOptions(r, 3, offset))
	case "tracks":
		vars.SearchTracks, err = client.SearchOpt(search, spotify.SearchTypeTrack, spot.GetOptions(r, 3, offset))
	case "playlists":
		vars.SearchPlaylists, err = client.SearchOpt(search, spotify.SearchTypePlaylist, spot.GetOptions(r, 3, offset))
	}

	if err != nil {
		pretty.Print(err)
	}

	returnTemplate(w, r, "random", vars, err)
	return
}
