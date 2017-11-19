package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
	spot "github.com/Jleagle/spotifyhelper/spotify"
	"github.com/Jleagle/spotifyhelper/structs"
	"github.com/kr/pretty"
	"github.com/zmb3/spotify"
)

func unpopularHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, "/unpopular")

	if !session.IsLoggedIn(r) {
		returnLoggedOutTemplate(w, r, nil)
		return
	}

	var err error
	vars := structs.TemplateVars{}

	client := spot.GetClient(r)
	vars.SearchAlbums, err = client.SearchOpt("tag:new", spotify.SearchTypeAlbum, spot.GetOptions(r, 3, 0))

	if err != nil {
		pretty.Print(err)
	}

	returnTemplate(w, r, "unpopular", vars, err)
	return
}
