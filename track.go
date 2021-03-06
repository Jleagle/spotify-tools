package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
	"github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
	"github.com/go-chi/chi"
	spot "github.com/zmb3/spotify"
)

func trackHandler(w http.ResponseWriter, r *http.Request) {

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

	trackID := spot.ID(chi.URLParam(r, "track"))
	client := spotify.GetClient(r)

	// Get track
	vars.Track, err = client.GetTrack(trackID)
	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = "Can't find track"
		returnTemplate(w, r, "error", vars, err)
		return
	}

	// Get audio features
	audioFeats, err := client.GetAudioFeatures(trackID)
	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = "Can't find audio features"
		returnTemplate(w, r, "error", vars, err)
		return
	}
	vars.AudioFeatures = audioFeats[0]

	returnTemplate(w, r, "track", vars, err)
	return
}
