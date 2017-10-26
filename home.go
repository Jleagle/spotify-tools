package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
	"github.com/Jleagle/spotifyhelper/structs"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, "/")

	returnTemplate(w, r, "home", structs.TemplateVars{}, nil)
	return
}
