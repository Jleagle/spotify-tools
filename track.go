package main

import (
	"net/http"
	"strings"

	"github.com/Jleagle/spotifyhelper/session"
	"github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
	"github.com/go-chi/chi"
	spot "github.com/zmb3/spotify"
)

func trackHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, r.URL.Path)

	if !session.IsLoggedIn(r) {
		returnLoggedOutTemplate(w, r, nil)
		return
	}

	trackID := spot.ID(chi.URLParam(r, "track"))
	client := spotify.GetClient(r)

	var err error
	vars := structs.TemplateVars{}

	vars.Track, err = client.GetTrack(trackID)
	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = "Can't find track"
		returnTemplate(w, r, "error", vars, err)
		return
	}

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

func albumHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, r.URL.Path)

	if !session.IsLoggedIn(r) {
		returnLoggedOutTemplate(w, r, nil)
		return
	}

	client := spotify.GetClient(r)

	var err error
	vars := structs.TemplateVars{}

	vars.Album, err = client.GetAlbum(spot.ID(chi.URLParam(r, "album")))
	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = err.Error() //todo, copy this to other places
		returnTemplate(w, r, "error", vars, err)
		return
	}
	vars.Album.AlbumType = strings.Title(vars.Album.AlbumType)

	returnTemplate(w, r, "album", vars, err)
	return
}

func artistHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, r.URL.Path)

	if !session.IsLoggedIn(r) {
		returnLoggedOutTemplate(w, r, nil)
		return
	}

	id := spot.ID(chi.URLParam(r, "artist"))

	client := spotify.GetClient(r)

	var err error
	vars := structs.TemplateVars{}

	vars.Artist, err = client.GetArtist(id)
	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = "Can't find artist"
		returnTemplate(w, r, "error", vars, err)
		return
	}

	vars.Tracks, err = client.GetArtistsTopTracks(id, session.Read(r, session.UserCountry))
	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = "Can't find artists top tracks"
		returnTemplate(w, r, "error", vars, err)
		return
	}

	vars.Albums, err = client.GetArtistAlbumsOpt(id, spotify.GetOptions(r, spotify.MaxArtistAlbums, 0), nil)
	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = "Can't find artists albums"
		returnTemplate(w, r, "error", vars, err)
		return
	}

	returnTemplate(w, r, "artist", vars, err)
	return
}

func playlistHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, r.URL.Path)

	if !session.IsLoggedIn(r) {
		returnLoggedOutTemplate(w, r, nil)
		return
	}

	client := spotify.GetClient(r)

	var err error
	vars := structs.TemplateVars{}

	vars.Playlist, err = client.GetPlaylist(chi.URLParam(r, "user"), spot.ID(chi.URLParam(r, "playlist")))
	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = "Can't get playlist"
		returnTemplate(w, r, "error", vars, err)
		return
	}

	returnTemplate(w, r, "playlist", vars, err)
	return
}

func userHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, r.URL.Path)

	if !session.IsLoggedIn(r) {
		returnLoggedOutTemplate(w, r, nil)
		return
	}

	id := chi.URLParam(r, "user")

	client := spotify.GetClient(r)

	var err error
	vars := structs.TemplateVars{}

	vars.User, err = client.GetUsersPublicProfile(spot.ID(id))
	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = "Can't find user"
		returnTemplate(w, r, "error", vars, err)
		return
	}
	if vars.User.DisplayName == "" {
		vars.User.DisplayName = vars.User.ID
	}

	vars.UserPlaylists, err = client.GetPlaylistsForUser(id)
	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = "Can't find user playlists"
		returnTemplate(w, r, "error", vars, err)
		return
	}

	returnTemplate(w, r, "user", vars, err)
	return
}
