package main

import (
	"net/http"

	"github.com/Jleagle/spotifyhelper/session"
	"github.com/Jleagle/spotifyhelper/structs"
)

func duplicatesHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, "/duplicates")

	returnTemplate(w, r, "duplicates", structs.TemplateVars{}, nil)
	return
}
