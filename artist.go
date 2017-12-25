package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
	"github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
	"github.com/go-chi/chi"
	spot "github.com/zmb3/spotify"
)

func artistHandler(w http.ResponseWriter, r *http.Request) {

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

	id := spot.ID(chi.URLParam(r, "artist"))

	client := spotify.GetClient(r)

	// Get artist
	vars.Artist, err = client.GetArtist(id)
	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = "Can't find artist"
		returnTemplate(w, r, "error", vars, err)
		return
	}

	// Get country
	country, err := session.Read(r, session.UserCountry)
	if err != nil {
		returnTemplate(w, r, "error", vars, err)
		return
	}

	// Get top tracks
	vars.Tracks, err = client.GetArtistsTopTracks(id, country)
	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = "Can't find artists top tracks"
		returnTemplate(w, r, "error", vars, err)
		return
	}

	// Get albums
	vars.Albums, err = client.GetArtistAlbumsOpt(id, spotify.GetOptions(r, spotify.MaxArtistAlbums, 0, ""), nil)
	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = "Can't find artists albums"
		returnTemplate(w, r, "error", vars, err)
		return
	}

	returnTemplate(w, r, "artist", vars, err)
	return
}
