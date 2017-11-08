package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
	"github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
	"github.com/go-chi/chi"
	"github.com/kr/pretty"
	spot "github.com/zmb3/spotify"
)

func trackHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, r.URL.Path)
	pretty.Print("Path:" + r.URL.Path)

	client := spotify.GetClient(r)

	var err error
	vars := structs.TemplateVars{}
	vars.Track, err = client.GetTrack(spot.ID(chi.URLParam(r, "track")))

	returnTemplate(w, r, "track", vars, err)
	return
}

func albumHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, r.URL.Path)
	pretty.Print("Path:" + r.URL.Path)

	client := spotify.GetClient(r)

	var err error
	vars := structs.TemplateVars{}
	vars.Album, err = client.GetAlbum(spot.ID(chi.URLParam(r, "album")))

	returnTemplate(w, r, "album", vars, err)
	return
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	session.Write(w, r, session.LastPage, r.URL.Path)

	client := spotify.GetClient(r)

	var err error
	vars := structs.TemplateVars{}
	vars.Artist, err = client.GetArtist(spot.ID(chi.URLParam(r, "artist")))

	if err != nil {
		vars.ErrorCode = "404"
		vars.ErrorMessage = "Can't find artist"
		returnTemplate(w, r, "error", vars, err)
	}

	returnTemplate(w, r, "artist", vars, err)
	return
}

func playlistHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, r.URL.Path)
	pretty.Print("Path:" + r.URL.Path)

	client := spotify.GetClient(r)

	var err error
	vars := structs.TemplateVars{}
	vars.Playlist, err = client.GetPlaylist(chi.URLParam(r, "user"), spot.ID(chi.URLParam(r, "playlist")))

	returnTemplate(w, r, "playlist", vars, err)
	return
}

func userHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, r.URL.Path)
	pretty.Print("Path:" + r.URL.Path)

	client := spotify.GetClient(r)

	var err error
	vars := structs.TemplateVars{}
	vars.User, err = client.GetUsersPublicProfile(spot.ID(chi.URLParam(r, "user")))

	returnTemplate(w, r, "user", vars, err)
	return
}
