package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
	spot "github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
	"github.com/go-chi/chi"
	"github.com/kr/pretty"
	"github.com/zmb3/spotify"
)

func randomHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, "/random")

	randomType := chi.URLParam(r, "type")

	client := spot.GetClient(r)

	vars := structs.TemplateVars{}

	searchString := "a OR e OR i OR o OR u"

	var err error

	switch randomType {
	case "albums":
		vars.SearchAlbums, err = client.SearchOpt(searchString, spotify.SearchTypeAlbum, spot.GetOptions(r, 3, 99))
	case "artists":
		vars.SearchArtists, err = client.SearchOpt(searchString, spotify.SearchTypeArtist, spot.GetOptions(r, 3, 99))
	case "playlists":
		vars.SearchPlaylists, err = client.SearchOpt(searchString, spotify.SearchTypePlaylist, spot.GetOptions(r, 3, 99))
	default:
		vars.SearchTracks, err = client.SearchOpt(searchString, spotify.SearchTypeTrack, spot.GetOptions(r, 3, 99))
	}

	if err != nil {
		pretty.Print(err)
	}

	returnTemplate(w, r, "random", vars, err)
	return
}
