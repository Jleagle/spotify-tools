package main

import (
	"fmt"
	"net/http"

	"github.com/Jleagle/go-helpers/helpers"
	"github.com/Jleagle/spotify-tools/session"
	spot "github.com/Jleagle/spotify-tools/spotify"
	"github.com/go-chi/chi"
	"github.com/kr/pretty"
	"github.com/zmb3/spotify"
)

func shuffleHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, "shuffle")

	vars := templateVars{}

	if session.IsLoggedIn(r) {

		client := spot.GetClient(r)
		options := spot.GetOptions(50, 0)

		playlist, err := client.CurrentUsersPlaylistsOpt(options)
		if err != nil {
			fmt.Println(err.Error())
		}

		for k, v := range playlist.Playlists {
			if v.Owner.DisplayName == "" {
				playlist.Playlists[k].Owner.DisplayName = v.Owner.ID
			}
		}

		// todo, sort playlists?
		vars.Playlists = playlist.Playlists
	}

	returnTemplate(w, r, "shuffle", vars, nil)
	return
}

func shuffleActionHandler(w http.ResponseWriter, r *http.Request) {

	playlist := chi.URLParam(r, "playlist")
	createNew := chi.URLParam(r, "new")

	client := spot.GetClient(r)

	// todo, make sure this gets every track so we dont lose any
	tracks, err := client.GetPlaylistTracks("jleagle", spotify.ID(playlist))
	if err != nil {
		pretty.Print(err.Error())
	}

	// Convert to strings
	trackStrings := []string{}
	for _, v := range tracks.Tracks {
		trackStrings = append(trackStrings, string(v.Track.ID))
	}

	// Shuffle
	helpers.ShuffleStrings(trackStrings)

	// Convert back to IDs
	trackIds := []spotify.ID{}
	for _, v := range trackStrings {
		trackIds = append(trackIds, spotify.ID(v))
	}

	if createNew == "1" {
		//todo
		session.SetFlash(w, r, "New playlist created!")
	} else {

		err = client.ReplacePlaylistTracks("jleagle", spotify.ID(playlist), trackIds...)
		if err != nil {
			pretty.Print(err.Error())
		}

		session.SetFlash(w, r, "Playlist shuffled!")
	}

	http.Redirect(w, r, "/shuffle", 302)
}
