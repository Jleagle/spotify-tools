package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"sync"

	"github.com/Jleagle/spotifyhelper/logging"
	"github.com/Jleagle/spotifyhelper/session"
	"github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
	"github.com/go-chi/chi"
	spot "github.com/zmb3/spotify"
)

func shuffleHandler(w http.ResponseWriter, r *http.Request) {

	vars := structs.TemplateVars{}

	err := session.Write(w, r, session.LastPage, "/shuffle")
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

	// Get playlists
	vars.Playlists, err = spotify.CurrentUsersPlaylists(r)
	if err != nil {
		returnTemplate(w, r, "error", vars, err)
		return
	}

	returnTemplate(w, r, "shuffle", vars, nil)
	return
}

func shuffleActionHandler(w http.ResponseWriter, r *http.Request) {

	vars := structs.TemplateVars{}

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

	playlistID := chi.URLParam(r, "playlist")
	createNew := chi.URLParam(r, "new")

	username, err := session.Read(r, session.UserID)
	if err != nil {
		logging.Error(err)
	}

	client := spotify.GetClient(r)

	// Get playlist
	playlist, err := client.GetPlaylist(username, spot.ID(playlistID))
	if err != nil {

		if err.Error() == "Not found." {
			session.SetFlash(w, r, "Playlist not found in your account")
		} else {
			session.SetFlash(w, r, err.Error())
		}

		logging.Error(err)
		http.Redirect(w, r, "/shuffle?highlight="+playlistID, 302)
		return
	}

	// Get tracks
	var trackChunks = int(math.Ceil(float64(playlist.Tracks.Total) / spotify.TracksLimit))
	var waitGroup sync.WaitGroup
	var trackStrings []string
	messages := make(chan string, playlist.Tracks.Total)

	for i := 0; i < trackChunks; i++ {

		waitGroup.Add(1)
		go func(chunk int) {
			defer waitGroup.Done()

			options := spotify.GetOptions(r, spotify.TracksLimit, chunk*spotify.TracksLimit, "")
			tracks, err := client.GetPlaylistTracksOpt(username, playlist.ID, options, "")
			if err != nil {

				logging.Info(err)
				session.SetFlash(w, r, "Failed to get playlist tracks")
				http.Redirect(w, r, "/shuffle?highlight="+playlistID, 302)
				return
			}

			for _, v := range tracks.Tracks {
				messages <- string(v.Track.ID)
			}

		}(i)
	}
	waitGroup.Wait()
	close(messages)

	var invalidTracks = 0
	for message := range messages {
		if message == "" {
			invalidTracks++
		} else {
			trackStrings = append(trackStrings, message)
		}
	}

	// Show message about invalid tracks
	if invalidTracks > 0 {
		session.SetFlash(w, r, fmt.Sprintf("%v %s", invalidTracks, " invalid tracks had to be removed"))
	}

	// Check we have as many tracks as we should
	if len(trackStrings)+invalidTracks != playlist.Tracks.Total {

		session.SetFlash(w, r, "Track count mismatch :(")
		http.Redirect(w, r, "/shuffle?highlight="+playlistID, 302)
		return
	}

	// Check playlist is worth shuffling
	if len(trackStrings) < 2 {

		session.SetFlash(w, r, "Playlist does not have enough track to shuffle.")
		http.Redirect(w, r, "/shuffle?highlight="+playlistID, 302)
		return
	}

	// Create new playlist
	if createNew == "1" {
		playlist, err = client.CreatePlaylistForUser(username, playlist.Name+" Shuffled", false)
		if err != nil {

			logging.Info(err)
			session.SetFlash(w, r, "Unable to create new playlist: "+err.Error())
			http.Redirect(w, r, "/shuffle?highlight="+playlistID, 302)
			return
		}

		playlistID = string(playlist.ID)
	}

	// Shuffle
	shuffleStrings(trackStrings)

	// Delete tracks
	err = client.ReplacePlaylistTracks("jleagle", playlist.ID)
	if err != nil {

		logging.Info(err)
		session.SetFlash(w, r, err.Error())
		http.Redirect(w, r, "/shuffle?highlight="+playlistID, 302)
		return
	}

	// Chunk the tracks to be added back
	chunks := arrayChunk(trackStrings, spotify.TracksLimit)
	for _, chunk := range chunks {

		// Convert back to IDs
		var trackIds []spot.ID
		for _, v := range chunk {
			trackIds = append(trackIds, spot.ID(v))
		}

		waitGroup.Add(1)
		func() {
			defer waitGroup.Done()
			_, err = client.AddTracksToPlaylist(username, playlist.ID, trackIds...)
			if err != nil {
				logging.Info(err)
			}
		}()
	}
	waitGroup.Wait()

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

func shuffleStrings(slc []string) {
	N := len(slc)
	for i := 0; i < N; i++ {
		// choose index uniformly in [i, N-1]
		r := i + rand.Intn(N-i)
		slc[r], slc[i] = slc[i], slc[r]
	}
}

func arrayChunk(logs []string, chunkSize int) (divided [][]string) {

	for i := 0; i < len(logs); i += chunkSize {
		end := i + chunkSize

		if end > len(logs) {
			end = len(logs)
		}

		divided = append(divided, logs[i:end])
	}

	return divided
}
