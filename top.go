package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
	spot "github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
	"github.com/go-chi/chi"
)

func topHandler(w http.ResponseWriter, r *http.Request) {

	vars := structs.TemplateVars{}

	err := session.Write(w, r, session.LastPage, "/top")
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

	var dateRange string

	switch chi.URLParam(r, "range") {
	case "years":
		dateRange = "long"
	case "months":
		dateRange = "medium"
	case "weeks", "":
		dateRange = "short"
	}

	artistTrack := chi.URLParam(r, "type")
	if artistTrack == "" {
		artistTrack = "artists"
	}

	client := spot.GetClient(r)

	switch artistTrack {
	case "artists", "":
		vars.FullArtistPage, err = client.CurrentUsersTopArtistsOpt(spot.GetOptions(r, 50, 0, dateRange))
	case "tracks":
		vars.FullTrackPage, err = client.CurrentUsersTopTracksOpt(spot.GetOptions(r, 50, 0, dateRange))
	}

	if err != nil {
		returnTemplate(w, r, "error", vars, err)
		return
	}

	vars.TimeRange = dateRange
	vars.ArtistTrack = artistTrack

	returnTemplate(w, r, "top", vars, nil)
	return
}
