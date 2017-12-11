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
