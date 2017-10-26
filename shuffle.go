package main

import (
	"fmt"
	"math"
	"net/http"
	"sync"

	"github.com/Jleagle/go-helpers/helpers"
	"github.com/Jleagle/spotifyhelper/session"
	spot "github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
	"github.com/go-chi/chi"
	"github.com/zmb3/spotify"
)

func shuffleHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, "/shuffle")

	vars := structs.TemplateVars{}

	if session.IsLoggedIn(r) {

		playlists, err := spot.CurrentUsersPlaylists(r)
		if err != nil {
			returnTemplate(w, r, "error", structs.TemplateVars{}, err)
			return
		}
		vars.Playlists = playlists
	}

	returnTemplate(w, r, "shuffle", vars, nil)
	return
}

func shuffleActionHandler(w http.ResponseWriter, r *http.Request) {

	playlistID := chi.URLParam(r, "playlist")
	createNew := chi.URLParam(r, "new")

	username := session.Read(r, session.UserID)

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

	// Get tracks
	var trackChunks = int(math.Ceil(float64(playlist.Tracks.Total) / spot.TracksLimit))
	var waitGroup sync.WaitGroup
	var trackStrings = []string{}
	messages := make(chan string, playlist.Tracks.Total)

	for i := 0; i < trackChunks; i++ {

		waitGroup.Add(1)
		go func(chunk int) {
			defer waitGroup.Done()

			options := spot.GetOptions(r, spot.TracksLimit, chunk*spot.TracksLimit)
			tracks, err := client.GetPlaylistTracksOpt(username, playlist.ID, options, "")
			if err != nil {
				fmt.Println("Getting tracks: " + err.Error())
				//return
			}

			for _, v := range tracks.Tracks {
				messages <- string(v.Track.ID)
			}

		}(i)
	}
	waitGroup.Wait()
	close(messages)

	for message := range messages {
		trackStrings = append(trackStrings, message)
	}

	if len(trackStrings) != playlist.Tracks.Total {
		session.SetFlash(w, r, "Track count mismatch :(")
		http.Redirect(w, r, "/shuffle?highlight="+playlistID, 302)
		return
	}

	// Check tracks count
	if len(trackStrings) < 2 {
		session.SetFlash(w, r, "Playlist does not have enough track to shuffle.")
		http.Redirect(w, r, "/shuffle?highlight="+playlistID, 302)
		return
	}

	// Create new playlist
	if createNew == "1" {
		playlist, err = client.CreatePlaylistForUser(username, playlist.Name+" Shuffled", false)
		if err != nil {
			session.SetFlash(w, r, "Unable to create new playlist: "+err.Error())
			http.Redirect(w, r, "/shuffle?highlight="+playlistID, 302)
			return
		}

		playlistID = string(playlist.ID)
	}

	// Shuffle
	helpers.ShuffleStrings(trackStrings)

	// Delete tracks
	err = client.ReplacePlaylistTracks("jleagle", playlist.ID)
	if err != nil {
		session.SetFlash(w, r, err.Error())
		http.Redirect(w, r, "/shuffle?highlight="+playlistID, 302)
		return
	}

	// Chunk the tracks to be added back
	chunks := helpers.ArrayChunk(trackStrings, spot.TracksLimit)
	for _, chunk := range chunks {

		// Convert back to IDs
		trackIds := []spotify.ID{}
		for _, v := range chunk {
			trackIds = append(trackIds, spotify.ID(v))
		}

		// Not in a go routine as that seemed to not always work..
		client.AddTracksToPlaylist(username, playlist.ID, trackIds...)
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
