package main

import (
	"fmt"
	"net/http"

	"github.com/Jleagle/go-helpers/helpers"
	"github.com/Jleagle/spotify-tools/session"
	spot "github.com/Jleagle/spotify-tools/spotify"
	"github.com/go-chi/chi"
	"github.com/zmb3/spotify"
	"github.com/kr/pretty"
	"math"
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

	playlistID := chi.URLParam(r, "playlist")
	createNew := chi.URLParam(r, "new")

	username := "jleagle"

	client := spot.GetClient(r)

	// Get playlist
	playlist, err := client.GetPlaylist(username, spotify.ID(playlistID))
	if err != nil {
		if err.Error() == "Not found." {
			session.SetFlash(w, r, "Playlist not found in your account")
		} else {
			session.SetFlash(w, r, err.Error())
		}

		http.Redirect(w, r, "/shuffle?highlight="+playlistID, 302)
		return
	}

	trackChunks := math.Ceil(float64(playlist.Tracks.Total)/100)

	for i := 1; i <= trackChunks; i++ {
		fmt.Println(i)
	}

	return

	// Get tracks
	// todo, make sure this gets every track so we dont lose any

	options:= spot.GetOptions(100, 0)

	tracks, err := client.GetPlaylistTracksOpt(username, playlist.ID, options, "")
	if err != nil {
		session.SetFlash(w, r, err.Error())
		http.Redirect(w, r, "/shuffle?highlight="+playlistID, 302)
		return
	}

	// Check tracks count
	if len(tracks.Tracks) < 2 {
		session.SetFlash(w, r, "Playlist does not have enough track to shuffle.")
		http.Redirect(w, r, "/shuffle?highlight="+playlistID, 302)
		return
	}

	pretty.Print(len(tracks.Tracks))
	return

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

	// Create new playlist
	if createNew == "1" {
		playlist, err = client.CreatePlaylistForUser(username, playlist.Name+" Shuffled", false)
		if err != nil {
			session.SetFlash(w, r, err.Error())
			http.Redirect(w, r, "/shuffle?highlight="+playlistID, 302)
			return
		}

		playlistID = string(playlist.ID)
	}

	// Replace tracks
	err = client.ReplacePlaylistTracks("jleagle", playlist.ID, trackIds...)
	if err != nil {
		session.SetFlash(w, r, err.Error())
		http.Redirect(w, r, "/shuffle?highlight="+playlistID, 302)
		return
	}

	// Set flash
	if createNew == "1" {
		session.SetFlash(w, r, "New playlist created!")
	} else {
		session.SetFlash(w, r, "Playlist shuffled!")
	}

	// Redirect
	http.Redirect(w, r, "/shuffle?highlight="+playlistID, 302)
	return
}

//func getPlaylistTracks(offset int) (count int, err error){
//
//	options:= spot.GetOptions(100, offset)
//
//	tracks, err := client.GetPlaylistTracksOpt(username, playlist.ID, options, "")
//	if err != nil {
//		session.SetFlash(w, r, err.Error())
//		http.Redirect(w, r, "/shuffle?highlight="+playlistID, 302)
//		return
//	}
//
//
//}
