package main

import (
	"github.com/Jleagle/spotify-tools/session"
	"net/http"
)

func infoHandler(w http.ResponseWriter, r *http.Request) {

	session.Write(w, r, session.LastPage, "/info")

	returnTemplate(w, r, "info", templateVars{}, nil)
	return
}
