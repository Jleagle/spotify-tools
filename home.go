package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
	"github.com/Jleagle/spotifyhelper/structs"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {

	vars := structs.TemplateVars{}

	err := session.Write(w, r, session.LastPage, "/")
	if err != nil {
		returnTemplate(w, r, "error", vars, err)
		return
	}

	returnTemplate(w, r, "home", vars, nil)
	return
}
