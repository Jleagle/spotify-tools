package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
	"github.com/Jleagle/spotifyhelper/spotify"
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

	var dateRange = chi.URLParam(r, "range")
	if !helpers.InArray(dateRange, []string{"long", "medium", "short"}) {
		dateRange = "short"
	}

	artistTrack := chi.URLParam(r, "type")
	if !helpers.InArray(artistTrack, []string{"artists", "tracks"}) {
		artistTrack = "artists"
	}

	client := spotify.GetClient(r)

	switch artistTrack {
	case "artists":
		vars.FullArtistPage, err = client.CurrentUsersTopArtistsOpt(spotify.GetOptions(r, 50, 0, dateRange))
	case "tracks":
		vars.FullTrackPage, err = client.CurrentUsersTopTracksOpt(spotify.GetOptions(r, 50, 0, dateRange))
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
